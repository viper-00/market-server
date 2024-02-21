package model

import "market/model/common"

type Affiliate struct {
	common.MARKET_MODEL
	ChainId int `json:"chain_id" gorm:"comment:chain_id"`
}

func (Affiliate) TableName() string {
	return "market_user_affiliates"
}
