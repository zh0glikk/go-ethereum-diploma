package implementations

import (
	"fmt"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/packer"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/protocol"
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/simulation_models"
	"github.com/ethereum/go-ethereum/log"
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

// func (i *SwapTransactor) PrepareTemplates(data *simulation_models.PrepareTemplatesDTO) error {
// 	var err error
//
// 	i.purchaseType = data.InputPairVersion
// 	i.sellType = data.OutputPairVersion
//
// 	log.Info(fmt.Sprintf("data.InputPairVersion: %d", data.InputPairVersion))
// 	log.Info(fmt.Sprintf("data.OutputPairVersion: %d", data.OutputPairVersion))
//
// 	log.Info(fmt.Sprintf("%s %s %s %s", data.InputPair.String(), data.OutputPair.String(), data.InputToken.String(), data.OutputToken.String()))
// 	if data.InputPairVersion == 2 {
//
// 	} else {
// 		i.purchaseTemplate, err = packer.PackerObj.PackSwapV3Template(data.InputToken, data.OutputToken, data.InputPair, data.Contract)
// 		if err != nil {
// 			return err
// 		}
// 		log.Info(fmt.Sprintf("packing swapV3 input template: %s", i.purchaseTemplate))
//
// 	}
//
// 	if data.OutputPairVersion == 2 {
// 		i.sellTemplate, err = packer.PackerObj.PackSwapV2Template(data.OutputPair, data.OutputToken, data.InputToken)
// 		if err != nil {
// 			return err
// 		}
// 		log.Info(fmt.Sprintf("packing swapV2 output template: %s", i.sellTemplate))
//
// 	} else {
// 		i.sellTemplate, err = packer.PackerObj.PackSwapV3Template(data.OutputToken, data.InputToken, data.OutputPair, data.Contract)
// 		if err != nil {
// 			return err
// 		}
// 		log.Info(fmt.Sprintf("packing swapV3 output template: %s", i.sellTemplate))
//
// 	}
//
// 	return nil
// }

func (i *SwapTransactor) Pack(data *simulation_models.PackFrontDTO) ([]byte, error) {
	log.Info(fmt.Sprintf("purchase %d %v", data.PairType, data.Value))
	if data.Value == nil {
		return nil, nil
	}
	if data.PairType == 2 {
		tmpl, err := packer.PackerObj.PackSwapV2Template(data.Pair, data.Input, data.Output)
		if err != nil {
			return nil, err
		}
		log.Info(fmt.Sprintf("packing swapV2 input template: %s", data.PairType))

		return packer.PackerObj.PackSwapV2(data.Value, tmpl)
	} else {
		tmpl, err := packer.PackerObj.PackSwapV3Template(data.Input, data.Output, data.Pair, data.Contract)
		if err != nil {
			return nil, err
		}
		log.Info(fmt.Sprintf("packing swapV3 input template: %s", data.PairType))

		return packer.PackerObj.PackSwapV3(data.Value, tmpl)
	}
}
