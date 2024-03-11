package model

import "market/model/common"

type UserNotificationSetting struct {
	common.MARKET_MODEL
	UserId            uint `json:"user_id" gorm:"user_id"`
	EmailUpdate       bool `json:"email_update" gorm:"comment:email_update"`
	MarketUpdate      bool `json:"market_update" gorm:"comment:market_update"`
	DailyUpdate       bool `json:"daily_update" gorm:"comment:daily_update"`
	IncomingUpdate    bool `json:"incoming_update" gorm:"comment:incoming_update"`
	OutgoingUpdate    bool `json:"outgoing_update" gorm:"comment:outgoing_update"`
	EventUpdate       bool `json:"event_update" gorm:"comment:event_update"`
	OrderUpdate       bool `json:"order_update" gorm:"comment:order_update"`
	CryptoPriceUpdate bool `json:"crypto_price_update" gorm:"comment:crypto_price_update"`
}

func (UserNotificationSetting) TableName() string {
	return "market_user_notification_settings"
}
