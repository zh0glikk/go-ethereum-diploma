package transact

import (
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/implementations"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/protocol"
)

func NewTransactor() protocol.Transactor {
	return implementations.NewSwapTransactor()
}
