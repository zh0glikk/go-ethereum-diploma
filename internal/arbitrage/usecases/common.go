package usecases

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/packer"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/protocol"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/simulation_models"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/unpacker"
	"github.com/ethereum/go-ethereum/internal/arbitrage/utils"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

func prepareSwapInitial(
	ctx context.Context,
	b ethapi.Backend,
	victim []ethapi.TransactionArgs,
	stateOverride *ethapi.StateOverride,
	blockNrOrHash rpc.BlockNumberOrHash,
	blockOverrides *ethapi.BlockOverrides,
) (
	prebundleExecution []models.CallManyResponseDTO,
	stateDB *state.StateDB,
	header *types.Header,
	vmctx vm.BlockContext,
	err error) {

	// on this step we apply overrides and execute prebundle transactions
	prebundleExecution, stateDB, header, vmctx, err = DoCallManyReturningState(
		ctx,
		b,
		victim,
		blockNrOrHash,
		stateOverride,
		b.RPCEVMTimeout(),
		b.RPCGasCap(),
		blockOverrides,
	)
	if err != nil {
		return
	}

	return
}

func applySwap(
	ctx context.Context,
	b ethapi.Backend,
	data *simulation_models.PackFrontDTO,
	transactor protocol.Transactor,
	contract common.Address,
	stateDB *state.StateDB,
	header *types.Header,
	vmctx vm.BlockContext,
	checkRevertIndex int,
) ([]models.CallManyResponseDTO, bool, []uint64, *state.StateDB, *types.Header, vm.BlockContext, error) {
	newTransactions := prepareSwapsTransactions(
		data,
		transactor,
		contract,
	)

	log.Info(fmt.Sprintf("%v", newTransactions))

	return DoCallManyOnStateReturningState(
		ctx,
		b,
		newTransactions,
		stateDB,
		header,
		vmctx,
		b.RPCEVMTimeout(),
		b.RPCGasCap(),
		checkRevertIndex, // we stop it only if buy failed
		nil,
	)
}

func prepareSwapsTransactions(
	frontDTO *simulation_models.PackFrontDTO,
	transactor protocol.Transactor,
	contract common.Address,
) []ethapi.TransactionArgs {
	var transactions []ethapi.TransactionArgs
	purchaseBB, _ := transactor.Pack(frontDTO)
	log.Info(fmt.Sprintf("swap: %s", hexutil.Encode(purchaseBB)))

	transactions = append(transactions, ethapi.TransactionArgs{
		To:   &contract,
		Data: utils.Ptr(hexutil.Bytes(purchaseBB)),
	})

	return transactions
}

func CallERC20BalanceOf(b ethapi.Backend, to common.Address, owner common.Address, stateDB *state.StateDB, header *types.Header, vmctx vm.BlockContext) (*big.Int, error) {
	data, err := packer.PackerObj.PackBalanceOf(owner)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return unpacker.UnpackerObj.ParseBalanceOf(execution)
}
