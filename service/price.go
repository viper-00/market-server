package service

import (
	"context"
	"encoding/json"
	"market/global"
	"market/global/constant"
	"market/model/market/response"

	"github.com/gin-gonic/gin"
)

func (m *MService) GetCryptoPrice(c *gin.Context) (result interface{}, err error) {
	var price response.CustomCoingeckoPrice

	data, err := global.MARKET_REDIS.Get(context.Background(), constant.CRYPTO_PRICE).Result()
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	err = json.Unmarshal([]byte(data), &price)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return price, nil
}
