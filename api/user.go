package api

import (
	"market/global"
	"market/model/common"
	"market/model/market/request"
	"market/service"
	"net/http"

	"github.com/gin-gonic/gin"
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
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithMessage("execution succeed")
	c.JSON(http.StatusOK, res)
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
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithMessage("execution succeed")
	c.JSON(http.StatusOK, res)
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
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithMessage("execution succeed")
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) LoginByCode(c *gin.Context) {
	var res common.Response
	var user request.UserLogin

	err := c.ShouldBind(&user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := service.MarketService.UserLoginByCode(user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) LoginByWallet(c *gin.Context) {
	var res common.Response
	var user request.UserLoginByWallet

	err := c.ShouldBind(&user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := service.MarketService.UserLoginByWallet(user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) GetUserInfo(c *gin.Context) {
	var res common.Response

	result, err := service.MarketService.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) GetUserProfile(c *gin.Context) {
	var res common.Response
	var profile request.GetUserProfile

	err := c.ShouldBind(&profile)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}
	result, err := service.MarketService.GetUserProfile(c, profile)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) GetUserSetting(c *gin.Context) {
	var res common.Response

	result, err := service.MarketService.GetUserSetting(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) UpdateUserSetting(c *gin.Context) {
	var res common.Response
	var user request.UpdateUserSetting

	err := c.ShouldBind(&user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	err = service.MarketService.UpdateUserSetting(c, user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithMessage("execution succeed")
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) UpdateUserNotificationSetting(c *gin.Context) {
	var res common.Response
	var user request.UpdateUserNotificationSetting

	err := c.ShouldBind(&user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	err = service.MarketService.UpdateUserNotificationSetting(c, user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithMessage("execution succeed")
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) CreateUserAffiliate(c *gin.Context) {
	var res common.Response
	var user request.CreateUserAffiliate

	err := c.ShouldBind(&user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := service.MarketService.CreateUserAffiliate(c, user)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) GetUserNotification(c *gin.Context) {
	var res common.Response

	result, err := service.MarketService.GetUserNotification(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) GetUserBalance(c *gin.Context) {
	var res common.Response

	result, err := service.MarketService.GetUserBalance(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}
