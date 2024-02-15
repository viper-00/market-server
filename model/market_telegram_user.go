package model

import "market/model/common"

type TelegramUser struct {
	common.MARKET_MODEL
	TGID         string `json:"tg_id" gorm:"comment:tg_id"`
	FirstName    string `json:"first_name" gorm:"comment:first_name"`
	LastName     string `json:"last_name" gorm:"comment:last_name"`
	Username     string `json:"username" gorm:"comment:username"`
	LanguageCode string `json:"language_code" gorm:"comment:language_code"`
}

func (TelegramUser) TableName() string {
	return "market_telegram_users"
}
