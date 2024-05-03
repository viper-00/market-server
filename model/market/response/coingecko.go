package response

type CoingeckoPrice struct {
	Ethereum CoingeckoPriceCore `json:"ethereum"`
	USDT     CoingeckoPriceCore `json:"tether"`
	USDC     CoingeckoPriceCore `json:"usd-coin"`
}

type CoingeckoPriceCore struct {
	USD           float64 `json:"usd"`
	USDMarketCap  float64 `json:"usd_market_cap"`
	USD24hVol     float64 `json:"usd_24h_vol"`
	USD24hChange  float64 `json:"usd_24h_change"`
	LastUpdatedAt int     `json:"last_updated_at"`
}

type CustomCoingeckoPrice struct {
	ETH  CoingeckoPriceCore `json:"eth"`
	USDT CoingeckoPriceCore `json:"usdt"`
	USDC CoingeckoPriceCore `json:"usdc"`
}
