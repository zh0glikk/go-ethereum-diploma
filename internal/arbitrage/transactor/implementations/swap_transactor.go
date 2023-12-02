package implementations

import (
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/packer"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/protocol"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/simulation_models"
)

type SwapTransactor struct {
	purchaseTemplate string
	purchaseType     int

	sellTemplate string
	sellType     int
}

func NewSwapTransactor() *SwapTransactor {
	return &SwapTransactor{}
}

func (i *SwapTransactor) New() protocol.Transactor {
	return &SwapTransactor{
		purchaseTemplate: i.purchaseTemplate,
		purchaseType:     i.purchaseType,
		sellTemplate:     i.sellTemplate,
		sellType:         i.sellType,
	}
}

func (i *SwapTransactor) PrepareTemplates(data *simulation_models.PrepareTemplatesDTO) error {
	var err error

	i.purchaseType = data.InputPairVersion
	i.sellType = data.OutputPairVersion

	if data.InputPairVersion == 2 {
		i.purchaseTemplate, err = packer.PackerObj.PackSwapV2Template(data.InputPair, data.InputToken, data.OutputToken)
		if err != nil {
			return err
		}
	} else {
		i.purchaseTemplate, err = packer.PackerObj.PackSwapV3Template(data.InputPair, data.InputToken, data.OutputToken, data.Contract)
		if err != nil {
			return err
		}
	}

	if data.OutputPairVersion == 2 {
		i.purchaseTemplate, err = packer.PackerObj.PackSwapV2Template(data.OutputPair, data.OutputToken, data.InputToken)
		if err != nil {
			return err
		}
	} else {
		i.purchaseTemplate, err = packer.PackerObj.PackSwapV3Template(data.OutputPair, data.OutputToken, data.InputToken, data.Contract)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *SwapTransactor) PackPurchase(data *simulation_models.PackFrontDTO) ([]byte, error) {
	if data.Value == nil {
		return nil, nil
	}
	if i.purchaseType == 2 {
		return packer.PackerObj.PackSwapV2(data.Value, i.purchaseTemplate)
	} else {
		return packer.PackerObj.PackSwapV3(data.Value, i.purchaseTemplate)
	}
}

func (i *SwapTransactor) PackSell(data *simulation_models.PackBackDTO) ([]byte, error) {
	if data.Value == nil {
		return nil, nil
	}

	if i.sellType == 2 {
		return packer.PackerObj.PackSwapV2(data.Value, i.sellTemplate)
	} else {
		return packer.PackerObj.PackSwapV3(data.Value, i.sellTemplate)
	}
}
