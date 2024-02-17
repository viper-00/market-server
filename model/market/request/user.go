package request

type UserRegister struct {
	Email   string `json:"email" form:"email"`
	Address string `json:"address" form:"address"`
	ChainId int    `json:"chain_id" form:"chain_id"`
}

type UserVerifyInvitation struct {
	InvitationCode string `json:"invitation_code" form:"invitation_code" binding:"required"`
	Email          string `json:"email" form:"email" binding:"required"`
}

type UserLogin struct {
	Email   string `json:"email" form:"email"`
	Address string `json:"address" form:"address"`
	ChainId int    `json:"chain_id" form:"chain_id"`
}
