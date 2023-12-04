package bot

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/arbitrage/bot/runners/mempool"
	"github.com/ethereum/go-ethereum/internal/arbitrage/bot/storage"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/arbitrage/usecases"
	"github.com/ethereum/go-ethereum/internal/arbitrage/utils"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"sync"
	"time"
)

type Service struct {
	cfg *ServiceConfig

	// chans here for conveyor
	blockIncludedQueue    chan *types.Block
	waitingTraceTxsQueue  chan models.WaitingTraceTxsJob
	simulationNotifyQueue chan models.SimulationNotifyJob
	simulationJobs        chan models.SimulationJob

	pendingTxsStorage *storage.PendingTxsStorage
	txsStorage        *storage.TxsStorage
	syncingCh         chan bool
	runCh             chan bool
}

func NewService(cfg *ServiceConfig) *Service {
	return &Service{
		cfg:                   cfg,
		waitingTraceTxsQueue:  make(chan models.WaitingTraceTxsJob, 10),
		blockIncludedQueue:    make(chan *types.Block, 1),
		simulationNotifyQueue: make(chan models.SimulationNotifyJob, 10),
		simulationJobs:        make(chan models.SimulationJob, 10),
		pendingTxsStorage:     storage.NewPendingTxsStorage(),
		syncingCh:             make(chan bool),
		runCh:                 make(chan bool),
	}
}

func (s *Service) Run() {
	s.cfg.wg.Add(1)
	go s.isSyncing()

	for {
		select {
		case <-s.runCh:
			s.cfg.wg.Add(7)

			go s.memPoolListener()
			go s.blockListener()

			go s.tracingWriter()
			go s.tracer()

			go s.pendingTxsStorageCleaner()

			go s.simulationWriter()
			go s.simulator()

		case <-s.syncingCh:
			s.cfg.wg.Add(1)
			go s.infuraCompare()
		}
	}

}

func (s *Service) memPoolListener() {
	streamerCtx, cancelF := context.WithCancel(context.Background())
	ch := make(chan types.Transaction, 100)

	log.Info("MemPoolListener started")

	wg := &sync.WaitGroup{}
	for _, memPool := range s.cfg.memPools {
		wg.Add(1)
		go func(memPool mempool.MemPool, streamerCtx context.Context, ch chan types.Transaction, wg *sync.WaitGroup) {
			err := memPool.Listen(streamerCtx, ch, wg)
			if err != nil {
				log.Error(fmt.Sprintf("failed to start mempool listener: %v", err))
				wg.Done()
				return
			}
		}(memPool, streamerCtx, ch, wg)
	}

	for {
		select {
		case memPoolTx := <-ch:
			from, err := types.Sender(types.LatestSignerForChainID(memPoolTx.ChainId()), &memPoolTx)
			if err != nil {
				log.Error("err", "failed to get from address with reason: %v", err)
				continue
			}

			s.pendingTxsStorage.Add(from, memPoolTx)

		case <-s.cfg.ctx.Done():
			log.Info("stopping MemPoolListener runner")

			cancelF()

			wg.Wait()

			s.cfg.wg.Done()
			log.Info("stopped MemPoolListener runner")
			return
		}
	}
}

func (s *Service) blockListener() {
	wg := &sync.WaitGroup{}

	blockListenerCtx, cancelF := context.WithCancel(context.Background())
	ch := s.cfg.blockListener.Listen(blockListenerCtx, wg)

	log.Info("block listener started")

	for {
		select {
		case block := <-ch:
			s.txsStorage = storage.NewTxsStorage()

			log.Info(fmt.Sprintf("Latest block number: %v", block.Number()))

			for _, tx := range block.Transactions() {
				from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
				if err != nil {
					log.Error("err", "failed to get from address with reason: %v", err)
					continue
				}

				s.pendingTxsStorage.Delete(from, tx.Hash())
			}

			log.Info(fmt.Sprintf("COUNT: %d, txs: %d", s.pendingTxsStorage.CountOfGroups(), s.pendingTxsStorage.CountOfTxs()))

			s.blockIncludedQueue <- block

		case <-s.cfg.ctx.Done():
			log.Info("stopping block listener runner")

			cancelF()

			wg.Wait()

			s.cfg.wg.Done()
			log.Info("stopped block listener runner")
			return
		}
	}
}

