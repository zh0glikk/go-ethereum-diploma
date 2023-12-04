package mempool

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"sync"
)

type InfuraMemPool struct {
	WsUrl string
}

func NewInfuraMemPool(WsUrl string) *InfuraMemPool {
	return &InfuraMemPool{WsUrl: WsUrl}
}

func (i *InfuraMemPool) Listen(ctx context.Context, ch chan types.Transaction, wg *sync.WaitGroup) error {
	streamer, err := NewStreamer(i.WsUrl, ch)
	if err != nil {
		return err
	}

	log.Info("infura mempool started")

	_ = streamer.Stream(ctx, wg)
	return err
}
