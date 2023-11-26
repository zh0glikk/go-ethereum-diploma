package arbitrageapi

import (
	"context"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/arbitrage/usecases"
	"github.com/ethereum/go-ethereum/rpc"
)

func (api *TraceAPI) TrackSwap(
	ctx context.Context,
	request models.TrackSwapBundle,
	blockNrOrHash rpc.BlockNumberOrHash,
	config *tracers.TraceCallConfig,
) (*models.TrackSwapResponse, error) {
	return usecases.TrackSwap(ctx, api.b, request, blockNrOrHash, config)
}
