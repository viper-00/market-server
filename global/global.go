package global

import (
	"market/config"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	MARKET_DB     *gorm.DB
	MARKET_DBList map[string]*gorm.DB
	MARKET_REDIS  *redis.Client
	MARKET_CONFIG config.Server
	MARKET_VP     *viper.Viper
	MARKET_LOG    *zap.Logger
	// MARKETs_Timer time.Timer = timer

	MARKET_MUTEX sync.Mutex
)
