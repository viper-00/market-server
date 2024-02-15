package model

import "market/model/common"

type TelegramAllInOneKey struct {
	common.MARKET_MODEL
	TGID     string `json:"tg_id" gorm:"comment:tg_id"`
	AuthKey  string `json:"auth_key" gorm:"comment:auth_key"`
	BotToken string `json:"bot_token" gorm:"comment:bot_token"`
}

func (TelegramAllInOneKey) TableName() string {
	return "market_telegram_allinone_keys"
}
