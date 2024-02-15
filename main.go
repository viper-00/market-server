package main

import (
	"market/core"
	"market/global"
	"market/initialize"

	"go.uber.org/zap"
)

func main() {
	global.MARKET_VP = core.Viper()
	global.MARKET_LOG = core.Zap()
	zap.ReplaceGlobals(global.MARKET_LOG)
	global.MARKET_DB = initialize.Gorm()

	if global.MARKET_DB != nil {
		initialize.RegisterTables()
		db, _ := global.MARKET_DB.DB()
		defer db.Close()
	}

	core.RunWindowsServer()
}
