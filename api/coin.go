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

func (n *MarketApi) GetFreeCoin(c *gin.Context) {
	var res common.Response
	var coin request.GetFreeCoin

	err := c.ShouldBind(&coin)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	err = service.MarketService.GetFreeCoin(c, coin)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithMessage("get successfully")
	c.JSON(http.StatusOK, res)
}
