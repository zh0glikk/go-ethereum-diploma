package usecases

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/packer"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/unpacker"
	"github.com/ethereum/go-ethereum/internal/arbitrage/utils"
	"github.com/ethereum/go-ethereum/internal/ethapi"
)

func CallToken0(b ethapi.Backend, to common.Address, stateDB *state.StateDB, header *types.Header, vmctx vm.BlockContext) (common.Address, error) {
	data, err := packer.PackerObj.PackToken0()
	if err != nil {
		return common.Address{}, err
	}

	execution, _, _, _, _, _, err := DoCallManyOnStateReturningState(
		context.Background(),
		b,
		[]ethapi.TransactionArgs{
			{
				To:   &to,
				Data: utils.Ptr(hexutil.Bytes(data)),
			},
		},
		stateDB,
		header,
		vmctx,
		b.RPCEVMTimeout(),
		b.RPCGasCap(),
		CheckRevertAll,
		nil,
	)
	if err != nil {
		return common.Address{}, err
	}

	return unpacker.UnpackerObj.ParseAddress(execution)
}

func CallToken1(b ethapi.Backend, to common.Address, stateDB *state.StateDB, header *types.Header, vmctx vm.BlockContext) (common.Address, error) {
	data, err := packer.PackerObj.PackToken1()
	if err != nil {
		return common.Address{}, err
	}

	execution, _, _, _, _, _, err := DoCallManyOnStateReturningState(
		context.Background(),
		b,
		[]ethapi.TransactionArgs{
			{
				To:   &to,
				Data: utils.Ptr(hexutil.Bytes(data)),
			},
		},
		stateDB,
		header,
		vmctx,
		b.RPCEVMTimeout(),
		b.RPCGasCap(),
		CheckRevertAll,
		nil,
	)
	if err != nil {
		return common.Address{}, err
	}

	return unpacker.UnpackerObj.ParseAddress(execution)
}
