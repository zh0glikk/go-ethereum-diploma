package usecases

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/internal/arbitrage/algo"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	transact "github.com/ethereum/go-ethereum/internal/arbitrage/transactor"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/protocol"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/simulation_models"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/unpacker"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"time"
)

func ExecuteSwaps(
	ctx context.Context,
	b ethapi.Backend,
	request models.SwapBundle,
	stateOverride *ethapi.StateOverride,
	blockNrOrHash rpc.BlockNumberOrHash,
	blockOverrides *ethapi.BlockOverrides,
) (*models.SwapResponse, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, _timeout)
	defer cancel()

	// also settings here
	// creates new arbitrage worker that perform searching for sandwich best parameters to receive
	// maximum profit
	startedAt := time.Now().UTC()
	worker := algo.NewWorker(algo.WorkerSettings{})
	transactor := transact.NewTransactor()

	applyContractCode(stateOverride, request.Contract, request.SimulationCode)

	log.Info(fmt.Sprintf("%s", request.SimulationCode))

	victimExecution, stateDB, header, vmctx, err := prepareSwapInitial(
		ctx,
		b,
		transactor,
		simulation_models.PrepareTemplatesDTO{
			InputPair:         request.InputPair,
			InputPairVersion:  request.InputPairVersion,
			OutputPair:        request.OutputPair,
			OutputPairVersion: request.OutputPairVersion,
			InputToken:        request.InputToken,
			OutputToken:       request.OutputToken,
			Contract:          request.Contract,
		},
		request.Transactions,
		stateOverride,
		blockNrOrHash,
		blockOverrides,
	)
	if err != nil {
		return nil, err
	}

	algoResult := worker.Execute(
		ctx,
		algo.InitialExtraData{
			MaxDepth:     request.MaxDepth,
			MaxBruteTime: request.MaxBruteTime,
			Points:       request.Points,
			SplitParam:   request.SplitParam,
		},
		func(value *big.Int) ([]models.CallManyResponseDTO, *big.Int) {
			return simulateSwaps(
				ctx,
				b,
				value,
				transactor,
				request.Contract,
				stateDB.Copy(), // required option .Copy()!!!!!
				header,
				vmctx,
			)
		})

	return &models.SwapResponse{
		OptimalValue: algoResult.BestValue,
		Profit:       algoResult.BestProfit,
		Reason:       algoResult.Reason,
		Duration:     time.Now().UTC().Sub(startedAt).Milliseconds(),
		Execution:    append(victimExecution, algoResult.Execution...),
	}, nil
}

func simulateSwaps(
	ctx context.Context,
	b ethapi.Backend,
	value *big.Int,
	transactor protocol.Transactor,
	contract common.Address,
	stateDB *state.StateDB,
	header *types.Header,
	vmctx vm.BlockContext,
) ([]models.CallManyResponseDTO, *big.Int) {
	executionFront, _, _, stateDB, _, _, err := applySwaps(
		ctx,
		b,
		&simulation_models.PackFrontDTO{
			Value: value,
		},
		&simulation_models.PackBackDTO{},
		transactor,
		contract,
		stateDB,
		header,
		vmctx,
		0,
	)
	if err != nil {
		return executionFront, big.NewInt(0)
	}

	outputAmountPurchase, err := unpacker.UnpackerObj.ParseOutputAmount(executionFront)
	if err != nil {
		return executionFront, big.NewInt(0)
	}

	log.Info(fmt.Sprintf("outputAmountPurchase: %s", outputAmountPurchase.String()))

	executionBack, _, _, _, _, _, err := applySwaps(
		ctx,
		b,
		&simulation_models.PackFrontDTO{},
		&simulation_models.PackBackDTO{
			Value: outputAmountPurchase,
		},
		transactor,
		contract,
		stateDB,
		header,
		vmctx,
		0,
	)
	if err != nil {
		return executionFront, big.NewInt(0)
	}

	resultExecution := append(executionFront, executionBack...)

	outputAmountSell, err := unpacker.UnpackerObj.ParseOutputAmount(executionFront)
	if err != nil {
		return resultExecution, big.NewInt(0)
	}

	log.Info(fmt.Sprintf("outputAmountSell: %s", outputAmountSell.String()))

	profit := new(big.Int).Sub(outputAmountSell, value)

	return resultExecution, profit
}
