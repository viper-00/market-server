package model

import "market/model/common"

type User struct {
	common.MARKET_MODEL
	ChainId         int    `json:"-" gorm:"comment:chain_id"`
	PrivateKey      string `json:"-" gorm:"comment:private_key"`
	Address         string `json:"-" gorm:"comment:address"`
	ContractAddress string `json:"contract_address" gorm:"comment:contract_address"`
	Email           string `json:"email" gorm:"comment:email"`
	SuperiorId      uint   `json:"-" gorm:"comment:superior_id"`
	InvitationCode  string `json:"invitation_code" gorm:"comment:invitation_code"`
	Way             uint   `json:"way" gorm:"way"` // 1: email 2: wallet
}

func (User) TableName() string {
	return "market_users"
}
