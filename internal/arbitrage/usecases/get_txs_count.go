package usecases

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
)

func GetTransactionsCount(ctx context.Context, b ethapi.Backend, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (uint64, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok && blockNr == rpc.PendingBlockNumber {
		nonce, err := b.GetPoolNonce(ctx, address)
		if err != nil {
			return 0, err
		}
		return nonce, nil
	}
	// Resolve block number and use its state to ask for the nonce
	state, _, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return 0, err
	}

	nonce := state.GetNonce(address)
	return nonce, state.Error()
}
