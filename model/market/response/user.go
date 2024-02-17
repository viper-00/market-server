package response

import "market/model/common"

type User struct {
	common.MARKET_MODEL
	ChainId         int    `json:"chain_id" gorm:"comment:chain_id"`
	PrivateKey      string `json:"-" gorm:"comment:private_key"`
	Address         string `json:"address" gorm:"comment:address"`
	ContractAddress string `json:"contract_address" gorm:"comment:contract_address"`
	Email           string `json:"email" gorm:"comment:email"`
	Auth            string `json:"auth" gorm:"comment:auth"`
}
