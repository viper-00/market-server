package api

import (
	"market/global"
	"market/model/common"
	"market/model/market/request"
	"market/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (n *MarketApi) CreateCommentLike(c *gin.Context) {
	var res common.Response
	var like request.CreateCommentLike

	err := c.ShouldBind(&like)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	err = service.MarketService.CreateCommentLike(c, like)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithMessage("execution succeed")
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) UpdateCommentLike(c *gin.Context) {
	var res common.Response
	var like request.UpdateCommentLike

	err := c.ShouldBind(&like)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	err = service.MarketService.UpdateCommentLike(c, like)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithMessage("execution succeed")
	c.JSON(http.StatusOK, res)
}
