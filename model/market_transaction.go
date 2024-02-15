package model

import "market/model/common"

type Transaction struct {
	common.MARKET_MODEL
	ChainId              int    `json:"chain_id" gorm:"comment:chain_id"`
	Hash                 string `json:"hash" gorm:"comment:hash"`
	BlockNumber          int64  `json:"block_number" gorm:"comment:block_number"`
	BlockHash            string `json:"block_hash" gorm:"comment:block_hash"`
	From                 string `json:"from" gorm:"comment:from"`
	To                   string `json:"to" gorm:"comment:to"`
	Gas                  int64  `json:"gas" gorm:"comment:gas"`
	GasPrice             int64  `json:"gasPrice" gorm:"comment:gasPrice"`
	Input                string `json:"input" gorm:"type:longtext;size:65535;comment:input"`
	MaxFeePerGas         int64  `json:"maxFeePerGas" gorm:"comment:maxFeePerGas"`
	MaxPriorityFeePerGas int64  `json:"maxPriorityFeePerGas" gorm:"comment:maxPriorityFeePerGas"`
	Nonce                int64  `json:"nonce" gorm:"comment:nonce"`
	TransactionIndex     int64  `json:"transactionIndex" gorm:"comment:transactionIndex"`
	Type                 int64  `json:"type" gorm:"comment:type"`
	Value                int64  `json:"value" gorm:"comment:value"`
	BlockTimestamp       int    `json:"block_timestamp" gorm:"comment:block_timestamp"`
}

func (Transaction) TableName() string {
	return "market_transactions"
}
