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

	log.Info(fmt.Sprintf("%s", request.Contract))
	applyContractCode(stateOverride, request.Contract, request.SimulationCode)

	victimExecution, stateDB, header, vmctx, err := prepareSwapInitial(
		ctx,
		b,
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
				request.Pairs,
				request.InputToken,
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
	pairs []models.SwapContractParams,
	input common.Address,
	contract common.Address,
	stateDB *state.StateDB,
	header *types.Header,
	vmctx vm.BlockContext,
) ([]models.CallManyResponseDTO, *big.Int) {

	var outputAmount *big.Int
	var execution []models.CallManyResponseDTO

	swapAmount := value
	for _, pair := range pairs {
		var executionTmp []models.CallManyResponseDTO

		token0, err := CallToken0(b, pair.Pair, stateDB, header, vmctx)
		if err != nil {
			return execution, big.NewInt(0)
		}
		token1, err := CallToken1(b, pair.Pair, stateDB, header, vmctx)
		if err != nil {
			return execution, big.NewInt(0)
		}

		var output common.Address

		if input == token0 {
			output = token1
		} else {
			output = token0
		}

		balance, err := CallERC20BalanceOf(b, input, contract, stateDB, header, vmctx)
		if err != nil {
			log.Info(fmt.Sprintf("%s", err.Error()))
			return execution, big.NewInt(0)
		}

		log.Info(fmt.Sprintf("swap %s %s %s %d; inputBalance: %s", pair.Pair.String(), input.String(), output.String(), swapAmount.String(), balance.String()))

		executionTmp, _, _, stateDB, _, _, err = applySwap(
			ctx,
			b,
			&simulation_models.PackFrontDTO{
				Value:    swapAmount,
				Pair:     pair.Pair,
				Input:    input,
				Output:   output,
				Contract: contract,
				PairType: pair.PairVersion,
			},
			transactor,
			contract,
			stateDB,
			header,
			vmctx,
			0,
		)
		execution = append(execution, executionTmp...)
		if err != nil {
			return execution, big.NewInt(0)
		}

		if pair.PairVersion == 2 {
			outputAmount, err = unpacker.UnpackerObj.ParseOutputAmount(executionTmp)
			if err != nil {
				return execution, big.NewInt(0)
			}
		} else {
			outputAmount, err = unpacker.UnpackerObj.ParseOutputAmount3(executionTmp)
			if err != nil {
				return execution, big.NewInt(0)
			}
		}

		swapAmount = outputAmount
		log.Info(fmt.Sprintf("outputAmount: %s", outputAmount.String()))

		input = output
	}

	profit := new(big.Int).Sub(outputAmount, value)

	return execution, profit
}
