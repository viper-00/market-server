package initialize

import (
	"context"
	"market/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Redis() {
	r := global.MARKET_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})

	if pong, err := client.Ping(context.Background()).Result(); err != nil {
		global.MARKET_LOG.Error("redis connect failed: ", zap.Error(err))
	} else {
		global.MARKET_LOG.Info("redis connect success: ", zap.String("pong", pong))
		global.MARKET_REDIS = client
	}
}
