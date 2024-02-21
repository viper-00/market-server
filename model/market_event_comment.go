package model

import (
	"market/model/common"
)

type EventComment struct {
	common.MARKET_MODEL
}

func (EventComment) TableName() string {
	return "market_event_comments"
}
