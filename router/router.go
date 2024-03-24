package router

import (
	"market/api"

	"market/middleware"

	"github.com/gin-gonic/gin"
)

type MarketRouter struct{}

func (mRouter *MarketRouter) InitRouter(Router *gin.RouterGroup) {
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
		clientRouter.GET("verify-invitation", api.VerifyInvitation)
		clientRouter.POST("login", api.Login)
		clientRouter.POST("login-by-code", api.LoginByCode)
		clientRouter.GET("crypto-price", api.GetCryptoPrice)
		clientRouter.GET("user-profile", api.GetUserProfile)
		clientRouter.GET("market-event-type", api.GetMarketEventTypeForHome)
		clientRouter.GET("market-event", api.GetMarketEventForHome)
	}

	userRouter := clientRouter.Group("/user")
	userRouter.Use(middleware.ClientAuth())
	{
		userRouter.GET("user-setting", api.GetUserSetting)
		userRouter.PUT("user-setting", api.UpdateUserSetting)
		userRouter.PUT("user-notification-setting", api.UpdateUserNotificationSetting)
		userRouter.POST("user-affiliate", api.CreateUserAffiliate)
		userRouter.GET("user-notification", api.GetUserNotification)
		userRouter.GET("user-balance", api.GetUserBalance)
	}

	eventRouter := clientRouter.Group("/event")
	eventRouter.Use(middleware.ClientAuth())
	{
		eventRouter.GET("market-event", api.GetMarketEvent)
		eventRouter.POST("market-event", api.CreateMarketEvent)
		eventRouter.PUT("market-event", api.UpdateMarketEvent)
		eventRouter.POST("market-event-play", api.CreateMarketEventPlay)
		eventRouter.PUT("market-event-play", api.UpdateMarketEventPlay)
		eventRouter.GET("market-event-play", api.GetMarketEventPlay)
		eventRouter.GET("market-event-type", api.GetMarketEventType)
	}

	orderRouter := eventRouter.Group("/order")
	orderRouter.Use(middleware.ClientAuth())
	{
		orderRouter.POST("market-event-order", api.CreateMarketEventOrder)
		orderRouter.POST("market-event-order-settle", api.SettleMarketEventOrder)
	}

	commentRouter := eventRouter.Group("/comment")
	commentRouter.Use(middleware.ClientAuth())
	{
		commentRouter.POST("market-event-comment", api.CreateEventComment)
		commentRouter.GET("market-event-comment", api.GetEventComment)
		commentRouter.DELETE("market-event-comment", api.RemoveEventComment)
	}

	likeRouter := commentRouter.Group("/like")
	likeRouter.Use(middleware.ClientAuth())
	{
		likeRouter.POST("market-event-comment-like", api.CreateCommentLike)
		likeRouter.PUT("market-event-comment-like", api.UpdateCommentLike)
	}

	uploadRouter := clientRouter.Group("upload")
	uploadRouter.Use(middleware.ClientAuth())
	{
		uploadRouter.POST("uploadFile", api.UploadFile)
	}
}
