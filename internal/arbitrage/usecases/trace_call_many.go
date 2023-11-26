package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
	"time"
)

const defaultTraceReexec = uint64(128)
const defaultTraceTimeout = 5 * time.Second

func TraceCallMany(ctx context.Context, b ethapi.Backend, txs []ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, config *tracers.TraceCallConfig) ([]interface{}, error) {
	config = setTraceCallConfigDefaultTracer(config)
	return prepareAndTraceCallMany(ctx, b, txs, blockNrOrHash, config)
}

func prepareAndTraceCallMany(ctx context.Context, b ethapi.Backend, args []ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, config *tracers.TraceCallConfig) ([]interface{}, error) {
	// Try to retrieve the specified block
	var (
		err   error
		block *types.Block
	)
	if hash, ok := blockNrOrHash.Hash(); ok {
		block, err = BlockByHash(ctx, b, hash)
	} else if number, ok := blockNrOrHash.Number(); ok {
		if number == rpc.PendingBlockNumber {
			// We don't have access to the miner here. For tracing 'future' transactions,
			// it can be done with block- and state-overrides instead, which offers
			// more flexibility and stability than trying to trace on 'pending', since
			// the contents of 'pending' is unstable and probably not a true representation
			// of what the next actual block is likely to contain.
			return nil, errors.New("tracing on top of pending is not supported")
		}
		block, err = BlockByNumber(ctx, b, number)
	} else {
		return nil, errors.New("invalid arguments; neither block nor hash specified")
	}
	if err != nil {
		return nil, err
	}
	// try to recompute the state
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	statedb, release, err := b.StateAtBlockCustom(ctx, block, reexec, nil, true, false)
	if err != nil {
		return nil, err
	}
	defer release()

	vmctx := core.NewEVMBlockContext(block.Header(), ethapi.NewChainContext(ctx, b), nil)
	// Apply the customization rules if required.
	if config != nil {
		if err := config.StateOverrides.Apply(statedb); err != nil {
			return nil, err
		}
		config.BlockOverrides.Apply(&vmctx)
	}

	var msgs []*core.Message

	for _, arg := range args {
		msg, err := arg.ToMessage(b.RPCGasCap(), block.BaseFee())
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	// Execute the trace

	var traceConfig *tracers.TraceConfig
	if config != nil {
		traceConfig = &config.TraceConfig
	}
	return traceTxs(ctx, b, msgs, new(tracers.Context), vmctx, statedb, traceConfig)
}

func traceTxs(ctx context.Context, b ethapi.Backend, messages []*core.Message, txctx *tracers.Context, vmctx vm.BlockContext, statedb *state.StateDB, config *tracers.TraceConfig) ([]interface{}, error) {
	var (
		tracer  tracers.Tracer
		err     error
		timeout = defaultTraceTimeout
	)
	if config == nil {
		config = &tracers.TraceConfig{}
	}
	// Default tracer is the struct logger
	tracer = logger.NewStructLogger(config.Config)
	if config.Tracer != nil {
		tracer, err = tracers.DefaultDirectory.New(*config.Tracer, txctx, config.TracerConfig)
		if err != nil {
			return nil, err
		}
	}

	// Define a meaningful timeout of a single transaction trace
	if config.Timeout != nil {
		if timeout, err = time.ParseDuration(*config.Timeout); err != nil {
			return nil, err
		}
	}

	var results []interface{}
	for _, message := range messages {
		// Call Prepare to clear out the statedb access list
		txContext := core.NewEVMTxContext(message)

		vmenv := vm.NewEVM(vmctx, txContext, statedb, b.ChainConfig(), vm.Config{Tracer: tracer, NoBaseFee: true})
		deadlineCtx, cancel := context.WithTimeout(ctx, timeout)
		go func() {
			<-deadlineCtx.Done()
			if errors.Is(deadlineCtx.Err(), context.DeadlineExceeded) {
				tracer.Stop(errors.New("execution timeout"))
				// Stop evm execution. Note cancellation is not necessarily immediate.
				vmenv.Cancel()
			}
		}()
		defer cancel()

		statedb.SetTxContext(txctx.TxHash, txctx.TxIndex)
		if _, err = core.ApplyMessage(vmenv, message, new(core.GasPool).AddGas(message.GasLimit)); err != nil {
			return nil, fmt.Errorf("tracing failed: %w", err)
		}

		result, err := tracer.GetResult()
		if err != nil {
			results = append(results, err)
			continue
		}

		results = append(results, result)
	}

	return results, err
}
