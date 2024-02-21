package model

import "market/model/common"

type User struct {
	common.MARKET_MODEL
	ChainId         int    `json:"chain_id" gorm:"comment:chain_id"`
	PrivateKey      string `json:"-" gorm:"comment:private_key"`
	Address         string `json:"address" gorm:"comment:address"`
	ContractAddress string `json:"contract_address" gorm:"comment:contract_address"`
	Email           string `json:"email" gorm:"comment:email"`
	SuperiorId      uint   `json:"-" gorm:"comment:superior_id"`
	InvitationCode  string `json:"invitation_code" gorm:"comment:invitation_code"`
}

func (User) TableName() string {
	return "market_users"
}
