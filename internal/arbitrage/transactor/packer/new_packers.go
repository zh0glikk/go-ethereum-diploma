package packer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type newPacker struct{}

func (_ newPacker) PackSwapV2(
	inputAmount *big.Int,
	template string) ([]byte, error) {
	return mustDecodeHex(fmt.Sprintf(
		template,
		hex.EncodeToString(encodeUint256(inputAmount)),
	)), nil
}

func (_ newPacker) PackSwapV2Template(
	pair common.Address,
	input common.Address,
	output common.Address) (string, error) {
	tmpl := hex.EncodeToString(bytes.Join([][]byte{
		encodeAddress(pair.Bytes()),
		encodeAddress(input.Bytes()),
		encodeAddress(output.Bytes()),
	}, nil))

	return swapV2Sig + tmpl + "%s", nil
}

func (_ newPacker) PackSwapV3(
	inputAmount *big.Int,
	template string) ([]byte, error) {
	return mustDecodeHex(fmt.Sprintf(
		template,
		hex.EncodeToString(encodeUint256(inputAmount)),
	)), nil
}

func (_ newPacker) PackSwapV3Template(
	input common.Address,
	output common.Address,
	pair common.Address,
	contract common.Address,
) (string, error) {
	tmpl := hex.EncodeToString(bytes.Join([][]byte{
		encodeAddress(input.Bytes()),
		encodeAddress(output.Bytes()),
		encodeAddress(pair.Bytes()),
		encodeUint256(zero),
		encodeAddress(contract.Bytes()),
	}, nil))

	return swapV3Sig + "%s" + tmpl, nil
}

func (_ newPacker) PackBalanceOf(
	address common.Address,
) ([]byte, error) {
	return bytes.Join(
		[][]byte{
			balanceOfSigBytes,
			encodeAddress(address.Bytes()),
		}, nil), nil
}

func (_ newPacker) PackToken0() ([]byte, error) {
	return mustDecodeHex(fmt.Sprintf("%s", token0Sig)), nil
}

func (_ newPacker) PackToken1() ([]byte, error) {
	return mustDecodeHex(fmt.Sprintf("%s", token1Sig)), nil
}
