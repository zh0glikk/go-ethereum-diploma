package storage

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func sortTxsByNonce(txs []types.Transaction) []types.Transaction {
	n := len(txs)
	// Perform n-1 passes
	for i := 0; i < n-1; i++ {
		// Last i elements are already in place, so we don't need to check them
		for j := 0; j < n-i-1; j++ {
			// Swap if the element found is greater than the next element
			if txs[j].Nonce() > txs[j+1].Nonce() {
				txs[j], txs[j+1] = txs[j+1], txs[j]
			}
		}
	}

	return txs
}

func uniqueTxs(txs []types.Transaction) []types.Transaction {
	unique := map[common.Hash]bool{}
	var result []types.Transaction

	for _, tx := range txs {
		if _, ok := unique[tx.Hash()]; !ok {
			result = append(result, tx)
			unique[tx.Hash()] = true
		}
	}

	return result
}
