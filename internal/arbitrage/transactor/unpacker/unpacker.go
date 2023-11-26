package unpacker

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"math/big"
)

type Unpacker interface {
	// ParseOutputAmount(
	// 	resp []models.CallManyResponseDTO,
	// 	tx TxType,
	// ) (*big.Int, error)

	// ParseInputAmount(resp []models.CallManyResponseDTO, tx TxType) (*big.Int, error)
	// ParseReceived(resp []models.CallManyResponseDTO, tx TxType) (*big.Int, error)
	ParseSqrtPriceLimit(resp []models.CallManyResponseDTO, tx TxType) (*big.Int, error)
	// ParseProductionGasSpent(resp []models.CallManyResponseDTO) (*big.Int, *big.Int, error)
	ParseBalanceOf(resp []models.CallManyResponseDTO) (*big.Int, error)
	// ParseReserveFirstToken(resp []models.CallManyResponseDTO) (*big.Int, error)
	// ParseReserveSecondToken(resp []models.CallManyResponseDTO) (*big.Int, error)
	ParseFee(resp []models.CallManyResponseDTO) (*big.Int, error)
	ParseGetPair(resp []models.CallManyResponseDTO) (common.Address, error)
	// ParseGetAmountOutput(resp []models.CallManyResponseDTO) (*big.Int, error)
	ParseFactory(resp []models.CallManyResponseDTO) (common.Address, error)
	ParseGetPool(resp []models.CallManyResponseDTO) (common.Address, error)
	IsTransfer(input []byte) bool
	IsTransferFrom(input []byte) bool
	IsSwapV2(input []byte) bool
	IsSwapV3(input []byte) bool
	UnpackTransfer(input []byte) (common.Address, *big.Int)
	UnpackTransferFrom(input []byte) (common.Address, common.Address, *big.Int)
}

func NewUnpacker() Unpacker {
	return newUnpacker{}
}
