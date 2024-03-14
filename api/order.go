package api

import (
	"market/global"
	"market/model/common"
	"market/model/market/request"
	"market/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (n *MarketApi) CreateMarketEventOrder(c *gin.Context) {
	var res common.Response
	var order request.CreateMarketEventOrder

	err := c.ShouldBind(&order)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	err = service.MarketService.CreateMarketEventOrder(c, order)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithMessage("execution succeed")
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) SettleMarketEventOrder(c *gin.Context) {
	var res common.Response
	var order request.SettleMarketOrder

	err := c.ShouldBind(&order)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := service.MarketService.SettleMarketEventOrder(c, order)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithDetailed(common.Error, err.Error(), result)
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}
