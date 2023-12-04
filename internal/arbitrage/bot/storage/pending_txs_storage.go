package storage

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/arbitrage/utils"
	"sync"
	"time"
)

type PendingTxsModel struct {
	From         common.Address
	CreatedAt    time.Time
	Transactions []types.Transaction
}

type PendingTxsStorage struct {
	data *sync.Map
}

func NewPendingTxsStorage() *PendingTxsStorage {
	return &PendingTxsStorage{
		data: &sync.Map{},
	}
}

func (s *PendingTxsStorage) Add(from common.Address, tx types.Transaction) {
	pendingTxs := s.Get(from)
	if pendingTxs == nil {
		s.data.Store(from, &PendingTxsModel{
			From:      from,
			CreatedAt: time.Now().UTC(),
			Transactions: []types.Transaction{
				tx,
			},
		})
		return
	}

	pendingTxs.Transactions = append(pendingTxs.Transactions, tx)
	pendingTxs.Transactions = sortTxsByNonce(uniqueTxs(pendingTxs.Transactions))

	s.data.Store(from, pendingTxs)
	return
}

func (s *PendingTxsStorage) Delete(from common.Address, hash common.Hash) {
	pendingTxs := s.Get(from)
	if pendingTxs == nil {
		return
	}

	for i, tx := range pendingTxs.Transactions {
		if tx.Hash() == hash {
			// log.Info(fmt.Sprintf("deleting tx with hash: %s from %s", hash.String(), pendingTxs.From.String()))
			pendingTxs.Transactions = utils.RemoveFromArray(pendingTxs.Transactions, i)
			break
		}
	}
	// if empty transactions list - we remove map entry
	if len(pendingTxs.Transactions) == 0 {
		s.data.Delete(from)
		return
	}

	s.data.Store(from, pendingTxs)

	return
}

func (s *PendingTxsStorage) Get(from common.Address) *PendingTxsModel {
	value, ok := s.data.Load(from)
	if !ok {
		return nil
	}

	return value.(*PendingTxsModel)
}

func (s *PendingTxsStorage) List() []*PendingTxsModel {
	var result []*PendingTxsModel

	s.data.Range(func(key, value any) bool {
		result = append(result, value.(*PendingTxsModel))
		return true
	})

	return result
}

func (s *PendingTxsStorage) ClearExpired() uint64 {
	clearedCount := uint64(0)

	s.data.Range(func(key, value any) bool {
		if value.(*PendingTxsModel).CreatedAt.After(time.Now().UTC().Add(-time.Minute * 60)) {
			s.data.Delete(key)
			clearedCount += 1
		}
		return true
	})

	return clearedCount
}

func (s *PendingTxsStorage) CountOfGroups() uint64 {
	var result uint64

	s.data.Range(func(key, value any) bool {
		result += 1
		return true
	})

	return result
}

func (s *PendingTxsStorage) CountOfTxs() uint64 {
	var result uint64

	s.data.Range(func(key, value any) bool {
		pendingTxs := s.Get(key.(common.Address))
		result += uint64(len(pendingTxs.Transactions))

		return true
	})

	return result
}
