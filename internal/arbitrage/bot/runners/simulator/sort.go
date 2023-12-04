package simulator

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"time"
)

type TxWrapper struct {
	TxArgs []ethapi.TransactionArgs
	Tx     []types.Transaction
	Profit *big.Int
}

// TODO: refactor somehow
// func SortSandwich2Process(
// 	ctx context.Context,
// 	b ethapi.Backend,
// 	s Sandwich2Process,
// 	contract common.Address,
// 	initialSplitParam *big.Int,
// 	slippageEstimatorWrapper simulation_wrappers.SlippageEstimatorWrapper,
// 	prebundleSlippageEstimatorWrapper simulation_wrappers.SlippageEstimatorWrapper,
// 	blockNrOrHash rpc.BlockNumberOrHash,
// 	blockOverrides *ethapi.BlockOverrides,
// ) ([]Sandwich2Process, error) {
// 	log.Info("start sandwich2Process sort")
// 	if len(s.prebundleTxs.Transactions) > 1 {
// 		log.Info(fmt.Sprintf("prebundle sorting, prebudnle len: %d", len(s.prebundleTxs.Transactions)))
//
// 		prebundleTxsWithProfit := []TxWrapper{}
// 		for i := range s.prebundleTxs.Args {
//
// 			min, max, err := usecases.GetMinMaxTokens(ctx, b, contract, s.output, s.pair, nil, blockNrOrHash, blockOverrides)
// 			if err != nil {
// 				return nil, err
// 			}
//
// 			points, _, _ := algo.PrepareNewPointsWithMerge(
// 				min, max,
// 				initialSplitParam,
// 				nil,
// 				false)
//
// 			data := models.Data2Simulate{
// 				SandwichContractParams: models.SandwichContractParams{
// 					Input:  s.output,
// 					Output: s.input,
// 					Pair:   s.pair,
// 				}, TxsArgs: models.TxsArgs{
// 					Transactions: s.prebundleTxs.Args[i],
// 				},
// 				BlockNumberOrHash: blockNrOrHash,
// 				BlockOverrides:    blockOverrides,
// 				Points:            points,
// 			}
// 			resp, err := prebundleSlippageEstimatorWrapper.EstimateSlippage(ctx, data)
// 			if err != nil {
// 				return nil, err
// 			}
// 			log.Info(fmt.Sprintf("Prebundletxs: %s, slippage: %s\nreason: %s", s.prebundleTxs.GroupHashes(i), resp.Profit.String(), resp.Reason))
//
// 			prebundleTxsWithProfit = append(prebundleTxsWithProfit, TxWrapper{Tx: s.prebundleTxs.Transactions[i], TxArgs: s.prebundleTxs.Args[i], Profit: resp.Profit})
// 		}
// 		log.Info(fmt.Sprintf("%v", len(prebundleTxsWithProfit)))
// 		for i, tx := range prebundleTxsWithProfit {
// 			for _, t := range tx.Tx {
// 				log.Info(fmt.Sprintf("PrebundleTxWrapper %v profit - %v", i, tx.Profit), fmt.Sprintf("%v", t.Hash()))
// 			}
// 		}
// 		s.prebundleTxs = bubbleSort(prebundleTxsWithProfit, true)
// 		for _, tx := range s.prebundleTxs.Transactions {
// 			for _, t := range tx {
// 				log.Info("prebundle sort", fmt.Sprintf("%v", t.Hash()))
// 			}
// 		}
// 	}
//
// 	if len(s.victimTxs.Args) <= 1 {
// 		log.Info(fmt.Sprintf("victim sorting case 1, victim len: %d", len(s.victimTxs.Transactions)))
//
// 		return []Sandwich2Process{s}, nil
// 	}
//
// 	if len(s.victimTxs.Args) == 2 && len(s.victimTxs.Transactions) == 2 {
// 		log.Info(fmt.Sprintf("victim sorting case 2, victim len: %d", len(s.victimTxs.Transactions)))
//
// 		ascSlippageVictim := &Sandwich2Process{
// 			input:        s.input,
// 			output:       s.output,
// 			prebundleTxs: s.prebundleTxs,
// 			victimTxs: models.TransactionsGroupsInfo{
// 				Transactions: [][]types.Transaction{s.victimTxs.Transactions[0], s.victimTxs.Transactions[1]},
// 				Args:         [][]ethapi.TransactionArgs{s.victimTxs.Args[0], s.victimTxs.Args[1]},
// 			},
// 			action: s.action,
// 		}
//
// 		descSlippageVictim := &Sandwich2Process{
// 			input:        s.input,
// 			output:       s.output,
// 			prebundleTxs: s.prebundleTxs,
// 			victimTxs: models.TransactionsGroupsInfo{
// 				Transactions: [][]types.Transaction{s.victimTxs.Transactions[1], s.victimTxs.Transactions[0]},
// 				Args:         [][]ethapi.TransactionArgs{s.victimTxs.Args[1], s.victimTxs.Args[0]},
// 			},
// 			action: s.action,
// 		}
//
// 		return []Sandwich2Process{*ascSlippageVictim, *descSlippageVictim}, nil
// 	}
//
// 	log.Info(fmt.Sprintf("victim sorting case 3, victim len: %d\n victims: %s", len(s.victimTxs.Transactions), s.victimTxs.Hashes()))
//
// 	var ascSlippageVictim = Sandwich2Process{
// 		input:        s.input,
// 		output:       s.output,
// 		prebundleTxs: s.prebundleTxs,
// 		action:       s.action,
// 	}
// 	var descSlippageVictim = Sandwich2Process{
// 		input:        s.input,
// 		output:       s.output,
// 		prebundleTxs: s.prebundleTxs,
// 		action:       s.action,
// 	}
// 	var ascAmountVictim = Sandwich2Process{
// 		input:        s.input,
// 		output:       s.output,
// 		prebundleTxs: s.prebundleTxs,
// 		action:       s.action,
// 	}
// 	var descAmountVictim = Sandwich2Process{
// 		input:        s.input,
// 		output:       s.output,
// 		prebundleTxs: s.prebundleTxs,
// 		action:       s.action,
// 	}
//
// 	victimTxsWithProfit := []TxWrapper{}
// 	log.Info("estimating slippage for txs")
// 	for i := range s.victimTxs.Args {
//
// 		min, max, err := usecases.GetMinMaxTokens(ctx, b, contract, s.input, s.pair, nil, blockNrOrHash, blockOverrides)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		points, _, _ := algo.PrepareNewPointsWithMerge(
// 			min, max,
// 			initialSplitParam,
// 			nil,
// 			false)
//
// 		data := models.Data2Simulate{
// 			SandwichContractParams: models.SandwichContractParams{
// 				Input:  s.input,
// 				Output: s.output,
// 				Pair:   s.pair,
// 			}, TxsArgs: models.TxsArgs{
// 				Prebundle:    s.prebundleTxs.ArgsList(),
// 				Transactions: s.victimTxs.Args[i],
// 			},
// 			BlockNumberOrHash: blockNrOrHash,
// 			BlockOverrides:    blockOverrides,
// 			Points:            points,
// 		}
// 		log.Info(fmt.Sprintf("txs: %s\n Input - %v\nOutput- %v\nPair - %v\nMethod - %v",
// 			s.victimTxs.GroupHashes(i), data.Input, data.Output, data.Pair, slippageEstimatorWrapper.Method()))
// 		resp, err := slippageEstimatorWrapper.EstimateSlippage(ctx, data)
// 		if err != nil {
// 			return nil, err
// 		}
// 		log.Info(fmt.Sprintf("txs: %s, slippage: %s\nreason: %s", s.victimTxs.GroupHashes(i), resp.Profit.String(), resp.Reason))
//
// 		victimTxsWithProfit = append(victimTxsWithProfit, TxWrapper{Tx: s.victimTxs.Transactions[i], TxArgs: s.victimTxs.Args[i], Profit: resp.Profit})
// 	}
//
// 	victimTxsWithInputAmount := []TxWrapper{}
// 	log.Info("estimating input amount for txs")
// 	for i := range s.victimTxs.Args {
// 		amount, err := usecases.EstimateInputAmount(ctx, b,
// 			models.EstimateInputAmountBundle{
// 				Transactions: s.victimTxs.Args[i],
// 				Pair:         s.pair,
// 				InputAddress: s.input,
// 			},
// 			blockNrOrHash,
// 			&tracers.TraceCallConfig{
// 				BlockOverrides: blockOverrides,
// 			},
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		log.Info(fmt.Sprintf("txs: %s, amount: %s", s.victimTxs.GroupHashes(i), amount.TotalAmount.String()))
//
// 		victimTxsWithInputAmount = append(victimTxsWithInputAmount, TxWrapper{Tx: s.victimTxs.Transactions[i], TxArgs: s.victimTxs.Args[i], Profit: amount.TotalAmount})
// 	}
//
// 	ascSlippageVictim.victimTxs = bubbleSort(victimTxsWithProfit, true)
// 	descSlippageVictim.victimTxs = bubbleSort(victimTxsWithProfit, false)
// 	ascAmountVictim.victimTxs = bubbleSort(victimTxsWithInputAmount, true)
// 	descAmountVictim.victimTxs = bubbleSort(victimTxsWithInputAmount, false)
//
// 	tmp := []Sandwich2Process{ascSlippageVictim, descSlippageVictim, ascAmountVictim, descAmountVictim}
// 	result := []Sandwich2Process{}
// 	hashes := make(map[string]bool)
// 	for _, sandwich := range tmp {
// 		if _, ok := hashes[sandwich.victimTxs.Hashes()]; !ok {
// 			hashes[sandwich.victimTxs.Hashes()] = true
// 			result = append(result, sandwich)
// 		} else {
// 			log.Info(fmt.Sprintf("collision found after sorting, skipped: %s", sandwich.victimTxs.Hashes()))
// 		}
// 	}
//
// 	log.Info(fmt.Sprintf("case 3 returning %v ", len(result)))
//
// 	return result, nil
// }

