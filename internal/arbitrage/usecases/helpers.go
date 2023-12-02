package usecases

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/internal/arbitrage/models"
	"github.com/ethereum/go-ethereum/internal/ethapi"
)

func convertCallManyResult(result *core.ExecutionResult) (models.CallManyResponseDTO, bool) {
	var resp models.CallManyResponseDTO

	revert := result.Revert()
	if len(revert) > 0 {
		resp = models.CallManyResponseDTO{
			Error: models.NewRevertErrorWithDecoding(result).ErrorData(),
		}
		return resp, false
	}
	if result.Unwrap() != nil {
		resp = models.CallManyResponseDTO{
			Error: result.Unwrap().Error(),
		}
		return resp, false
	}

	return models.CallManyResponseDTO{
		Value: hexutil.Bytes(result.Return()),
	}, true
}

func convertCallManyResults(results []*core.ExecutionResult) []models.CallManyResponseDTO {
	resp := make([]models.CallManyResponseDTO, len(results))

	for ind := range results {
		resp[ind], _ = convertCallManyResult(results[ind])
	}

	return resp
}

func convertTraceCallManyResult(traceResponse []interface{}) ([][]models.TransactionTrace, error) {
	var traces [][]models.TransactionTrace
	for _, tx := range traceResponse {
		bb, err := json.Marshal(tx)
		if err != nil {
			return nil, err
		}

		var txs []models.TransactionTrace
		err = json.Unmarshal(bb, &txs)
		if err != nil {
			return nil, err
		}

		traces = append(traces, txs)
	}
	return traces, nil
}

func applyContractCode(stateOverride *ethapi.StateOverride, contract common.Address, code *hexutil.Bytes) {
	if stateOverride == nil {
		stateOverride = &ethapi.StateOverride{}
	}

	(*stateOverride)[contract] = ethapi.OverrideAccount{
		Code: code,
	}
}

func convertCallManyResultWithDecoding(result *core.ExecutionResult) (models.CallManyResponseDTO, bool) {
	var resp models.CallManyResponseDTO

	revert := result.Revert()
	if len(revert) > 0 {
		resp = models.CallManyResponseDTO{
			Error: models.NewRevertErrorWithDecoding(result).ErrorData(),
		}
		return resp, false
	}
	if result.Unwrap() != nil {
		resp = models.CallManyResponseDTO{
			Error: result.Unwrap().Error(),
		}
		return resp, false
	}

	return models.CallManyResponseDTO{
		Value: hexutil.Bytes(result.Return()),
	}, true
}

func convertCallManyResultsWithDecoding(results []*core.ExecutionResult) []models.CallManyResponseDTO {
	resp := make([]models.CallManyResponseDTO, len(results))

	for ind := range results {
		resp[ind], _ = convertCallManyResultWithDecoding(results[ind])
	}

	return resp
}
