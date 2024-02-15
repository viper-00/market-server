package router

import (
	"market/api"

	"market/middleware"

	"github.com/gin-gonic/gin"
)

type MainRouter struct{}

func (mRouter *MainRouter) InitRouter(Router *gin.RouterGroup) {
	internalRouter := Router.Group("/internal")
	api := new(api.MarketApi)
	{
		internalRouter.GET("test", api.Test)
		internalRouter.GET("networkInfo", api.GetNetworkInfo)
		internalRouter.POST("storeWalletAddress", api.StoreWalletAddress)
		internalRouter.POST("bulkStoreUserWallet", api.BulkStoreUserWallet)
		internalRouter.POST("storeChainContract", api.StoreChainContract)
		internalRouter.POST("bulkStoreChainContract", api.BulkStoreChainContract)
		internalRouter.GET("getTransactionByChainAndHash", api.GetTransactionByChainAndHash)
		internalRouter.GET("getTransactionsByChainAndAddress", api.GetTransactionsByChainAndAddress)
		internalRouter.POST("sendMessageToTelegram", api.SendMessageToTelegram)
		internalRouter.POST("revokeTelegramKey", api.RevokeTelegramKey)
	}

	internalRouter.Use(middleware.Wss())
	{
		internalRouter.GET("ws", api.WsForTxInfo)
	}

	clientRouter := Router.Group("/client")
	{
		clientRouter.GET("test", api.Test)

		// login and register
		clientRouter.POST("register", api.Register)
		clientRouter.POST("verify-invitation", api.VerifyInvitation)
		clientRouter.POST("login", api.Login)
	}

}
