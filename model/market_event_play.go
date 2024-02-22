package model

import "market/model/common"

type EventPlay struct {
	common.MARKET_MODEL
	UserId         uint    `json:"user_id" gorm:"comment:user_id"`
	Introduce      string  `json:"introduce" gorm:"comment:introduce"`
	InitialCapital float64 `json:"initial_capital" gorm:"comment:initial_capital"`
	Coin           string  `json:"coin" gorm:"comment:coin"`
}

func (EventPlay) TableName() string {
	return "market_event_plays"
}
