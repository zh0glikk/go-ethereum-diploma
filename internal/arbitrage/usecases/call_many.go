package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"sync/atomic"
	"time"
)

const (
	_timeout       = time.Millisecond * time.Duration(5000)
	CheckRevertAll = -1
)

func InitCallMany(ctx context.Context, b ethapi.Backend, blockNrOrHash rpc.BlockNumberOrHash, overrides *ethapi.StateOverride, timeout time.Duration, blockOverrides *ethapi.BlockOverrides) (*state.StateDB, *types.Header, vm.BlockContext, context.CancelFunc, error) {
	state, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, nil, vm.BlockContext{}, nil, err
	}
	if err := overrides.Apply(state); err != nil {
		return nil, nil, vm.BlockContext{}, nil, err
	}

	vmctx := core.NewEVMBlockContext(header, b.GetChainContext(), nil)
	blockOverrides.Apply(&vmctx)

	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	return state, header, vmctx, cancel, nil
}

func DoCallMany(ctx context.Context, b ethapi.Backend, args []ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *ethapi.StateOverride, timeout time.Duration, globalGasCap uint64, blockOverrides *ethapi.BlockOverrides) ([]models.CallManyResponseDTO, error) {
	defer func(start time.Time) { log.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	var results []*core.ExecutionResult

	state, header, vmctx, cancel, err := InitCallMany(ctx, b, blockNrOrHash, overrides, timeout, blockOverrides)
	if err != nil {
		return nil, err
	}
	defer cancel()

	for _, arg := range args {
		msg, err := arg.ToMessage(globalGasCap, header.BaseFee)
		if err != nil {
			return nil, err
		}

		evm, vmError, err := b.GetEVMWithContext(ctx, msg, state, header, &vm.Config{NoBaseFee: true}, vmctx)
		if err != nil {
			return nil, err
		}

		gp := new(core.GasPool).AddGas(math.MaxUint64)

		result, err := core.ApplyMessage(evm, msg, gp)
		if err := vmError(); err != nil {
			return nil, err
		}

		// If the timer caused an abort, return an appropriate error message
		if evm.Cancelled() {
			return nil, fmt.Errorf("execution aborted (timeout = %v)", timeout)
		}

		if err != nil {
			execResult := core.ExecutionResult{
				UsedGas:    0,
				Err:        fmt.Errorf("err: %w (supplied gas %d)", err, msg.GasLimit),
				ReturnData: nil,
			}

			results = append(results, &execResult)
			continue
		}

		results = append(results, result)
	}

	return convertCallManyResults(results), nil
}

func DoCallManyWithErrorDecoding(ctx context.Context, b ethapi.Backend, args []ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *ethapi.StateOverride, timeout time.Duration, globalGasCap uint64, blockOverrides *ethapi.BlockOverrides) ([]models.CallManyResponseDTO, error) {
	defer func(start time.Time) { log.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	var results []*core.ExecutionResult

	state, header, vmctx, cancel, err := InitCallMany(ctx, b, blockNrOrHash, overrides, timeout, blockOverrides)
	if err != nil {
		return nil, err
	}
	defer cancel()

	for _, arg := range args {
		msg, err := arg.ToMessage(globalGasCap, header.BaseFee)
		if err != nil {
			return nil, err
		}

		evm, vmError, err := b.GetEVMWithContext(ctx, msg, state, header, &vm.Config{NoBaseFee: true}, vmctx)
		if err != nil {
			return nil, err
		}

		gp := new(core.GasPool).AddGas(math.MaxUint64)

		result, err := core.ApplyMessage(evm, msg, gp)
		if err := vmError(); err != nil {
			return nil, err
		}

		// If the timer caused an abort, return an appropriate error message
		if evm.Cancelled() {
			return nil, fmt.Errorf("execution aborted (timeout = %v)", timeout)
		}

		if err != nil {
			execResult := core.ExecutionResult{
				UsedGas:    0,
				Err:        fmt.Errorf("err: %w (supplied gas %d)", err, msg.GasLimit),
				ReturnData: nil,
			}

			results = append(results, &execResult)
			continue
		}

		results = append(results, result)
	}

	return convertCallManyResultsWithDecoding(results), nil
}

func DoCallManyReturningState(
	ctx context.Context,
	b ethapi.Backend,
	args []ethapi.TransactionArgs,
	blockNrOrHash rpc.BlockNumberOrHash,
	overrides *ethapi.StateOverride,
	timeout time.Duration,
	globalGasCap uint64,
	blockOverrides *ethapi.BlockOverrides,
) ([]models.CallManyResponseDTO, *state.StateDB, *types.Header, vm.BlockContext, error) {
	defer func(start time.Time) { log.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	var results []models.CallManyResponseDTO

	stateDb, header, vmctx, _, err := InitCallMany(ctx, b, blockNrOrHash, overrides, timeout, blockOverrides)
	if err != nil {
		return nil, nil, nil, vm.BlockContext{}, err
	}

	for _, arg := range args {
		msg, err := arg.ToMessage(globalGasCap, header.BaseFee)
		if err != nil {
			return nil, nil, nil, vmctx, err
		}

		evm, vmError, err := b.GetEVMWithContext(ctx, msg, stateDb, header, &vm.Config{NoBaseFee: true}, vmctx)
		if err != nil {
			return nil, nil, nil, vmctx, err
		}

		gp := new(core.GasPool).AddGas(math.MaxUint64)

		result, err := core.ApplyMessage(evm, msg, gp)
		if err := vmError(); err != nil {
			return nil, nil, nil, vmctx, err
		}

		// If the timer caused an abort, return an appropriate error message
		if evm.Cancelled() {
			return nil, nil, nil, vmctx, fmt.Errorf("execution aborted (timeout = %v)", timeout)
		}

		if err != nil {
			execResult := core.ExecutionResult{
				UsedGas:    0,
				Err:        fmt.Errorf("err: %w (supplied gas %d)", err, msg.GasLimit),
				ReturnData: nil,
			}
			convertedResult, _ := convertCallManyResult(&execResult)
			results = append(results, convertedResult)
			continue
		}

		convertedResult, _ := convertCallManyResult(result)
		results = append(results, convertedResult)
	}

	return results, stateDb, header, vmctx, nil
}

func DoCallManyOnStateReturningState(
	ctx context.Context,
	b ethapi.Backend,
	args []ethapi.TransactionArgs,
	stateDB *state.StateDB,
	header *types.Header,
	vmctx vm.BlockContext,
	timeout time.Duration,
	globalGasCap uint64,
	checkRevertIndex int,
	returnGasIndexes []int,
) ([]models.CallManyResponseDTO, bool, []uint64, *state.StateDB, *types.Header, vm.BlockContext, error) {
	defer func(start time.Time) { log.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	var results []models.CallManyResponseDTO
	var spentGas []uint64

	for i := range args {
		msg, err := args[i].ToMessage(globalGasCap, header.BaseFee)
		if err != nil {
			return nil, false, spentGas, stateDB, header, vmctx, err
		}

		evm, vmError, err := b.GetEVMWithContext(ctx, msg, stateDB, header, &vm.Config{NoBaseFee: true}, vmctx)
		if err != nil {
			return nil, false, spentGas, stateDB, header, vmctx, err
		}

		gp := new(core.GasPool).AddGas(math.MaxUint64)

		result, err := core.ApplyMessage(evm, msg, gp)
		if err := vmError(); err != nil {
			return nil, false, spentGas, stateDB, header, vmctx, err
		}

		// If the timer caused an abort, return an appropriate error message
		if evm.Cancelled() {
			return nil, false, spentGas, stateDB, header, vmctx, fmt.Errorf("execution aborted (timeout = %v)", timeout)
		}

		if err != nil {
			execResult := core.ExecutionResult{
				UsedGas:    0,
				Err:        fmt.Errorf("err: %w (supplied gas %d)", err, msg.GasLimit),
				ReturnData: nil,
			}

			convertedResult, success := convertCallManyResult(&execResult)
			results = append(results, convertedResult)
			if (checkRevertIndex == CheckRevertAll || checkRevertIndex == i) && !success {
				return results, true, spentGas, stateDB, header, vmctx, nil
			}

			continue
		}

		convertedResult, success := convertCallManyResult(result)
		results = append(results, convertedResult)
		if (checkRevertIndex == CheckRevertAll || checkRevertIndex == i) && !success {
			return results, true, spentGas, stateDB, header, vmctx, nil
		}

		for _, returnGasIndex := range returnGasIndexes {
			if i == returnGasIndex {
				spentGas = append(spentGas, result.UsedGas)
			}
		}
	}

	return results, false, spentGas, stateDB, header, vmctx, nil
}

func SimpleBatchCallMany(ctx context.Context, b ethapi.Backend, args [][]ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *ethapi.StateOverride, timeout time.Duration, globalGasCap uint64, blockOverrides *ethapi.BlockOverrides) ([][]models.CallManyResponseDTO, error) {
	resp := make([][]models.CallManyResponseDTO, len(args))

	for i, arg := range args {
		results, err := DoCallMany(ctx, b, arg, blockNrOrHash, overrides, timeout, globalGasCap, blockOverrides)
		if err != nil {
			return nil, err
		}

		resp[i] = results
	}

	return resp, nil
}

func BatchCallMany(ctx context.Context, b ethapi.Backend, args [][]ethapi.TransactionArgs, blockNrOrHash []rpc.BlockNumberOrHash, overrides *[]*ethapi.StateOverride, blockOverrides *[]*ethapi.BlockOverrides) ([][]models.CallManyResponseDTO, error) {
	if len(args) != len(blockNrOrHash) {
		return nil, errors.New("length of args and blocks not equal")
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, _timeout)
	defer cancel()

	resp := make([][]models.CallManyResponseDTO, len(args))

	errChan := make(chan models.ChanCallManyError)
	respChan := make(chan models.ChanCallManyResponse)

	var counter atomic.Int64
	for i, arg := range args {
		if len(arg) == 0 {
			return nil, errors.New("invalid length of args")
		}
		go func(i int, arg []ethapi.TransactionArgs) {
			defer counter.Add(1)

			var override *ethapi.StateOverride
			if overrides != nil && len(*overrides) > i {
				override = (*overrides)[i]
			}

			var blockOverride *ethapi.BlockOverrides
			if blockOverrides != nil && len(*blockOverrides) > i {
				blockOverride = (*blockOverrides)[i]
			}

			results, err := DoCallMany(ctx, b, arg, blockNrOrHash[i], override, b.RPCEVMTimeout(), b.RPCGasCap(), blockOverride)
			if err != nil {
				errChan <- models.ChanCallManyError{Err: err, Ind: i}
				return
			}

			respChan <- models.ChanCallManyResponse{Response: results, Ind: i}
		}(i, arg)
	}

	for {
		select {
		case <-ctxTimeout.Done():
			return nil, ctxTimeout.Err()
		case err := <-errChan:
			log.Error(fmt.Sprintf("Error in eth_batchCallMany by index %v", err.Ind), "err", err.Err)
			resp[err.Ind] = []models.CallManyResponseDTO{}
		case r := <-respChan:
			resp[r.Ind] = r.Response
		default:
			if counter.Load() == int64(len(args)) {
				return resp, nil
			}
		}
	}
}

func SimpleBatchCallManyOnState(
	ctx context.Context,
	b ethapi.Backend,
	args [][]ethapi.TransactionArgs,
	stateDB *state.StateDB,
	header *types.Header,
	vmctx vm.BlockContext,
	timeout time.Duration,
	globalGasCap uint64) ([][]models.CallManyResponseDTO, error) {
	resp := make([][]models.CallManyResponseDTO, len(args))

	for i, arg := range args {
		results, _, _, _, _, _, err := DoCallManyOnStateReturningState(ctx, b, arg, stateDB.Copy(), header, vmctx, timeout, globalGasCap, 0, nil)
		if err != nil {
			return nil, err
		}

		resp[i] = results
	}

	return resp, nil
}
