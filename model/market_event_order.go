package model

import "market/model/common"

type EventOrder struct {
	common.MARKET_MODEL
	UserId      uint    `json:"-" gorm:"comment:user_id"`
	EventId     uint    `json:"-" gorm:"event_id"`
	Amount      float64 `json:"amount" gorm:"amount"`
	PlayValue   string  `json:"play_value" gorm:"play_value"`
	OrderStatus uint    `json:"order_status" gorm:"order_status"` // 1:process 2:settlement 3:complete 4:failed
	OrderType   uint    `json:"order_type" gorm:"order_type"`     // 1: buy 2: sell
}

func (EventOrder) TableName() string {
	return "market_event_orders"
}
