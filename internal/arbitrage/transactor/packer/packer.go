package packer

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Packer interface {
	PackSwapV2(
		inputAmount *big.Int,
		template string) ([]byte, error)

	PackSwapV2Template(
		pair common.Address,
		input common.Address,
		output common.Address) (string, error)

	PackSwapV3(
		inputAmount *big.Int,
		template string) ([]byte, error)

	PackSwapV3Template(
		input common.Address,
		output common.Address,
		pair common.Address,
		contract common.Address,
	) (string, error)

	PackBalanceOf(
		address common.Address,
	) ([]byte, error)

	PackToken0() ([]byte, error)
	PackToken1() ([]byte, error)
}

func NewPacker() Packer {
	return newPacker{}
}
