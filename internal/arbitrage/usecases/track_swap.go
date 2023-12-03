package usecases

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/unpacker"
	"github.com/ethereum/go-ethereum/internal/arbitrage/utils"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"time"
)

func TrackSwap(
	ctx context.Context,
	b ethapi.Backend,
	request models.TrackSwapBundle,
	blockNrOrHash rpc.BlockNumberOrHash,
	config *tracers.TraceCallConfig,
) (*models.TrackSwapResponse, error) {
	startedAt := time.Now().UTC()

	traceResponse, err := TraceCallMany(ctx, b, request.Transactions, blockNrOrHash, config)
	if err != nil {
		return nil, err
	}

	traces, err := convertTraceCallManyResult(traceResponse)
	if err != nil {
		return nil, err
	}

	swaps := [][]models.SwapInfo{}
	for _, txTraces := range traces {
		swaps = append(swaps, parseV2V3(txTraces))
	}

	return &models.TrackSwapResponse{
		Swaps:    swaps,
		Duration: time.Now().UTC().Sub(startedAt).Milliseconds(),
	}, nil
}

type transferWrapper struct {
	From   common.Address
	To     common.Address
	Amount *big.Int
	Token  common.Address
}

func parseV2V3(traces []models.TransactionTrace) []models.SwapInfo {
	swaps := []models.SwapInfo{}

	transfersByFrom := make(map[common.Address][]transferWrapper)
	transfersByTo := make(map[common.Address][]transferWrapper)

	for _, trace := range traces {
		bb, err := hexutil.Decode(trace.Action.Input)
		if err != nil {
			continue
		}
		if unpacker.UnpackerObj.IsSwapV2(bb) {
			swaps = append(swaps, models.SwapInfo{Type: "v2", Pair: trace.Action.To})
		}
		if unpacker.UnpackerObj.IsSwapV3(bb) {
			swaps = append(swaps, models.SwapInfo{Type: "v3", Pair: trace.Action.To})
		}
		if unpacker.UnpackerObj.IsTransfer(bb) {
			to, amount := unpacker.UnpackerObj.UnpackTransfer(bb)
			transfersByFrom[trace.Action.From] = append(
				transfersByFrom[trace.Action.From], transferWrapper{From: trace.Action.From, To: to, Amount: amount, Token: trace.Action.To})
			transfersByTo[to] = append(
				transfersByTo[to], transferWrapper{From: trace.Action.From, To: to, Amount: amount, Token: trace.Action.To})
		}
		if unpacker.UnpackerObj.IsTransferFrom(bb) {
			from, to, amount := unpacker.UnpackerObj.UnpackTransferFrom(bb)
			transfersByFrom[trace.Action.From] = append(
				transfersByFrom[trace.Action.From], transferWrapper{From: from, To: to, Amount: amount, Token: trace.Action.To})
			transfersByTo[to] = append(
				transfersByTo[to], transferWrapper{From: from, To: to, Amount: amount, Token: trace.Action.To})
		}
	}

	if len(swaps) == 0 {
		return swaps
	}

	for i := range swaps {
		if len(transfersByTo[swaps[i].Pair]) == 0 ||
			len(transfersByFrom[swaps[i].Pair]) == 0 {
			continue
		}

		swaps[i].InputAmount = transfersByTo[swaps[i].Pair][0].Amount

		swaps[i].Input = transfersByTo[swaps[i].Pair][0].Token
		swaps[i].Output = transfersByFrom[swaps[i].Pair][0].Token

		utils.RemoveFromArray(transfersByTo[swaps[i].Pair], 0)
		utils.RemoveFromArray(transfersByFrom[swaps[i].Pair], 0)
	}

	return swaps
}
