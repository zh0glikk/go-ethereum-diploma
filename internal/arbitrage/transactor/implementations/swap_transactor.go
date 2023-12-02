package implementations

import (
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/packer"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/protocol"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/simulation_models"
)

type SwapTransactor struct {
	// purchaseTemplate string
	// purchaseType     int
	//
	// sellTemplate string
	// sellType     int
}

func NewSwapTransactor() *SwapTransactor {
	return &SwapTransactor{}
}

func (i *SwapTransactor) New() protocol.Transactor {
	return &SwapTransactor{}
}

func (i *SwapTransactor) Pack(data *simulation_models.PackFrontDTO) ([]byte, error) {
	if data.Value == nil {
		return nil, nil
	}
	if data.PairType == 2 {
		tmpl, err := packer.PackerObj.PackSwapV2Template(data.Pair, data.Input, data.Output)
		if err != nil {
			return nil, err
		}
		return packer.PackerObj.PackSwapV2(data.Value, tmpl)
	} else {
		tmpl, err := packer.PackerObj.PackSwapV3Template(data.Input, data.Output, data.Pair, data.Contract)
		if err != nil {
			return nil, err
		}
		return packer.PackerObj.PackSwapV3(data.Value, tmpl)
	}
}
