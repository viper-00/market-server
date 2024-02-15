package service

import (
	MARKET_Client "market/utils/http"
)

type MService struct{}

var (
	MarketService = new(MService)
	client        MARKET_Client.Client
)
