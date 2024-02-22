package model

import (
	"market/model/common"
)

type EventComment struct {
	common.MARKET_MODEL
	UserId  uint   `json:"user_id" gorm:"comment:user_id"`
	Content string `json:"content" gorm:"comment:content"`
	ReplyId uint   `json:"reply_id" gorm:"comment:reply_id"`
}

func (EventComment) TableName() string {
	return "market_event_comments"
}
