package blocklistener

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core"
	"sync"
	"testing"
	"time"
)

func TestBlockListener(t *testing.T) {
	fmt.Println("start block listener test")
	cli := NewBlockListener(&core.BlockChain{})
	wg := &sync.WaitGroup{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	ch := cli.Listen(ctx, wg)

	for {
		select {
		case v := <-ch:
			fmt.Println(v.Number())
			time.Sleep(time.Second * 5)
		case <-ctx.Done():
			wg.Wait()
			fmt.Println("block listener stopped")
			return
		}
	}
}
