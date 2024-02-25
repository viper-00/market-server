package model

import "market/model/common"

type UserAffiliate struct {
	common.MARKET_MODEL
	UserId uint `json:"user_id" gorm:"user_id"`
}

func (UserAffiliate) TableName() string {
	return "market_user_affiliates"
}
