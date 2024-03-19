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
	ChainId int    `json:"chain_id" form:"chain_id"`
	Code    string `json:"code" form:"code"`
}

type UpdateUserSetting struct {
	Username  string `json:"username" form:"username"`
	AvatarUrl string `json:"avatar_url" form:"avatar_url"`
	Bio       string `json:"bio" form:"bio"`
}

type GetUserProfile struct {
	Address string `json:"address" form:"address" binding:"required"`
}

type UpdateUserNotificationSetting struct {
	// EmailUpdate       bool `json:"email_update"`
	// MarketUpdate      bool `json:"market_update"`
	// DailyUpdate       bool `json:"daily_update"`
	// IncomingUpdate    bool `json:"incoming_update"`
	// OutgoingUpdate    bool `json:"outgoing_update"`
	// EventUpdate       bool `json:"event_update"`
	// OrderUpdate       bool `json:"order_update"`
	// CryptoPriceUpdate bool `json:"crypto_price_update"`
	Type   string `json:"type"`
	Status int    `json:"status"`
}

type CreateUserAffiliate struct{}
