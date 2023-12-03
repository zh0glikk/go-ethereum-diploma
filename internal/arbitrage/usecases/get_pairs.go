package usecases

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

func GetPairs(
	ctx context.Context,
	b ethapi.Backend,
	request models.GetPairsRequest,
	stateOverride *ethapi.StateOverride,
	blockNrOrHash rpc.BlockNumberOrHash,
	blockOverrides *ethapi.BlockOverrides,
) ([]models.PairInfo, error) {

	_, state, header, vmctx, err := DoCallManyReturningState(ctx, b, nil, blockNrOrHash,
		stateOverride, b.RPCEVMTimeout(), b.RPCGasCap(), blockOverrides)
	if err != nil {
		return nil, err
	}

	var result []models.PairInfo

	pairAddr, err := CallGetPair(ctx, b, uniSwapV2FactoryAddr, request.Token0, request.Token1, state, header, vmctx)
	if err == nil {
		result = append(result, models.PairInfo{
			Pair:    pairAddr,
			Version: 2,
			Factory: "univ2",
		})
	} else {
		log.Info(fmt.Sprintf("%s", err.Error()))
	}

	pairAddr, err = CallGetPair(ctx, b, uniSwapV2FactoryAddr, request.Token0, request.Token1, state, header, vmctx)
	if err == nil {
		result = append(result, models.PairInfo{
			Pair:    pairAddr,
			Version: 2,
			Factory: "sushiv2",
		})
	} else {
		log.Info(fmt.Sprintf("%s", err.Error()))
	}

	for _, fee := range []*big.Int{
		big.NewInt(1000),
		big.NewInt(3000),
		big.NewInt(5000),
		big.NewInt(10000),
	} {
		pairAddr, err = CallGetPool(ctx, b, request.Token0, request.Token1, fee, state, header, vmctx)
		if err == nil {
			result = append(result, models.PairInfo{
				Pair:    pairAddr,
				Version: 2,
				Factory: fmt.Sprintf("univ3_%s", fee.String()),
			})
		} else {
			log.Info(fmt.Sprintf("%s", err.Error()))
		}
	}

	return result, nil
}
