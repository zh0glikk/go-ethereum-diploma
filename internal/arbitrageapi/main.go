package arbitrageapi

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/internal/ethapi"
)

type ArbitrageAPI struct {
	b ethapi.Backend
}

func NewArbitrageAPI(b ethapi.Backend) *ArbitrageAPI {
	return &ArbitrageAPI{b}
}

type TraceAPI struct {
	b ethapi.Backend
}

func NewTraceAPI(b ethapi.Backend) *TraceAPI {
	return &TraceAPI{b: b}
}

type BundleAPI struct {
	b     ethapi.Backend
	chain *core.BlockChain
}

func NewBundleAPI(b ethapi.Backend, chain *core.BlockChain) *BundleAPI {
	return &BundleAPI{b, chain}
}
