package api

import (
	"market/global"
	"market/model/common"
	"market/model/market/request"
	"market/model/market/response"
	"market/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (n *MarketApi) Test(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (n *MarketApi) GetNetworkInfo(c *gin.Context) {
	var res common.Response
	var info request.GetNetworkInfo

	err := c.ShouldBind(&info)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := service.MarketService.GetInfo(info)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) StoreWalletAddress(c *gin.Context) {
	var res common.Response
	var wallet request.StoreUserWallet

	err := c.ShouldBind(&wallet)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	err = service.MarketService.StoreUserWallet(wallet)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithMessage("store successfully")
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) BulkStoreUserWallet(c *gin.Context) {
	var res common.Response
	var wallets request.BulkStoreUserWallet

	err := c.ShouldBind(&wallets)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := service.MarketService.BulkStorageUserWallets(wallets)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithDetailed(common.Error, err.Error(), result)
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) StoreChainContract(c *gin.Context) {
	var res common.Response
	var contract request.StoreChainContract

	err := c.ShouldBind(&contract)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	err = service.MarketService.StoreChainContract(contract)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithMessage("store successfully")
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) BulkStoreChainContract(c *gin.Context) {
	var res common.Response
	var contracts request.BulkStoreChainContract

	err := c.ShouldBind(&contracts)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := service.MarketService.BulkStorageChainContract(contracts)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithDetailed(common.Error, err.Error(), result)
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) GetTransactionByChainAndHash(c *gin.Context) {
	var res common.Response
	var tx request.TransactionByChainAndHash

	err := c.ShouldBind(&tx)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := service.MarketService.GetTransactionByChainAndHash(tx.ChainId, tx.Hash)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithDetailed(common.Error, err.Error(), result)
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) GetTransactionsByChainAndAddress(c *gin.Context) {
	var res common.Response
	var tx request.TransactionsByChainAndAddress

	err := c.ShouldBind(&tx)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	result, total, err := service.MarketService.GetTransactionsByChainAndAddress(tx)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithDetailed(common.Error, err.Error(), result)
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithDetailed(common.Success, "Request data successful", response.OwnListResponse{
		Transactions: result,
		Total:        total,
		Page:         tx.Page,
		PageSize:     tx.PageSize,
	})

	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) SendMessageToTelegram(c *gin.Context) {
	var res common.Response
	var message request.SendMessageToTelegram

	err := c.ShouldBind(&message)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := service.MarketService.SendMessageToTelegram(message)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithDetailed(common.Error, err.Error(), result)
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(result)
	c.JSON(http.StatusOK, res)
}

func (n *MarketApi) RevokeTelegramKey(c *gin.Context) {
	var res common.Response
	var key request.RevokeTelegramKey

	err := c.ShouldBind(&key)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	err = service.MarketService.RevokeTelegramKey(key)
	if err != nil {
		global.MARKET_LOG.Error(err.Error(), zap.Error(err))
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OkWithMessage("revoke successfully")
	c.JSON(http.StatusOK, res)
}
