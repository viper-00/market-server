package request

type UserRegister struct {
	Email   string `json:"email" form:"email"`
	Address string `json:"address" form:"address"`
	ChainId int    `json:"chain_id" form:"chain_id" binding:"required"`
}

type UserVerifyInvitation struct {
	Code string `json:"code" form:"code" binding:"required"`
}

type UserLogin struct {
	Email   string `json:"email" form:"email"`
	Address string `json:"address" form:"address"`
	ChainId int    `json:"chain_id" form:"chain_id" binding:"required"`
}
