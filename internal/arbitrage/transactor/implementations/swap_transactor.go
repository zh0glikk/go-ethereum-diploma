package implementations

import (
	"fmt"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/packer"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/protocol"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/simulation_models"
	"github.com/ethereum/go-ethereum/log"
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

	log.Info(fmt.Sprintf("data.InputPairVersion: %d", data.InputPairVersion))
	log.Info(fmt.Sprintf("data.OutputPairVersion: %d", data.OutputPairVersion))

	log.Info(fmt.Sprintf("%s %s %s %s", data.InputPair.String(), data.OutputPair.String(), data.InputToken.String(), data.OutputToken.String()))
	if data.InputPairVersion == 2 {
		i.purchaseTemplate, err = packer.PackerObj.PackSwapV2Template(data.InputPair, data.InputToken, data.OutputToken)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("packing swapV2 input template: %s", i.purchaseTemplate))

	} else {
		i.purchaseTemplate, err = packer.PackerObj.PackSwapV3Template(data.InputPair, data.InputToken, data.OutputToken, data.Contract)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("packing swapV3 input template: %s", i.purchaseTemplate))

	}

	if data.OutputPairVersion == 2 {
		i.sellTemplate, err = packer.PackerObj.PackSwapV2Template(data.OutputPair, data.OutputToken, data.InputToken)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("packing swapV2 output template: %s", i.sellTemplate))

	} else {
		i.sellTemplate, err = packer.PackerObj.PackSwapV3Template(data.OutputPair, data.OutputToken, data.InputToken, data.Contract)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("packing swapV3 output template: %s", i.sellTemplate))

	}

	return nil
}

func (i *SwapTransactor) PackPurchase(data *simulation_models.PackFrontDTO) ([]byte, error) {
	log.Info(fmt.Sprintf("purchase %d %v", i.purchaseType, data.Value))
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
	log.Info(fmt.Sprintf("sell %d %v", i.sellType, data.Value))

	if data.Value == nil {
		return nil, nil
	}

	if i.sellType == 2 {
		return packer.PackerObj.PackSwapV2(data.Value, i.sellTemplate)
	} else {
		return packer.PackerObj.PackSwapV3(data.Value, i.sellTemplate)
	}
}
