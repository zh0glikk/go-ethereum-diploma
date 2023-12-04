package simulation_wrappers

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/arbitrage/usecases"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"math/big"
)

type Wrapper struct {
	b ethapi.Backend

	contract     common.Address
	maxDepth     int
	maxBruteTime int64
	splitParam   *big.Int
}

func NewWrapper(
	b ethapi.Backend,
	contract common.Address,
	maxDepth int,
	maxBruteTime int64,
	splitParam *big.Int,
) *Wrapper {
	return &Wrapper{
		b:            b,
		contract:     contract,
		maxDepth:     maxDepth,
		maxBruteTime: maxBruteTime,
		splitParam:   splitParam,
	}
}

func (s *Wrapper) Simulate(ctx context.Context, data models.Data2Simulate) (*models.SwapResponse, error) {
	return usecases.ExecuteSwaps(
		ctx,
		s.b,
		models.SwapBundle{
			Pairs:      data.Pairs,
			InputToken: data.InputToken,
			Contract:   data.Contract,
			AlgoCommon: models.AlgoCommon{
				MaxDepth:     s.maxDepth,
				MaxBruteTime: s.maxBruteTime,
				SplitParam:   s.splitParam,
				Points:       data.Points,
			},
			Transactions:   data.Transactions,
			SimulationCode: &byteCode,
		},
		&ethapi.StateOverride{},
		data.BlockNumberOrHash,
		data.BlockOverrides,
	)
}
