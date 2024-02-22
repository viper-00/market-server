package model

import "market/model/common"

type UserNotificationSetting struct {
	common.MARKET_MODEL
	UserId       uint   `json:"user_id" gorm:"user_id"`
	Email        string `json:"email" gorm:"comment:email"`
	MarketUpdate bool   `json:"market_update" gorm:"comment:market_update"`
}

func (UserNotificationSetting) TableName() string {
	return "market_user_notification_settings"
}
