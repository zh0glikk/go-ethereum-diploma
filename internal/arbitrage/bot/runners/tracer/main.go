package tracer

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/arbitrage/usecases"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"strings"
	"sync"
)

const workersCount = 10

type Tracer struct {
	b        ethapi.Backend
	resultCh chan models.WaitingTraceTxsJob
}

func NewTracer(b ethapi.Backend) *Tracer {
	return &Tracer{b: b, resultCh: make(chan models.WaitingTraceTxsJob, 10)}
}

func (t *Tracer) Trace(ctx context.Context, waitingTraceTxsQueue chan models.WaitingTraceTxsJob, wg *sync.WaitGroup) chan models.WaitingTraceTxsJob {
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go t.worker(ctx, waitingTraceTxsQueue, wg)
	}

	return t.resultCh
}

func (t *Tracer) worker(ctx context.Context, waitingTraceTxsQueue chan models.WaitingTraceTxsJob, wg *sync.WaitGroup) {
	log.Info("tracing worker started")
	for {
		select {
		case tx := <-waitingTraceTxsQueue:
			err := t.processTransaction(ctx, tx)
			if err != nil {
				log.Error("processTransaction", fmt.Sprintf("tracing worker err: %v", err))
			}
		case <-ctx.Done():
			wg.Done()
			log.Info("tracing worker stopped")
			return
		}
	}
}

func (t *Tracer) processTransaction(ctx context.Context, job models.WaitingTraceTxsJob) error {
	count, err := usecases.GetTransactionsCount(
		ctx,
		t.b,
		job.From,
		job.BlockNumberOrHash,
	)
	if err != nil {
		return err
	}

	for _, tx := range job.Args {
		if count != *(*uint64)(tx.Nonce) {
			return nil
		}
		count += 1
	}

	response, err := usecases.TrackSwap(
		ctx,
		t.b,
		models.TrackSwapBundle{
			Transactions: job.Args,
		},
		job.BlockNumberOrHash,
		&tracers.TraceCallConfig{
			TraceConfig:       tracers.TraceConfig{},
			StateOverrides:    nil,
			BlockOverrides:    job.BlockOverrides,
			NestedTraceOutput: false,
		})
	if err != nil {
		if strings.Contains(err.Error(), "max fee per gas less than block base fee") ||
			strings.Contains(err.Error(), "insufficient funds for gas * price + value") {
			return nil
		}
		return errors.Wrap(err, "track swap failed")
	}
	if len(response.Swaps) == 0 {
		log.Info(fmt.Sprintf("len(response.Swaps) == 0 case for tx %s", job.RawTxs[0].Hash()))
		return nil
	}

	response, err = usecases.FilterSwaps(
		ctx,
		t.b,
		response.Swaps,
		nil,
		job.BlockNumberOrHash,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "filter swap failed")
	}

	job.Swaps = response.Swaps

	t.resultCh <- job

	return nil
}
