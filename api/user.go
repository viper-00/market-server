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

func (n *MarketApi) Register(c *gin.Context) {
	var res common.Response
	var user request.UserRegister

	err := c.ShouldBind(&user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	err = service.MarketService.UserRegister(user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (n *MarketApi) VerifyInvitation(c *gin.Context) {
	var res common.Response
	var invitation request.UserVerifyInvitation

	err := c.ShouldBind(&invitation)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	err = service.MarketService.UserVerifyInvitation(invitation)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (n *MarketApi) Login(c *gin.Context) {
	var res common.Response
	var user request.UserLogin

	err := c.ShouldBind(&user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	err = service.MarketService.UserLogin(user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusOK, nil)
}