func (s *Service) tracingWriter() {
	log.Info("tracingWriter started")

	for {
		select {
		case block := <-s.blockIncludedQueue:

			header, err := usecases.GetHeaderByNumber(s.cfg.ctx, s.cfg.b, rpc.BlockNumber(block.Number().Int64()))
			if err != nil {
				log.Info(fmt.Sprintf("failed to GetHeaderByNumber due to err: %s", err.Error()))
				continue
			}
			nextBaseFee, ok := header["nextBaseFee"]
			if !ok || nextBaseFee == nil {
				log.Info(fmt.Sprintf("nextBaseFee not found or nil: %v", header["nextBaseFee"]))
				continue
			}
			nextBaseFeeBig := nextBaseFee.(*hexutil.Big)

			pendingTxsGroups := s.pendingTxsStorage.List()

			for _, pendingTxsGroup := range pendingTxsGroups {
				var args []ethapi.TransactionArgs
				for _, tx := range pendingTxsGroup.Transactions {
					args = append(args, ethapi.TransactionArgs{
						From:                 utils.Ptr(pendingTxsGroup.From),
						To:                   tx.To(),
						Gas:                  utils.Ptr((hexutil.Uint64)(tx.Gas())),
						MaxFeePerGas:         (*hexutil.Big)(tx.GasFeeCap()),
						MaxPriorityFeePerGas: (*hexutil.Big)(tx.GasTipCap()),
						Value:                (*hexutil.Big)(tx.Value()),
						Nonce:                utils.Ptr((hexutil.Uint64)(tx.Nonce())),
						Data:                 utils.Ptr(hexutil.Bytes(tx.Data())),
						ChainID:              (*hexutil.Big)(tx.ChainId()),
					})
				}

				s.waitingTraceTxsQueue <- models.WaitingTraceTxsJob{
					From:   pendingTxsGroup.From,
					Args:   args,
					RawTxs: pendingTxsGroup.Transactions,
					Swaps:  nil,
					BlockNumberOrHash: rpc.BlockNumberOrHash{
						BlockNumber: utils.Ptr(rpc.BlockNumber(block.Number().Int64())),
					},
					BlockOverrides: &ethapi.BlockOverrides{
						Number:  (*hexutil.Big)(new(big.Int).Add(block.Number(), big.NewInt(1))),
						BaseFee: nextBaseFeeBig,
					},
					CreatedAt: time.Now().UTC(),
				}
			}

		case <-s.cfg.ctx.Done():
			log.Info("stopping tracingWriter runner")
			s.cfg.wg.Done()
			log.Info("stopped tracingWriter runner")
			return
		}
	}
}

func (s *Service) tracer() {
	wg := &sync.WaitGroup{}

	tracerCtx, cancelF := context.WithCancel(context.Background())
	resultsChan := s.cfg.tracer.Trace(tracerCtx, s.waitingTraceTxsQueue, wg)
	log.Info("tracing runner started")

	for {
		select {
		case job := <-resultsChan:
			swapFound := false
			for _, swaps := range job.Swaps {
				for range swaps {
					swapFound = true
				}
			}

			if !swapFound {
				for _, tx := range job.RawTxs {
					s.pendingTxsStorage.Delete(job.From, tx.Hash())
				}
				continue
			}

			uniquePairs := make(map[common.Address]bool)
			for _, swaps := range job.Swaps {
				for _, swap := range swaps {
					if _, ok := uniquePairs[swap.Pair]; !ok {
						uniquePairs[swap.Pair] = true
					}

					s.txsStorage.AddOrUpdate(swap.Pair, swap.Input, swap.Output, swap.Type, job.RawTxs, job.Args)
				}
			}

			var pairs []common.Address
			for pair := range uniquePairs {
				pairs = append(pairs, pair)
			}

			s.simulationNotifyQueue <- models.SimulationNotifyJob{
				BlockNumberOrHash: job.BlockNumberOrHash,
				BlockOverrides:    job.BlockOverrides,
				Pairs:             pairs,
			}
		case <-s.cfg.ctx.Done():
			log.Info("stopping tracing runner")

			cancelF()

			// wait until all tracing workers stop
			wg.Wait()

			s.cfg.wg.Done()
			log.Info("stopped tracing runner")
			return
		}
	}
}

