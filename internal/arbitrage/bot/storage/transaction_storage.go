package storage

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"sync"
)

type TxsModel struct {
	Pair     common.Address
	PairType string
	Token0   common.Address
	Token1   common.Address
	TxsIn    models.TransactionsGroupsInfo
	TxsOut   models.TransactionsGroupsInfo
}

type TxsStorage struct {
	data *sync.Map
}

func NewTxsStorage() *TxsStorage {
	return &TxsStorage{
		data: &sync.Map{},
	}
}

func (s *TxsStorage) AddOrUpdate(pair, input, output common.Address, pairType string, txs []types.Transaction, args []ethapi.TransactionArgs) {
	model := s.Get(pair)
	if model == nil {
		s.data.Store(pair, &TxsModel{
			Pair:     pair,
			PairType: pairType,
			Token0:   input,
			Token1:   output,
			TxsIn: models.TransactionsGroupsInfo{
				Transactions: [][]types.Transaction{
					txs,
				},
				Args: [][]ethapi.TransactionArgs{
					args,
				},
			},
		})
	} else {
		if input == model.Token0 {
			log.Info(fmt.Sprintf("updating storage, current in: %s", model.TxsIn.Hashes()))
			var conflicts bool
			for _, tx := range txs {
				conflicts = model.TxsIn.Contains(tx.Hash())
				if conflicts {
					log.Info(fmt.Sprintf("conflict: %s", tx.Hash().String()))
					break
				}
			}

			if !conflicts {
				model.TxsIn.Transactions = append(model.TxsIn.Transactions, txs)
				model.TxsIn.Args = append(model.TxsIn.Args, args)
				log.Info(fmt.Sprintf("updating storage, updated in: %s", model.TxsIn.Hashes()))
			}
		} else if input == model.Token1 {
			log.Info(fmt.Sprintf("updating storage, current out: %s", model.TxsOut.Hashes()))

			var conflicts bool
			for _, tx := range txs {
				conflicts = model.TxsOut.Contains(tx.Hash())
				if conflicts {
					log.Info(fmt.Sprintf("conflict: %s", tx.Hash().String()))
					break
				}
			}

			if !conflicts {
				model.TxsOut.Transactions = append(model.TxsOut.Transactions, txs)
				model.TxsOut.Args = append(model.TxsOut.Args, args)
				log.Info(fmt.Sprintf("updating storage, updated out: %s", model.TxsOut.Hashes()))
			}
		}

		s.data.Store(pair, model)
	}
}

func (s *TxsStorage) Get(pair common.Address) *TxsModel {
	value, ok := s.data.Load(pair)
	if !ok {
		return nil
	}
	return value.(*TxsModel)
}

func (s *TxsStorage) ListFiltered(pairs []common.Address) []*TxsModel {
	var result []*TxsModel

	for _, pair := range pairs {
		result = append(result, s.Get(pair))
	}

	return result
}

func (s *TxsStorage) List() []*TxsModel {
	var result []*TxsModel

	s.data.Range(func(key, value any) bool {
		result = append(result, value.(*TxsModel))
		return true
	})

	return result
}
