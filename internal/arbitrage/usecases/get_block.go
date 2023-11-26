package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

func BlockByNumber(ctx context.Context, b ethapi.Backend, number rpc.BlockNumber) (*types.Block, error) {
	block, err := b.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, fmt.Errorf("block #%d not found", number)
	}
	return block, nil
}

func BlockByHash(ctx context.Context, b ethapi.Backend, hash common.Hash) (*types.Block, error) {
	block, err := b.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, fmt.Errorf("block %s not found", hash.Hex())
	}
	return block, nil
}

func BlockByNumberAndHash(ctx context.Context, b ethapi.Backend, number rpc.BlockNumber, hash common.Hash) (*types.Block, error) {
	block, err := BlockByNumber(ctx, b, number)
	if err != nil {
		return nil, err
	}
	if block.Hash() == hash {
		return block, nil
	}
	return BlockByHash(ctx, b, hash)
}

func BlockByNumberOrHash(ctx context.Context, b ethapi.Backend, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block, error) {
	var (
		err   error
		block *types.Block
	)
	if hash, ok := blockNrOrHash.Hash(); ok {
		block, err = BlockByHash(ctx, b, hash)
	} else if number, ok := blockNrOrHash.Number(); ok {
		if number == rpc.PendingBlockNumber {
			return nil, errors.New("tracing on top of pending is not supported")
		}
		block, err = BlockByNumber(ctx, b, number)
	} else {
		return nil, errors.New("invalid arguments; neither block nor hash specified")
	}
	return block, err
}

func GetHeaderByNumber(ctx context.Context, b ethapi.Backend, number rpc.BlockNumber) (map[string]interface{}, error) {
	headerPending, err := b.HeaderByNumber(ctx, rpc.PendingBlockNumber)
	headerNext, err := b.HeaderByNumber(ctx, number+1)
	header, err := b.HeaderByNumber(ctx, number)

	if header != nil && err == nil {
		var nextBaseFee *big.Int
		if headerNext == nil {
			nextBaseFee = headerPending.BaseFee
		} else {
			nextBaseFee = headerNext.BaseFee
		}
		response := rpcMarshalHeader(ctx, b, header)

		response["nextBaseFee"] = (*hexutil.Big)(nextBaseFee)

		if number == rpc.PendingBlockNumber {
			// Pending header need to nil out a few fields
			for _, field := range []string{"hash", "nonce", "miner"} {
				response[field] = nil
			}
		}
		return response, err
	}
	return nil, err
}

func rpcMarshalHeader(ctx context.Context, b ethapi.Backend, header *types.Header) map[string]interface{} {
	fields := RPCMarshalHeader(header)
	fields["totalDifficulty"] = (*hexutil.Big)(b.GetTd(ctx, header.Hash()))
	return fields
}

func RPCMarshalHeader(head *types.Header) map[string]interface{} {
	result := map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             head.Hash(),
		"parentHash":       head.ParentHash,
		"nonce":            head.Nonce,
		"mixHash":          head.MixDigest,
		"sha3Uncles":       head.UncleHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"miner":            head.Coinbase,
		"difficulty":       (*hexutil.Big)(head.Difficulty),
		"extraData":        hexutil.Bytes(head.Extra),
		"gasLimit":         hexutil.Uint64(head.GasLimit),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        hexutil.Uint64(head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}
	if head.BaseFee != nil {
		result["baseFeePerGas"] = (*hexutil.Big)(head.BaseFee)
	}
	if head.WithdrawalsHash != nil {
		result["withdrawalsRoot"] = head.WithdrawalsHash
	}
	if head.BlobGasUsed != nil {
		result["blobGasUsed"] = hexutil.Uint64(*head.BlobGasUsed)
	}
	if head.ExcessBlobGas != nil {
		result["excessBlobGas"] = hexutil.Uint64(*head.ExcessBlobGas)
	}
	if head.ParentBeaconRoot != nil {
		result["parentBeaconBlockRoot"] = head.ParentBeaconRoot
	}
	return result
}
