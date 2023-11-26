package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"math/big"
)

type TrackSwapBundle struct {
	Transactions []ethapi.TransactionArgs `json:"transactions"`
}

type TrackSwapResponse struct {
	Swaps    [][]SwapInfo `json:"swaps"`
	Duration int64        `json:"duration"`
}

type TrackSwapsBlockResponse struct {
	Swaps    []SwapBlockInfo `json:"blockSwaps"`
	Duration int64           `json:"duration"`
}

type SwapBlockInfo struct {
	Swaps  []SwapInfo  `json:"swaps"`
	TxHash common.Hash `json:"txHash"`
}

type SwapInfo struct {
	Type        string         `json:"type"`
	Pair        common.Address `json:"pair"`
	Input       common.Address `json:"input"`
	Output      common.Address `json:"output"`
	InputAmount *big.Int       `json:"inputAmount"`
}
