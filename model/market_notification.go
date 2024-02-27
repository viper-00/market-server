package model

import (
	"market/model/common"
)

type Notification struct {
	common.MARKET_MODEL
	Hash             string `json:"hash" gorm:"comment:hash"`
	ChainId          int    `json:"chain_id" gorm:"comment:chain_id"`
	UserId           uint   `json:"user_id" gorm:"user_id"`
	Title            string `json:"title" gorm:"title"`
	Description      string `json:"description" gorm:"description"`
	Content          string `json:"content" gorm:"content"`
	NotificationType string `json:"notification_type" gorm:"notification_type"` // incoming outgoing system other
	OwnId            uint   `json:"own_id" gorm:"own_id"`
	IsRead           int    `json:"is_read" gorm:"is_read"` // 1:read 2:unread 3:delete
}

func (Notification) TableName() string {
	return "market_notifications"
}
