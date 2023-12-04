package mempool

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"sync"
	"testing"
	"time"
)

func TestStreamer(t *testing.T) {
	chTx := make(chan types.Transaction)
	cli, err := NewStreamer("", chTx)
	if err != nil {
		fmt.Println(err)
		return
	}
	wg := &sync.WaitGroup{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	ch := cli.Stream(ctx, wg)

	for {
		select {
		case v := <-ch:
			fmt.Println(v.Hash())
		case <-ctx.Done():
			fmt.Println("main stopped")
			return
		}
	}
}
