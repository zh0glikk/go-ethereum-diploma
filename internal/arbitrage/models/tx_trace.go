package models

import "github.com/ethereum/go-ethereum/common"

type Action struct {
	CallType string         `json:"callType"`
	From     common.Address `json:"from"`
	Gas      string         `json:"gas"`
	Input    string         `json:"input"`
	To       common.Address `json:"to"`
	Value    string         `json:"value"`
}

type Result struct {
	GasUsed string `json:"gasUsed"`
	Output  string `json:"output"`
}

type TransactionTrace struct {
	Action              Action      `json:"action"`
	BlockHash           interface{} `json:"blockHash"`
	BlockNumber         int         `json:"blockNumber"`
	Error               string      `json:"error"`
	Result              Result      `json:"result"`
	Subtraces           int         `json:"subtraces"`
	TraceAddress        []int       `json:"traceAddress"`
	TransactionHash     interface{} `json:"transactionHash"`
	TransactionPosition int         `json:"transactionPosition"`
	Type                string      `json:"type"`
}
