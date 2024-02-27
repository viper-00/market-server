package initialize

import (
	"market/global"
	"market/model"
	"os"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.MARKET_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	default:
		return GormMysql()
	}
}

func RegisterTables() {
	db := global.MARKET_DB
	if err := db.AutoMigrate(
		model.Chain{},
		model.User{},
		model.UserSetting{},
		model.UserNotificationSetting{},
		model.UserAffiliate{},
		model.Event{},
		model.EventComment{},
		model.EventCommentLike{},
		model.EventPlay{},
		model.Transaction{},
		model.OwnTransaction{},
		model.Notification{},
	); err != nil {
		global.MARKET_LOG.Error("db: register table failed", zap.Error(err))
		os.Exit(0)
	}

	global.MARKET_LOG.Info("db: register table success")
}
