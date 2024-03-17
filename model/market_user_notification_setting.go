package model

import "market/model/common"

type UserNotificationSetting struct {
	common.MARKET_MODEL
	UserId            uint `json:"user_id" gorm:"user_id"`
	EmailUpdate       int  `json:"email_update" gorm:"comment:email_update"`
	DailyUpdate       int  `json:"daily_update" gorm:"comment:daily_update"`
	IncomingUpdate    int  `json:"incoming_update" gorm:"comment:incoming_update"`
	OutgoingUpdate    int  `json:"outgoing_update" gorm:"comment:outgoing_update"`
	EventUpdate       int  `json:"event_update" gorm:"comment:event_update"`
	OrderUpdate       int  `json:"order_update" gorm:"comment:order_update"`
	CryptoPriceUpdate int  `json:"crypto_price_update" gorm:"comment:crypto_price_update"` // 1: on 2: off
}

func (UserNotificationSetting) TableName() string {
	return "market_user_notification_settings"
}
