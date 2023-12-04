package mempool

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"sync"
	"time"
)

type Streamer struct {
	wsCli *rpc.Client

	ch chan types.Transaction
}

type MemPool interface {
	Listen(ctx context.Context, ch chan types.Transaction, wg *sync.WaitGroup) error
}

func NewStreamer(wsUrl string, ch chan types.Transaction) (*Streamer, error) {
	wsCli, err := rpc.Dial(wsUrl)
	if err != nil {
		return nil, err
	}

	return &Streamer{
		wsCli: wsCli,
		ch:    ch,
	}, nil
}

func (s *Streamer) Stream(ctx context.Context, wg *sync.WaitGroup) chan types.Transaction {
	go s.safeSubscribe(ctx, wg)

	return s.ch
}

func (s *Streamer) safeSubscribe(ctx context.Context, wg *sync.WaitGroup) {
	sub, err := s.wsCli.Subscribe(ctx, "eth", s.ch, "newPendingTransactions", true)
	if err != nil {
		log.Info("failed to subscribe: %s", err.Error())
		time.Sleep(time.Second * 5)
		s.safeSubscribe(ctx, wg)
		return
	}

	for {
		select {
		case err = <-sub.Err():
			log.Info("mempool wss sub failed with reason: %s", err.Error())
			s.safeSubscribe(ctx, wg)
			return
		case <-ctx.Done():
			sub.Unsubscribe()
			wg.Done()
			log.Info("mempool listener stopped")
			return
		}
	}
}
