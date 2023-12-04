package algo

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"math"
	"math/big"
	"time"
)

const (
	defaultMaxBruteTime   = 5000 * time.Millisecond
	subWorkersCount       = 100
	iterationsLimit       = 20
	DefaultSplitParam     = 10
	significantBytesCount = 5
)

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)
)

type Worker struct {
	ch     chan *SubWorkerJob
	chResp chan *SubWorkerJobResponse

	workerChanResp     chan *ResultWorkerResponse
	newIterationJob    chan *NewIterationJob
	finishSubWorkersCh chan struct{}

	gasUsageF GasUsageFunc

	splitParam          *big.Int
	maxDepth            int
	maxBruteTime        time.Duration
	exceededMaxDuration bool

	settings WorkerSettings

	startPoints []*big.Int
}

func NewWorker(settings WorkerSettings) *Worker {
	w := &Worker{
		ch:                  make(chan *SubWorkerJob, len(criticalPoints)),
		chResp:              make(chan *SubWorkerJobResponse, len(criticalPoints)),
		workerChanResp:      make(chan *ResultWorkerResponse),
		newIterationJob:     make(chan *NewIterationJob),
		finishSubWorkersCh:  make(chan struct{}),
		splitParam:          big.NewInt(DefaultSplitParam),
		startPoints:         []*big.Int{},
		maxDepth:            iterationsLimit,
		maxBruteTime:        defaultMaxBruteTime,
		gasUsageF:           defaultGasUsageFunc,
		exceededMaxDuration: false,
		settings:            settings,
	}
	w.initSubWorkers()

	return w
}

func (w *Worker) initSubWorkers() {
	for i := 0; i < subWorkersCount; i++ {
		go runSubWorker(w.ch, w.chResp, w.finishSubWorkersCh)
	}
}

func (w *Worker) runJobManager() {
	subWorkersFinished := 0
	currentPoints := w.startPoints
	currentIteration := 0

	var bestExecution []models.CallManyResponseDTO
	bestProfit := big.NewInt(0)
	bestIndex := 0
	bestValue := big.NewInt(0)
	lastDiff := big.NewInt(0)

	for {
		select {
		case resp := <-w.chResp:
			subWorkersFinished += 1

			if resp.Profit.Cmp(bestProfit) != -1 {
				bestProfit = resp.Profit
				bestValue = resp.Point
				bestExecution = resp.Execution
				bestIndex = resp.Index
			}

			// log.Info(fmt.Sprintf("tick: %d/%d point: %s profit: %s", subWorkersFinished, len(currentPoints), resp.Point.String(), resp.Profit.String()))

			// check that iteration finished
			if subWorkersFinished == len(currentPoints) {
				currentIteration += 1
				if w.exceededMaxDuration {
					w.workerChanResp <- &ResultWorkerResponse{
						BestValue:  bestValue,
						BestProfit: bestProfit,
						Execution:  bestExecution,
						Reason:     fmt.Sprintf("Success in %d iterations. Max brute time reached. ", currentIteration),
					}
					return
				}

				if currentIteration >= w.maxDepth {
					w.workerChanResp <- &ResultWorkerResponse{
						BestValue:  bestValue,
						BestProfit: bestProfit,
						Execution:  bestExecution,
						Reason:     fmt.Sprintf("Success in %d iterations. Max depth reached. ", currentIteration),
					}
					return
				}

				index := bestIndex
				amountFrom := big.NewInt(0)
				amountTo := big.NewInt(0)

				if index == 0 {
					amountFrom = bestValue
					amountTo = currentPoints[index+1]
				} else if index == len(currentPoints)-1 {
					amountFrom = currentPoints[index-1]
					amountTo = bestValue
				} else {
					amountFrom = currentPoints[index-1]
					amountTo = currentPoints[index+1]
				}

				// log.Info(fmt.Sprintf("preparing new points: min: %s max %s", amountFrom.String(), amountTo.String()))
				newPoints, diff, ind := PrepareNewPointsWithMerge(
					amountFrom,
					amountTo,
					w.splitParam,
					bestValue,
				)
				bestIndex = ind
				if diff.Cmp(one) != 1 || diff.Cmp(lastDiff) == 0 {
					w.workerChanResp <- &ResultWorkerResponse{
						BestValue:  bestValue,
						BestProfit: bestProfit,
						Execution:  bestExecution,
						Reason:     fmt.Sprintf("Success in %d iterations. Min diff reached.", currentIteration),
					}
					return
				}
				lastDiff = diff

				subWorkersFinished = 0
				currentPoints = newPoints

				w.newIterationJob <- &NewIterationJob{Points: newPoints}

				continue
			}
		}
	}
}

func (w *Worker) Execute(
	ctx context.Context,
	extraData InitialExtraData,
	f AlgoFunc,
) *ResultWorkerResponse {
	w.prepareInitialParams(extraData)

	go w.runJobManager()

	// start first iteration
	for ind, point := range w.startPoints {
		w.ch <- &SubWorkerJob{
			Point: point,
			F:     f,
			Index: ind,
		}
	}

	algoTimeoutTicker := time.NewTicker(w.maxBruteTime)

	for {
		select {
		case <-algoTimeoutTicker.C:
			// after this tick we must wait until current iteration finished and stop algo
			w.exceededMaxDuration = true
		case <-ctx.Done():
			close(w.finishSubWorkersCh)
			return &ResultWorkerResponse{
				BestValue:  big.NewInt(0),
				BestProfit: big.NewInt(0),
				Execution:  nil,
				GasInfo:    nil,
				Reason:     "Context canceled. ",
			}
		case res := <-w.workerChanResp:
			close(w.finishSubWorkersCh)
			return res
		case newIteration := <-w.newIterationJob:
			for ind, point := range newIteration.Points {
				w.ch <- &SubWorkerJob{
					Point: point,
					F:     f,
					Index: ind,
				}
			}
		}
	}
}

func (w *Worker) prepareInitialParams(extraData InitialExtraData) {
	points := extraData.Points
	if len(points) == 0 {
		points = criticalPoints
	}

	w.startPoints = points

	if extraData.SplitParam != nil {
		w.splitParam = extraData.SplitParam
	}
	if extraData.MaxDepth != 0 {
		w.maxDepth = extraData.MaxDepth
	}
	if extraData.MaxBruteTime != 0 {
		w.maxBruteTime = time.Duration(extraData.MaxBruteTime) * time.Millisecond
	}
}

func RoundPoint(point *big.Int) *big.Int {
	if point.Cmp(big.NewInt(0)) == 0 {
		return point
	}

	pointF, _ := new(big.Float).SetInt(point).Float64()

	k := math.Ceil(math.Log2(pointF) / 8)
	y := k - significantBytesCount
	multiplierF, _ := big.NewFloat(math.Pow(256, y)).Float64()
	res := math.Floor(pointF/multiplierF) * multiplierF

	value := new(big.Int)
	_, _ = new(big.Float).SetFloat64(res).Int(value)

	return value
}
