package core

import (
	"fmt"
	"market/global"
	"market/initialize"
	"market/service"
	"market/service/task"
	"market/sweep"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.MARKET_CONFIG.System.UseInit {
		err := service.MarketService.InitChainList()
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	}

	if global.MARKET_CONFIG.System.UseRedis {
		initialize.Redis()
	}

	if global.MARKET_CONFIG.System.UseTask {
		task.RunTask()
	}

	if global.MARKET_CONFIG.Blockchain.OpenSweepBlock {
		sweep.RunBlockSweep()
	}

	router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.MARKET_CONFIG.System.Addr)
	server := initServer(address, router)

	global.MARKET_LOG.Info("server run success on", zap.String("address", address))

	global.MARKET_LOG.Error(server.ListenAndServe().Error())
}

// Hot reload for gin server
func initServer(address string, router *gin.Engine) Server {
	server := endless.NewServer(address, router)
	server.ReadHeaderTimeout = 20 * time.Second
	server.WriteTimeout = 20 * time.Second
	server.MaxHeaderBytes = 1 << 20

	return server
}
