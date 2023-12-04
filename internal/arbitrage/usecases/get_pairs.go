package usecases

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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
	if err != nil {
		log.Info(fmt.Sprintf("%s", err.Error()))
	} else {
		if pairAddr != common.HexToAddress("0x0") {
			result = append(result, models.PairInfo{
				Pair:    pairAddr,
				Version: 2,
				DEX:     "univ2",
			})
		}
	}

	pairAddr, err = CallGetPair(ctx, b, uniSwapV2FactoryAddr, request.Token0, request.Token1, state, header, vmctx)
	if err != nil {
		log.Info(fmt.Sprintf("%s", err.Error()))
	} else {
		if pairAddr != common.HexToAddress("0x0") {
			result = append(result, models.PairInfo{
				Pair:    pairAddr,
				Version: 2,
				DEX:     "sushiv2",
			})
		}
	}

	for _, fee := range []*big.Int{
		big.NewInt(100),
		big.NewInt(500),
		big.NewInt(3000),
		big.NewInt(10000),
	} {
		pairAddr, err = CallGetPool(ctx, b, request.Token0, request.Token1, fee, state, header, vmctx)
		if err != nil {
			log.Info(fmt.Sprintf("%s", err.Error()))
		} else {
			if pairAddr != common.HexToAddress("0x0") {
				result = append(result, models.PairInfo{
					Pair:    pairAddr,
					Version: 3,
					DEX:     fmt.Sprintf("univ3_%s", fee.String()),
				})
			}
		}
	}

	return result, nil
}
