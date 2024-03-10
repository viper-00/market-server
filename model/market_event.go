package model

import (
	"market/model/common"
	"time"
)

type Event struct {
	common.MARKET_MODEL
	UserId            uint      `json:"-" gorm:"comment:user_id"`
	Title             string    `json:"title" gorm:"comment:title"`
	UniqueWebsiteCode string    `json:"unique_website_code" gorm:"comment:unique_website_code"`
	ExpireTime        time.Time `json:"expire_time" gorm:"comment:expire_time"`
	Type              string    `json:"type" gorm:"comment:type"`
	PlayId            uint      `json:"play_id" gorm:"comment:play_id"`
	EventLogo         string    `json:"event_logo" gorm:"comment:event_logo"`
	SettlementAddress string    `json:"settlement_address" gorm:"comment:settlement_address"`
	ResolverAddress   string    `json:"rosolver_address" gorm:"comment:rosolver_address"`
	EventStatus       int       `json:"event_status" gorm:"comment:event_status"` //1: normal
	Password          string    `json:"-" gorm:"comment:password"`
}

func (Event) TableName() string {
	return "market_events"
}
