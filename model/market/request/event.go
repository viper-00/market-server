package request

type CreateMarketEvent struct {
	Title             string `json:"title" form:"title" binding:"required"`
	ExpireTime        int64  `json:"expire_time" form:"expire_time" binding:"required"`
	Type              string `json:"type" form:"type" binding:"required"`
	PlayType          string `json:"play_type" from:"play_type" binding:"required"`
	EventLogo         string `json:"event_logo" form:"event_logo" binding:"required"`
	SettlementAddress string `json:"settlement_address" form:"settlement_address" binding:"required"`
	ResolverAddress   string `json:"resolver_address" form:"resolver_address" binding:"required"`
	Password          string `json:"password" form:"password" binding:"required"`
}

type UpdateMarketEvent struct{}

type CreateMarketEventPlay struct{}

type UpdateMarketEventPlay struct{}
