package simulator

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/internal/arbitrage/bot/runners/simulator/simulation_wrappers"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/arbitrage/usecases"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"sync"
)

const workersCount = 1

type Simulator struct {
	b ethapi.Backend

	contractAddr             common.Address
	initialSplitParam        *big.Int
	minTokenBalance          *big.Int
	additionalTvlTrackTokens []common.Address
	wrapper                  *simulation_wrappers.Wrapper

	resultCh chan models.SimulationJobResponse
}

func NewSimulator(
	b ethapi.Backend,
	contract common.Address,
	maxDepth int,
	maxBruteTime int64,
	splitParam *big.Int,
	initialSplitParam *big.Int,
	minTokenBalance *big.Int,
	additionalTvlTrackTokens []common.Address,
) *Simulator {
	return &Simulator{
		b:                        b,
		contractAddr:             contract,
		initialSplitParam:        initialSplitParam,
		minTokenBalance:          minTokenBalance,
		additionalTvlTrackTokens: additionalTvlTrackTokens,
		resultCh:                 make(chan models.SimulationJobResponse, 10),
		wrapper:                  simulation_wrappers.NewWrapper(b, contract, maxDepth, maxBruteTime, splitParam),
	}
}

func (t *Simulator) Run(ctx context.Context, simulationJobs chan models.SimulationJob, wg *sync.WaitGroup) chan models.SimulationJobResponse {
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go t.worker(ctx, simulationJobs, wg)
	}

	return t.resultCh
}

func (t *Simulator) worker(ctx context.Context, simulationJobs chan models.SimulationJob, wg *sync.WaitGroup) {
	log.Info("Simulator worker started")
	for {
		select {
		case job := <-simulationJobs:
			result, err := t.processTransaction(ctx, job)
			if err != nil {
				log.Error("%s", err.Error())
				continue
			}

			t.resultCh <- *result
		case <-ctx.Done():
			log.Info("Simulator worker stopped")
			wg.Done()
			return
		}
	}
}

func (t *Simulator) processTransaction(ctx context.Context, job models.SimulationJob) (*models.SimulationJobResponse, error) {
	pairs, err := usecases.GetPairs(ctx, t.b, models.GetPairsRequest{
		Token0: job.Token0,
		Token1: job.Token1,
	}, &ethapi.StateOverride{}, job.BlockNumberOrHash, job.BlockOverrides)
	if err != nil {
		return nil, err
	}

	var simulationData []SimulationData

	if usecases.IsWeth(job.Token0) {
		if len(job.TxsIn.Transactions) > 0 {

			for _, pair := range pairs {
				if pair.Pair == job.Pair {
					continue
				}

				simulationData = append(simulationData, SimulationData{
					Pairs: []models.SwapContractParams{
						{
							Pair:        pair.Pair,
							PairVersion: pair.Version,
						},
						{
							Pair: job.Pair,
							PairVersion: func() int {
								if job.PairType == "v2" {
									return 2
								} else {
									return 3
								}
							}(),
						},
					},
					Txs: job.TxsIn,
				},
				)
			}

			// so it was a purchase
			//  we need to buy on another pair and sell here
		} else {
			for _, pair := range pairs {
				if pair.Pair == job.Pair {
					continue
				}

				simulationData = append(simulationData, SimulationData{

					Pairs: []models.SwapContractParams{
						{
							Pair: job.Pair,
							PairVersion: func() int {
								if job.PairType == "v2" {
									return 2
								} else {
									return 3
								}
							}(),
						},
						{
							Pair:        pair.Pair,
							PairVersion: pair.Version,
						},
					},
					Txs: job.TxsOut,
				},
				)
			}
			// its sell
			// we need to buy here and sell on another pair
		}
	} else if usecases.IsWeth(job.Token1) {
		if len(job.TxsOut.Transactions) > 0 {
			for _, pair := range pairs {
				if pair.Pair == job.Pair {
					continue
				}

				simulationData = append(simulationData, SimulationData{
					Pairs: []models.SwapContractParams{
						{
							Pair:        pair.Pair,
							PairVersion: pair.Version,
						},
						{
							Pair: job.Pair,
							PairVersion: func() int {
								if job.PairType == "v2" {
									return 2
								} else {
									return 3
								}
							}(),
						},
					},
					Txs: job.TxsOut,
				},
				)
			}
			// so it was a purchase
			//  we need to buy on another pair and sell here
		} else {
			// its sell
			// we need to buy here and sell on another pair
			for _, pair := range pairs {
				if pair.Pair == job.Pair {
					continue
				}

				simulationData = append(simulationData, SimulationData{

					Pairs: []models.SwapContractParams{
						{
							Pair: job.Pair,
							PairVersion: func() int {
								if job.PairType == "v2" {
									return 2
								} else {
									return 3
								}
							}(),
						},
						{
							Pair:        pair.Pair,
							PairVersion: pair.Version,
						},
					},
					Txs: job.TxsIn,
				},
				)
			}
		}
	}

	var mostProfitableResult *models.SwapResponse
	var mostProfitableData SimulationData

	for _, data := range simulationData {
		result, err := t.wrapper.Simulate(ctx, models.Data2Simulate{
			Pairs:             data.Pairs,
			InputToken:        usecases.Weth(),
			Contract:          t.contractAddr,
			Transactions:      data.Txs.ArgsList(),
			BlockNumberOrHash: job.BlockNumberOrHash,
			BlockOverrides:    job.BlockOverrides,
		})
		if err != nil {
			log.Info(fmt.Sprintf("%s", err.Error()))
			continue
		}

		pairsStr := ""
		for _, pair := range data.Pairs {
			pairsStr = fmt.Sprintf("%s_%d  ->", pair.Pair, pair.PairVersion)
		}

		data2print := fmt.Sprintf(
			"path: %s\n "+
				"txs: %s\n"+
				"amount: %s\n"+
				"profit: %s\n"+
				"duration: %d\n"+
				"reason: %s\n",
			pairsStr,
			data.Txs.Hashes(),
			result.OptimalValue.String(),
			result.Profit.String(),
			result.Duration,
			result.Reason,
		)

		if mostProfitableResult == nil {
			mostProfitableResult = result
			mostProfitableData = data
		} else {
			if result.Profit.Cmp(mostProfitableResult.Profit) == 1 {
				mostProfitableResult = result
				mostProfitableData = data
			}
		}

		log.Info(data2print)
	}

	pairsStr := ""
	for _, pair := range mostProfitableData.Pairs {
		pairsStr = fmt.Sprintf("%s_%d  ->", pair.Pair, pair.PairVersion)
	}

	data2print := fmt.Sprintf(
		"path: %s\n "+
			"txs: %s\n"+
			"amount: %s\n"+
			"profit: %s\n"+
			"duration: %d\n"+
			"reason: %s\n",
		pairsStr,
		mostProfitableData.Txs.Hashes(),
		mostProfitableResult.OptimalValue.String(),
		mostProfitableResult.Profit.String(),
		mostProfitableResult.Duration,
		mostProfitableResult.Reason,
	)

	log.Info(data2print)

	return &models.SimulationJobResponse{}, nil
}

type SimulationData struct {
	Pairs []models.SwapContractParams
	Txs   models.TransactionsGroupsInfo
}
