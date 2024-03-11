package initialize

import (
	"market/global"
	"market/middleware"
	"market/router"
	"net/http"

	"github.com/gin-gonic/gin"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

func Routers() *gin.Engine {

	SetGinMode(global.MARKET_CONFIG.System.Env)

	newRouter := gin.New()

	newRouter.Use(middleware.Cors())

	// newRouter.GET(global.MARKET_CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	newRouter.StaticFS("images", http.Dir("resource/images"))
	newRouter.StaticFS("files", http.Dir("resource/files"))

	newRouter.MaxMultipartMemory = 1 << 20

	MarketRouter := new(router.MarketRouter)

	Group := newRouter.Group(global.MARKET_CONFIG.System.RouterPrefix)
	{
		Group.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})

		MarketRouter.InitRouter(Group)
	}

	global.MARKET_LOG.Info("router register success")
	return newRouter
}

func SetGinMode(mode string) {
	switch mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case
		"release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
