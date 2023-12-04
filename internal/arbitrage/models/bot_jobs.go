package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"time"
)

type WaitingTraceTxsJob struct {
	From   common.Address
	Args   []ethapi.TransactionArgs
	RawTxs []types.Transaction

	Swaps [][]SwapInfo

	BlockNumberOrHash rpc.BlockNumberOrHash
	BlockOverrides    *ethapi.BlockOverrides

	CreatedAt time.Time
}

type SimulationNotifyJob struct {
	BlockNumberOrHash rpc.BlockNumberOrHash
	BlockOverrides    *ethapi.BlockOverrides

	Pairs []common.Address
}

type SimulationJob struct {
	BlockNumberOrHash rpc.BlockNumberOrHash
	BlockOverrides    *ethapi.BlockOverrides
	Pair              common.Address
	PairType          string
	Token0            common.Address
	Token1            common.Address
	TxsIn             TransactionsGroupsInfo
	TxsOut            TransactionsGroupsInfo
}

type SimulationJobResponse struct {
	// Value     *big.Int
	// Indices   []uint64
	// Execution []CallManyResponseDTO
	//
	// Input        common.Address
	// Output       common.Address
	// Pair         common.Address
	// PrebundleTxs TransactionsGroupsInfo
	// VictimTxs    TransactionsGroupsInfo

	// Method SandwichMethod
}

type TransactionsGroupsInfo struct {
	Transactions [][]types.Transaction
	Args         [][]ethapi.TransactionArgs
}

func (t *TransactionsGroupsInfo) ArgsList() []ethapi.TransactionArgs {
	var result []ethapi.TransactionArgs

	for _, arg := range t.Args {
		result = append(result, arg...)
	}
	return result
}

func (t *TransactionsGroupsInfo) Hashes() string {
	var result string

	for _, txs := range t.Transactions {
		for _, tx := range txs {
			// result = append(result, tx.Hash())
			result += tx.Hash().String() + ", "
		}
	}

	return result
}

func (t *TransactionsGroupsInfo) Contains(hash common.Hash) bool {
	for _, txs := range t.Transactions {
		for _, tx := range txs {
			if tx.Hash() == hash {
				return true
			}
		}
	}
	return false
}

func (t *TransactionsGroupsInfo) GroupHashes(i int) string {
	var result string

	txs := t.Transactions[i]
	for _, tx := range txs {
		result += tx.Hash().String() + ", "
	}

	return result
}

type Data2Simulate struct {
	Pairs      []SwapContractParams
	InputToken common.Address
	Contract   common.Address

	Transactions []ethapi.TransactionArgs

	BlockNumberOrHash rpc.BlockNumberOrHash
	BlockOverrides    *ethapi.BlockOverrides

	Points []*big.Int
}
