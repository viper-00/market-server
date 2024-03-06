package model

import (
	"market/model/common"
	"time"
)

type Event struct {
	common.MARKET_MODEL
	Title             string    `json:"title" gorm:"comment:title"`
	UniqueWebsiteLink string    `json:"unique_website_link" gorm:"comment:unique_website_link"`
	ExpireTime        time.Time `json:"expire_time" gorm:"comment:expire_time"`
	Type              string    `json:"type" gorm:"comment:type"`
	PlayId            uint      `json:"play_id" gorm:"comment:play_id"`
	EventLogo         string    `json:"event_logo" gorm:"comment:event_logo"`
	SettlementAddress string    `json:"settlement_address" gorm:"comment:settlement_address"`
	ResolverAddress   string    `json:"rosolver_address" gorm:"comment:rosolver_address"`
	EventStatus       int       `json:"event_status" gorm:"comment:event_status"`
	Password          string    `json:"password" gorm:"comment:password"`
}

func (Event) TableName() string {
	return "market_events"
}
