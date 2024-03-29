package request

type GetMarketEvent struct {
	Code string `json:"code" form:"code" binding:"required"`
}

type CreateMarketEvent struct {
	Title           string `json:"title" form:"title" binding:"required"`
	ExpireTime      int64  `json:"expire_time" form:"expire_time" binding:"required"`
	Type            string `json:"type" form:"type" binding:"required"`
	PlayType        string `json:"play_type" from:"play_type" binding:"required"`
	EventLogo       string `json:"event_logo" form:"event_logo" binding:"required"`
	ResolverAddress string `json:"resolver_address" form:"resolver_address" binding:"required"`
	Password        string `json:"password" form:"password" binding:"required"`
}

type UpdateMarketEvent struct {
	Code      string `json:"code" form:"code" binding:"required"`
	Title     string `json:"title" form:"title" binding:"required"`
	Type      string `json:"type" form:"type" binding:"required"`
	EventLogo string `json:"event_logo" form:"event_logo" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required"`
}

type CreateMarketEventPlay struct{}

type UpdateMarketEventPlay struct{}

type CreateMarketEventOrder struct {
	EventUniqueCode string  `json:"event_unique_code" form:"event_unique_code" binding:"required"`
	Amount          float64 `json:"amount" form:"amount"`
	PlayValue       string  `json:"play_value" form:"play_value" binding:"required"`
	Type            uint    `json:"type" form:"type" binding:"required"`
}

type SettleMarketOrder struct {
	EventUniqueCode string `json:"event_unique_code" form:"event_unique_code" binding:"required"`
	Password        string `json:"password" form:"password" binding:"required"`
}

type GetMarketEventForHome struct {
	Type string `json:"type" form:"type" binding:"required"`
}
