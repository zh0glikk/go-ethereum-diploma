package models

import "github.com/ethereum/go-ethereum/common"

type GetPairsRequest struct {
	Token0 common.Address `json:"token0"`
	Token1 common.Address `json:"token1"`
}

type PairInfo struct {
	Pair    common.Address `json:"pair"`
	Version int            `json:"pair_version"`
	DEX     string         `json:"dex"` // for debug
}
