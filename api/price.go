package api

import (
	"market/global"
	"market/model/common"
	"market/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (n *MarketApi) GetCryptoPrice(c *gin.Context) {
	var res common.Response

	result, err := service.MarketService.GetCryptoPrice(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}
