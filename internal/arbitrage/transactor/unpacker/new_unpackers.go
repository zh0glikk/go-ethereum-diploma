package unpacker

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

type newUnpacker struct{}

func (_ newUnpacker) ParseSqrtPriceLimit(resp []models.CallManyResponseDTO, tx TxType) (*big.Int, error) {
	if tx == Front {
		if resp[0].Value == nil {
			return nil, errors.New("sandwich failed")
		}
		return new(big.Int).SetBytes(resp[0].Value.(hexutil.Bytes)[64:96]), nil

	} else if tx == Back {
		if resp[len(resp)-1].Value == nil {
			return nil, errors.New("sandwich failed")
		}
		return new(big.Int).SetBytes(resp[len(resp)-1].Value.(hexutil.Bytes)[64:96]), nil
	}
	return nil, nil
}

func (_ newUnpacker) ParseBalanceOf(resp []models.CallManyResponseDTO) (*big.Int, error) {
	if resp[0].Value == nil || len(resp[0].Value.(hexutil.Bytes)) < 32 {
		return nil, errors.New("balanceOf failed")
	}
	return new(big.Int).SetBytes(resp[0].Value.(hexutil.Bytes)[:32]), nil
}

func (_ newUnpacker) ParseFee(resp []models.CallManyResponseDTO) (*big.Int, error) {
	if resp[0].Value == nil {
		return nil, errors.New("pairFee failed")
	}
	return new(big.Int).SetBytes(resp[0].Value.(hexutil.Bytes)[:32]), nil
}

func (_ newUnpacker) ParseOutputAmount(resp []models.CallManyResponseDTO) (*big.Int, error) {
	if resp[0].Value == nil || resp[0].Error != nil {
		log.Info(fmt.Sprintf("%v", resp))
		return nil, errors.New("ParseOutputAmount failed")
	}
	return new(big.Int).SetBytes(resp[0].Value.(hexutil.Bytes)[:32]), nil
}

func (_ newUnpacker) ParseGetPair(resp []models.CallManyResponseDTO) (common.Address, error) {
	if resp[0].Value == nil {
		return common.Address{}, errors.New("getPair failed")
	}

	return common.BytesToAddress(resp[0].Value.(hexutil.Bytes)[:32]), nil
}

func (_ newUnpacker) ParseAddress(resp []models.CallManyResponseDTO) (common.Address, error) {
	if resp[0].Value == nil {
		return common.Address{}, errors.New("parse addr failed")
	}

	return common.BytesToAddress(resp[0].Value.(hexutil.Bytes)[:32]), nil
}

func (_ newUnpacker) ParseGetAmountOutput(resp []models.CallManyResponseDTO) (*big.Int, error) {
	if len(resp) == 0 {
		return big.NewInt(0), nil
	}

	if resp[0].Value == nil {
		return nil, errors.New("getAmountOutput failed")
	}
	return new(big.Int).SetBytes(resp[0].Value.(hexutil.Bytes)[:32]), nil
}

func (_ newUnpacker) ParseFactory(resp []models.CallManyResponseDTO) (common.Address, error) {
	if resp[0].Value == nil {
		return common.Address{}, errors.New("getFactory failed")
	}
	if cap(resp[0].Value.(hexutil.Bytes)) < 32 {
		return common.Address{}, errors.New("getFactory failed")
	}
	return common.BytesToAddress(resp[0].Value.(hexutil.Bytes)[:32]), nil
}

func (_ newUnpacker) ParseGetPool(resp []models.CallManyResponseDTO) (common.Address, error) {
	if resp[0].Value == nil {
		return common.Address{}, errors.New("getPair failed")
	}

	return common.BytesToAddress(resp[0].Value.(hexutil.Bytes)[:32]), nil
}

func (_ newUnpacker) IsTransfer(input []byte) bool {
	if len(input) < 4 {
		return false
	}
	sig, err := hexutil.Decode(transferSign)
	if err != nil {
		return false
	}
	if bytes.Equal(sig, input[:4]) {
		return true
	}
	return false
}

func (_ newUnpacker) IsTransferFrom(input []byte) bool {
	if len(input) < 4 {
		return false
	}
	sig, err := hexutil.Decode(transferFromSign)
	if err != nil {
		return false
	}
	if bytes.Equal(sig, input[:4]) {
		return true
	}
	return false
}

func (_ newUnpacker) IsSwapV2(input []byte) bool {
	if len(input) < 4 {
		return false
	}
	sig, err := hexutil.Decode(swapSigV2Sign)
	if err != nil {
		return false
	}
	if bytes.Equal(sig, input[:4]) {
		return true
	}
	return false
}

func (_ newUnpacker) IsSwapV3(input []byte) bool {
	if len(input) < 4 {
		return false
	}
	sig, err := hexutil.Decode(swapSigV3Sign)
	if err != nil {
		return false
	}
	if bytes.Equal(sig, input[:4]) {
		return true
	}
	return false
}

// TODO: rename to parse or rename all to unpack
func (_ newUnpacker) UnpackTransfer(input []byte) (common.Address, *big.Int) {
	if len(input) < 68 {
		return common.Address{}, big.NewInt(0)
	}

	return common.BytesToAddress(input[4:36]), new(big.Int).SetBytes(input[36:68])
}

func (_ newUnpacker) UnpackTransferFrom(input []byte) (common.Address, common.Address, *big.Int) {
	if len(input) < 100 {
		return common.Address{}, common.Address{}, big.NewInt(0)
	}

	return common.BytesToAddress(input[4:36]),
		common.BytesToAddress(input[36:68]),
		new(big.Int).SetBytes(input[68:100])
}