func (s *Service) pendingTxsStorageCleaner() {
	log.Info("pendingTxsStorageCleaner started")

	ticker := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-ticker.C:

			count := s.pendingTxsStorage.ClearExpired()
			log.Info(fmt.Sprintf("cleared expired entries: %d", count))

		case <-s.cfg.ctx.Done():
			s.cfg.wg.Done()
			log.Info("pendingTxsStorageCleaner runner")
			return
		}
	}
}

func (s *Service) simulationWriter() {
	log.Info("simulationWriter started")

	for {
		select {
		case job := <-s.simulationNotifyQueue:
			// for debug
			tmp := ""
			for _, pair := range job.Pairs {
				tmp += pair.String() + "\n"
			}
			log.Info(fmt.Sprintf("pairs: %s", tmp))

			// simulationJobs := []models.SimulationJob{}
			for _, txsModel := range s.txsStorage.ListFiltered(job.Pairs) {
				log.Info(fmt.Sprintf("selected: \npair: %s, token0: %s, token1: %s\nin: %s\nout: %s\n", txsModel.Pair.String(), txsModel.Token0.String(), txsModel.Token1.String(), txsModel.TxsIn.Hashes(), txsModel.TxsOut.Hashes()))

				s.simulationJobs <- models.SimulationJob{
					BlockNumberOrHash: job.BlockNumberOrHash,
					BlockOverrides:    job.BlockOverrides,
					Pair:              txsModel.Pair,
					PairType:          txsModel.PairType,
					Token0:            txsModel.Token0,
					Token1:            txsModel.Token1,
					TxsIn:             txsModel.TxsIn,
					TxsOut:            txsModel.TxsOut,
				}

			}
			log.Info(fmt.Sprintf("txsStorage count: %d", len(s.txsStorage.List())))

		case <-s.cfg.ctx.Done():
			s.cfg.wg.Done()
			log.Info("simulationWriter stopped")
			return
		}
	}
}

func (s *Service) simulator() {
	wg := &sync.WaitGroup{}

	ctx, cancelF := context.WithCancel(context.Background())
	resultsChan := s.cfg.simulator.Run(ctx, s.simulationJobs, wg)
	log.Info("simulator started")

	for {
		select {
		case job := <-resultsChan:
			_ = job
		// TODO: implement

		case <-s.cfg.ctx.Done():
			log.Info("stopping simulator")

			cancelF()

			// wait until all simulator workers stop
			wg.Wait()

			s.cfg.wg.Done()
			log.Info("stopped simulator")
			return
		}
	}
}

func (s *Service) isSyncing() {
	defer s.cfg.wg.Done()
	for {
		select {
		case <-s.cfg.ctx.Done():
			log.Info("stopping syncing")
			return
		default:
			progress := s.cfg.b.SyncProgress()

			if progress.CurrentBlock <= progress.HighestBlock {
				log.Info(fmt.Sprintf("Sando bot not started because node is syncing, current block : %d highest block : %d", progress.CurrentBlock, progress.HighestBlock))
				time.Sleep(1 * time.Second)
				continue
			}

			log.Info("Syncing completed!!!")
			s.syncingCh <- true
			return
		}

	}
}

func (s *Service) infuraCompare() {
	defer s.cfg.wg.Done()
	for {
		select {
		case <-s.cfg.ctx.Done():
			log.Info("stopping compare")
			return

		default:
			infuraBlockNumber, err := s.cfg.infuraEthClient.BlockNumber(s.cfg.ctx)
			if err != nil {
				log.Error(fmt.Sprintf("failed to get infura block number due to err :%s", err.Error()))
				return
			}
			progress := s.cfg.b.SyncProgress()

			if infuraBlockNumber-progress.CurrentBlock >= 5 {
				log.Info(fmt.Sprintf("Sando bot not started because node is syncing, node block: %d, infura block: %d", progress.CurrentBlock, infuraBlockNumber))
				time.Sleep(1 * time.Second)
				continue
			}
			log.Info("Syncing completed!!!")
			s.runCh <- true
			return
		}

	}
}
