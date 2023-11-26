package arbitrageapi

import (
	"context"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/arbitrage/usecases"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
)

func (s *ArbitrageAPI) CallMany(ctx context.Context, args []ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *ethapi.StateOverride, blockOverrides *ethapi.BlockOverrides) ([]models.CallManyResponseDTO, error) {
	return usecases.DoCallMany(ctx, s.b, args, blockNrOrHash, overrides, s.b.RPCEVMTimeout(), s.b.RPCGasCap(), blockOverrides)
}

// for decoding
func (s *ArbitrageAPI) CallMany2(ctx context.Context, args []ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *ethapi.StateOverride, blockOverrides *ethapi.BlockOverrides) ([]models.CallManyResponseDTO, error) {
	return usecases.DoCallManyWithErrorDecoding(ctx, s.b, args, blockNrOrHash, overrides, s.b.RPCEVMTimeout(), s.b.RPCGasCap(), blockOverrides)
}

func (s *ArbitrageAPI) SimpleBatchCallMany(ctx context.Context, args [][]ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *ethapi.StateOverride, blockOverrides *ethapi.BlockOverrides) ([][]models.CallManyResponseDTO, error) {
	return usecases.SimpleBatchCallMany(ctx, s.b, args, blockNrOrHash, overrides, s.b.RPCEVMTimeout(), s.b.RPCGasCap(), blockOverrides)
}

func (s *ArbitrageAPI) BatchCallMany(ctx context.Context, args [][]ethapi.TransactionArgs, blockNrOrHash []rpc.BlockNumberOrHash, overrides *[]*ethapi.StateOverride, blockOverrides *[]*ethapi.BlockOverrides) ([][]models.CallManyResponseDTO, error) {
	return usecases.BatchCallMany(ctx, s.b, args, blockNrOrHash, overrides, blockOverrides)
}