func bubbleSort(txs []TxWrapper, asc bool) models.TransactionsGroupsInfo {
	startedAt := time.Now().UTC()
	for i := 0; i < len(txs)-1; i++ {
		for j := 0; j < len(txs)-i-1; j++ {
			if asc && txs[j].Profit.Cmp(txs[j+1].Profit) == 1 {
				txs[j], txs[j+1] = txs[j+1], txs[j]
			}
			if !asc && txs[j].Profit.Cmp(txs[j+1].Profit) == -1 {
				txs[j], txs[j+1] = txs[j+1], txs[j]
			}
		}
	}

	sortedArgs := [][]ethapi.TransactionArgs{}
	sortedTxs := [][]types.Transaction{}
	for _, tx := range txs {
		fmt.Println(tx.Profit)
		sortedArgs = append(sortedArgs, tx.TxArgs)
		sortedTxs = append(sortedTxs, tx.Tx)
	}

	result := models.TransactionsGroupsInfo{
		Transactions: sortedTxs,
		Args:         sortedArgs,
	}

	log.Info(fmt.Sprintf("Bubble sort: txsGroupsCount %v in %v\ntxs: %s\n",
		len(sortedArgs),
		time.Now().UTC().Sub(startedAt).Milliseconds(),
		result.Hashes(),
	))

	return result
}
