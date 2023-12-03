package algo

import (
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"math/big"
)

// returns execution, profit
type AlgoFunc func(v *big.Int) ([]models.CallManyResponseDTO, *big.Int)

// GasUsageFunc - returns fee and map with gas info
type GasUsageFunc func(inputAmount *big.Int, execution []models.CallManyResponseDTO) (*big.Int, map[string]*big.Int)

func defaultGasUsageFunc(inputAmount *big.Int, execution []models.CallManyResponseDTO) (*big.Int, map[string]*big.Int) {
	return big.NewInt(0), map[string]*big.Int{}
}

type WorkerSettings struct {
	// DisableFirstPointCheck bool
	// DisableMinProfitCheck  bool
	// DisablePointRounding   bool // rounded first 5 byte
}

type InitialExtraData struct {
	MaxDepth     int
	MaxBruteTime int64
	Points       []*big.Int
	SplitParam   *big.Int
	GasUsageF    *GasUsageFunc
}

type SubWorkerJob struct {
	Point *big.Int
	Index int
	F     AlgoFunc
}

type SubWorkerJobResponse struct {
	Point     *big.Int
	Index     int
	Profit    *big.Int
	Execution []models.CallManyResponseDTO
}

type ResultWorkerResponse struct {
	BestValue  *big.Int
	BestProfit *big.Int
	Execution  []models.CallManyResponseDTO
	GasInfo    map[string]*big.Int
	Reason     string
}

type NewIterationJob struct {
	Points []*big.Int
}
