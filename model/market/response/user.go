package response

import (
	"market/model/common"
)

type User struct {
	common.MARKET_MODEL
	ChainId         int    `json:"chain_id" gorm:"comment:chain_id"`
	Address         string `json:"address" gorm:"comment:address"`
	ContractAddress string `json:"contract_address" gorm:"comment:contract_address"`
	Email           string `json:"email" gorm:"comment:email"`
	Auth            string `json:"auth" gorm:"comment:auth"`
	UserName        string `json:"username" gorm:"comment:username"`
	InviteCode      string `json:"invite_code" gorm:"comment:invite_code"`
	AvatarUrl       string `json:"avatar_url" gorm:"comment:avatar_url"`
	Bio             string `json:"bio" gorm:"comment:bio"`
	JoinedDate      int64  `json:"joined_date" gorm:"comment:joined_date"`
}

type UserBalance struct {
	ETH  string `json:"eth"`
	USDT string `json:"usdt"`
	USDC string `json:"usdc"`
}
