package mempool

import (
	"context"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"sync"
)

type InternalMemPool struct {
	b ethapi.Backend
}

func NewInternalMemPool(b ethapi.Backend) *InternalMemPool {
	return &InternalMemPool{b: b}
}

func (i *InternalMemPool) Listen(ctx context.Context, ch chan types.Transaction, wg *sync.WaitGroup) error {
	txsEvent := make(chan core.NewTxsEvent)
	sub := i.b.SubscribeNewTxsEvent(txsEvent)

	log.Info("internal mempool started")

	for {
		select {
		case txs := <-txsEvent:
			for _, tx := range txs.Txs {
				ch <- *tx
			}
		case <-ctx.Done():
			sub.Unsubscribe()
			wg.Done()
			log.Info("internal mempool listener stopped")
			return nil
		}
	}
}
