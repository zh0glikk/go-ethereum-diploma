package usecases

import (
	"context"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
	"time"
)

func FilterSwaps(
	ctx context.Context,
	b ethapi.Backend,
	request [][]models.SwapInfo,
	stateOverride *ethapi.StateOverride,
	blockNrOrHash rpc.BlockNumberOrHash,
	blockOverrides *ethapi.BlockOverrides,
) (*models.TrackSwapResponse, error) {
	startedAt := time.Now().UTC()

	_, state, header, vmctx, err := DoCallManyReturningState(ctx, b, nil, blockNrOrHash,
		stateOverride, b.RPCEVMTimeout(), b.RPCGasCap(), blockOverrides)
	if err != nil {
		return nil, err
	}

	responses, err := FilterSwapsOnState(ctx, b, request, state, header, vmctx)
	if err != nil {
		return nil, err
	}

	return &models.TrackSwapResponse{
		Swaps:    responses,
		Duration: time.Now().UTC().Sub(startedAt).Milliseconds(),
	}, nil
}

func FilterSwapsOnState(
	ctx context.Context,
	b ethapi.Backend,
	requests [][]models.SwapInfo,
	stateDB *state.StateDB,
	header *types.Header,
	vmctx vm.BlockContext,
) ([][]models.SwapInfo, error) {
	responses := [][]models.SwapInfo{}
	for _, request := range requests {
		resp := []models.SwapInfo{}
		for _, r := range request {
			factoryAddr, err := CallFactory(ctx, b, r.Pair, stateDB, header, vmctx)
			if err != nil {
				continue
			}

			if factoryAddr == uniSwapV2FactoryAddr {
				pair, err := CallGetPair(ctx, b, uniSwapV2FactoryAddr, r.Input, r.Output, stateDB, header, vmctx)
				if err != nil {
					continue
				}
				if pair == r.Pair {
					resp = append(resp, r)
				}
			}
			if factoryAddr == sushiSwapV2FactoryAddr {
				pair, err := CallGetPair(ctx, b, sushiSwapV2FactoryAddr, r.Input, r.Output, stateDB, header, vmctx)
				if err != nil {
					continue
				}
				if pair == r.Pair {
					resp = append(resp, r)
				}
			}
			if factoryAddr == uniSwapV3FactoryAddr {
				fee, err := CallPairFee(ctx, b, r.Pair, stateDB, header, vmctx)
				if err != nil {
					continue
				}

				pair, err := CallGetPool(ctx, b, r.Input, r.Output, fee, stateDB, header, vmctx)
				if err != nil {
					continue
				}

				if pair == r.Pair {
					resp = append(resp, r)
				}
			}
		}

		responses = append(responses, resp)
	}
	
	return responses, nil
}
