package simulator

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"math/big"
	"testing"
	"time"
)

func TestSort(t *testing.T) {
	mock := []TxWrapper{
		{
			Profit: big.NewInt(100000000000000000),
			TxArgs: []ethapi.TransactionArgs{{}},
			Tx:     []types.Transaction{},
		},
		{
			Profit: big.NewInt(100000000000000000),
			TxArgs: []ethapi.TransactionArgs{},
			Tx:     []types.Transaction{},
		},
		{
			Profit: big.NewInt(349571194872656306),
			TxArgs: []ethapi.TransactionArgs{},
			Tx:     []types.Transaction{},
		},
		{
			Profit: big.NewInt(349571194872656306),
			TxArgs: []ethapi.TransactionArgs{},
			Tx:     []types.Transaction{},
		},
	}

	startedAt := time.Now().UTC()
	_ = bubbleSort(mock, true)
	fmt.Println(time.Now().UTC().Sub(startedAt).Milliseconds())
}
