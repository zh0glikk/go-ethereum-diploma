package simulation_models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"math/big"
)

type PrepareTemplatesDTO struct {
	InputPair         common.Address
	InputPairVersion  int
	OutputPair        common.Address
	OutputPairVersion int

	InputToken  common.Address
	OutputToken common.Address
	Contract    common.Address
}

type PackFrontDTO struct {
	Value *big.Int
}

type PackBackDTO struct {
	Value *big.Int
}

type CalculateProfit struct {
	Execution []models.CallManyResponseDTO
	Value     *big.Int
	ActionBuy bool
}
