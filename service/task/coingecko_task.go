package task

import (
	"context"
	"encoding/json"
	"fmt"
	"market/global"
	"market/global/constant"
	"market/model/market/response"
	"time"
)

func RunCoingeckoTask() {
	for {
		now := time.Now()

		nextSecond := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+10, 0, now.Location())
		durationUntilNextHour := nextSecond.Sub(now)

		ticker := time.NewTicker(durationUntilNextHour)

		<-ticker.C

		RunCoingeckoCore()
	}
}

func RunCoingeckoCore() {
	global.MARKET_LOG.Info("---------- Run Coingecko Task ----------")

	ids := fmt.Sprintf("%s,%s,%s", constant.IDS_Ethereum, constant.IDS_USDT, constant.IDS_USDC)
	currency := "USD"
	include_market_cap_string := "true"
	include_24hr_vol_string := "true"
	include_24hr_change_string := "true"
	include_last_updated_at_string := "true"

	client.URL = fmt.Sprintf("%s?ids=%s&vs_currencies=%s&include_market_cap=%s&include_24hr_vol=%s&include_24hr_change=%s&include_last_updated_at=%s", constant.CoingeckoGetPrice, ids, currency, include_market_cap_string, include_24hr_vol_string, include_24hr_change_string, include_last_updated_at_string)
	client.Headers = map[string]string{
		"x_cg_demo_api_key": global.MARKET_CONFIG.Coingecko.ApiKey,
	}

	var cryptoResponse response.CoingeckoPrice
	err := client.HTTPGet(&cryptoResponse)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	var customCoingeckoPrice response.CustomCoingeckoPrice
	customCoingeckoPrice.ETH = cryptoResponse.Ethereum
	customCoingeckoPrice.USDT = cryptoResponse.USDT
	customCoingeckoPrice.USDC = cryptoResponse.USDC

	cryptoByte, err := json.Marshal(customCoingeckoPrice)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	_, err = global.MARKET_REDIS.Set(context.Background(), constant.CRYPTO_PRICE, string(cryptoByte), time.Minute*10).Result()
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}
}
