package model

import "market/model/common"

type EventPlay struct {
	common.MARKET_MODEL
	UserId             uint    `json:"user_id" gorm:"comment:user_id"`
	Title              string  `json:"title" gorm:"comment:title"`
	Introduce          string  `json:"introduce" gorm:"type:text;comment:introduce"`
	GuessNumber        int     `json:"guess_number" gorm:"comment:guess_number"`
	MinimumCapitalPool float64 `json:"minimum_capital_pool" gorm:"comment:minimum_capital_pool"`
	MaximumCapitalPool float64 `json:"maximum_capital_pool" gorm:"comment:maximum_capital_pool"`
	Coin               string  `json:"coin" gorm:"comment:coin"` // USDT USDC
	PledgeAmount       float64 `json:"pledge_amount" gorm:"comment:pledge_amount"`
}

func (EventPlay) TableName() string {
	return "market_event_plays"
}
