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
	}

	userRouter := clientRouter.Group("/user")
	{
		userRouter.GET("user-info", api.GetUserInfo)
		userRouter.PUT("user-info", api.UpdateUserInfo)
		userRouter.PUT("user-setting", api.UpdateUserSetting)
		userRouter.PUT("user-notification-setting", api.UpdateUserNotificationSetting)
		userRouter.POST("user-affiliate", api.CreateUserAffiliate)
	}

	eventRouter := clientRouter.Group("/event")
	{
		eventRouter.POST("market-event", api.CreateMarketEvent)
		eventRouter.PUT("market-event", api.UpdateMarketEvent)
		eventRouter.POST("market-event-play", api.CreateMarketEventPlay)
		eventRouter.PUT("market-event-play", api.UpdateMarketEventPlay)
	}

	commentRouter := eventRouter.Group("/comment")
	{
		commentRouter.POST("market-event-comment", api.CreateEventComment)
		commentRouter.GET("market-event-comment", api.FindEventComment)
		commentRouter.DELETE("market-event-comment", api.RemoveEventComment)
	}

	likeRouter := eventRouter.Group("/like")
	{
		likeRouter.POST("market-event-comment-like", api.CreateCommentLike)
		likeRouter.PUT("market-event-comment-like", api.UpdateCommentLike)
	}

	uploadRouter := eventRouter.Group("upload")
	{
		uploadRouter.POST("uploadFile", api.UploadFile)
	}
}
