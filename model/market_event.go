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
	EventLogo         string    `json:"event_logo" gorm:"comment:event_logo"`
	SettlementAddress string    `json:"settlement_address" gorm:"comment:settlement_address"`
	CapitalPool       float64   `json:"capital_pool" gorm:"comment:capital_pool"`
	ResolverAddress   string    `json:"rosolver_address" gorm:"comment:rosolver_address"`
	RuleDetails       string    `json:"rule_details" gorm:"comment:rule_details"`
	EventStatus       int       `json:"event_status" gorm:"comment:event_status"`
}

func (Event) TableName() string {
	return "market_events"
}
