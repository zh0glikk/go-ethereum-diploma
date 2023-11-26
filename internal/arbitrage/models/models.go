package models

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"strings"
)

type CallManyResponseDTO struct {
	Value interface{} `json:"value,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

type ChanCallManyResponse struct {
	Response []CallManyResponseDTO
	Ind      int
}

type ChanCallManyError struct {
	Err error
	Ind int
}

func NewRevertError(result *core.ExecutionResult) *RevertError {
	reason, errUnpack := abi.UnpackRevert(result.Revert())
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &RevertError{
		error:  err,
		reason: hexutil.Encode(result.Revert()),
	}
}

type RevertError struct {
	error
	reason string // revert reason hex encoded
}

func (e *RevertError) ErrorCode() int {
	return 3
}

func (e *RevertError) ErrorData() interface{} {
	return e.reason
}

func NewRevertErrorWithDecoding(result *core.ExecutionResult) *RevertError {
	reason, errUnpack := abi.UnpackRevert(result.Revert())
	err := errors.New("execution reverted")

	reason = strings.Replace(reason, "\"", "'", -1)
	reason = strings.Replace(reason, "{", "(", -1)
	reason = strings.Replace(reason, "}", ")", -1)

	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &RevertError{
		error:  err,
		reason: reason,
	}
}
