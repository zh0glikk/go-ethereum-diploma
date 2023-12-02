package protocol

import (
	"github.com/ethereum/go-ethereum/internal/arbitrage/transactor/simulation_models"
)

type Transactor interface {
	New() Transactor

	// PrepareTemplates(data *simulation_models.PrepareTemplatesDTO) error

	Pack(data *simulation_models.PackFrontDTO) ([]byte, error)
	// PackSell(data *simulation_models.PackBackDTO) ([]byte, error)
}
