package model

import "market/model/common"

type Wallet struct {
	common.MARKET_MODEL
	Address     string `json:"address" gorm:"comment:address"`
	ChainId     int    `json:"chain_id" gorm:"comment:chain_id"`
	NetworkName string `json:"network_name" gorm:"comment:network_name"`
}

func (Wallet) TableName() string {
	return "market_wallets"
}
