package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"math/big"
)

type SwapBundle struct {
	SwapContractParams
	AlgoCommon

	Transactions   []ethapi.TransactionArgs `json:"transactions"`
	SimulationCode *hexutil.Bytes           `json:"simulationCode"`
}

type SwapContractParams struct {
	InputPair         common.Address `json:"input_pair"`
	InputPairVersion  int            `json:"input_pair_version"`
	OutputPair        common.Address `json:"output_pair"`
	OutputPairVersion int            `json:"output_pair_version"`

	InputToken  common.Address `json:"input_token"`
	OutputToken common.Address `json:"output_token"`
	Contract    common.Address `json:"contract"`
}

type AlgoCommon struct {
	MaxDepth          int        `json:"maxDepth"`
	MaxBruteTime      int64      `json:"maxBruteTime"`
	SplitParam        *big.Int   `json:"splitParam"`
	Points            []*big.Int `json:"points"`
	InitialSplitParam *big.Int   `json:"initialSplitParam"` // used for first iteration points formation
}

type SwapResponse struct {
	OptimalValue *big.Int              `json:"optimal_amount_dbg"`
	Profit       *big.Int              `json:"profit_dbg"`
	Reason       string                `json:"reason"`
	Duration     int64                 `json:"duration_dbg"`
	Execution    []CallManyResponseDTO `json:"execution"`
	// GasInfo         map[string]*big.Int   `json:"gas_info,omitempty"`
	// Optimized       *bool                 `json:"optimized,omitempty"`
}
