package api

import (
	"market/global"
	"market/model/common"
	"market/model/market/request"
	"market/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (n *MarketApi) GetMarketEventTypeForHome(c *gin.Context) {
	var res common.Response

	result, err := service.MarketService.GetMarketEventTypeForHome(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) GetMarketEventForHome(c *gin.Context) {
	var res common.Response
	var event request.GetMarketEventForHome

	err := c.ShouldBind(&event)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := service.MarketService.GetMarketEventForHome(c, event)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}
