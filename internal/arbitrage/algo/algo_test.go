package algo

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"math/big"
	"testing"
	"time"
)

func TestAlgo(t *testing.T) {
	w := NewWorker(WorkerSettings{})

	min := big.NewInt(100000000000000)
	max, _ := big.NewInt(0).SetString("987654321123456789000", 10)

	initialPoints, _, _ := PrepareNewPointsWithMerge(min, max, big.NewInt(DefaultSplitParam), nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	res := w.Execute(
		ctx,
		InitialExtraData{
			MaxDepth:     100,
			MaxBruteTime: 100,
			Points:       initialPoints,
			SplitParam:   big.NewInt(50),
		},
		func(
			v *big.Int,
		) ([]models.CallManyResponseDTO, *big.Int) {
			valueToFind, _ := big.NewInt(0).SetString("230399959613046783", 10)

			if v.Cmp(valueToFind) == 1 {
				return []models.CallManyResponseDTO{}, big.NewInt(0)
			}

			if v.Cmp(valueToFind) == -1 {
				return []models.CallManyResponseDTO{}, new(big.Int).Quo(v, big.NewInt(10))
			} else {
				return []models.CallManyResponseDTO{}, new(big.Int).Sub(v, big.NewInt(1))
			}
		})

	fmt.Println(fmt.Sprintf("RES: %s", res.BestValue.String()))
	fmt.Println(fmt.Sprintf("reason: %s", res.Reason))
	fmt.Println(fmt.Sprintf("results: %v", res.Execution))

	// time.Sleep(time.Second * 1)
}

func TestInitialPoints(t *testing.T) {
	balance, _ := new(big.Int).SetString("364510576149187356599391034456", 10)

	// 364510576149187356599391034456
	points, _, index := PrepareNewPointsWithMerge(new(big.Int).Quo(balance, big.NewInt(200)), balance, big.NewInt(50), nil)
	fmt.Println(index)
	for _, p := range points {
		fmt.Println(p.String())
	}
}
