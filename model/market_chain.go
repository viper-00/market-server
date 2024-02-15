package model

import "market/model/common"

var ChainList []ChainInfo

type Coin struct {
	Symbol     string `json:"symbol"`
	Decimals   int    `json:"decimals"`
	Contract   string `json:"contract"`
	IsMainCoin bool   `json:"isMainCoin"`
}

type ChainInfo struct {
	Name      string `json:"name"`
	Chain     string `json:"chain"`
	ChainId   int    `json:"chainId"`
	NetworkId int    `json:"networkId"`
	Coins     []Coin `json:"coins"`
}

type Chain struct {
	common.MARKET_MODEL
	Name       string `json:"name" gorm:"comment:name"`
	Chain      string `json:"chain" gorm:"comment:chain"`
	ChainId    int    `json:"chain_id" gorm:"comment:chain_id"`
	NetworkId  int    `json:"network_id" gorm:"comment:network_id"`
	Symbol     string `json:"symbol" gorm:"comment:symbol"`
	Decimals   int    `json:"decimals" gorm:"comment:decimals"`
	Contract   string `json:"contract" gorm:"comment:contract"`
	IsMainCoin bool   `json:"is_main_coin" gorm:"comment:is_main_coin"`
}

func (Chain) TableName() string {
	return "market_chains"
}
