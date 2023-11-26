package arbitrageapi

import (
	"context"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/internal/arbitrage/usecases"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
)

func (api *TraceAPI) Call(ctx context.Context, args ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, config *tracers.TraceCallConfig) (interface{}, error) {
	return usecases.TraceCall(ctx, api.b, args, blockNrOrHash, config)
}

func (api *TraceAPI) CallMany(ctx context.Context, txs []ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, config *tracers.TraceCallConfig) (interface{}, error) {
	return usecases.TraceCallMany(ctx, api.b, txs, blockNrOrHash, config)
}
