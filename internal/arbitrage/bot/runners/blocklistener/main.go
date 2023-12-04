package blocklistener

import (
	"context"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"sync"
)

type BlockListener struct {
	chain *core.BlockChain

	lastestBlocks chan *types.Block
}

func NewBlockListener(chain *core.BlockChain) *BlockListener {
	return &BlockListener{lastestBlocks: make(chan *types.Block), chain: chain}
}

func (b *BlockListener) Listen(ctx context.Context, wg *sync.WaitGroup) chan *types.Block {
	wg.Add(1)
	go b.listenLatestBlocks(ctx, wg)

	return b.lastestBlocks
}

func (b *BlockListener) listenLatestBlocks(ctx context.Context, wg *sync.WaitGroup) {
	headEventChan := make(chan core.ChainHeadEvent)
	sub := b.chain.SubscribeChainHeadEvent(headEventChan)

	for {
		select {
		case headEvent := <-headEventChan:
			b.lastestBlocks <- headEvent.Block
		case <-ctx.Done():
			sub.Unsubscribe()
			log.Info("BlockListener stopped")
			wg.Done()
			return
		}
	}
}
