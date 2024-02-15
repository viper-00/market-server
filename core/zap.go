package core

import (
	"market/core/internal"
	"market/global"
	"market/utils"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Zap() *zap.Logger {
	if ok, _ := utils.PathExists(global.MARKET_CONFIG.Zap.Director); !ok {
		os.Mkdir(global.MARKET_CONFIG.Zap.Director, os.ModePerm)
	}

	cores := internal.Zap.GetZapCores()
	logger := zap.New(zapcore.NewTee(cores...))

	if global.MARKET_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	globalHook := func(entry zapcore.Entry) error {
		switch entry.Level {
		case zapcore.ErrorLevel:
			// utils.InformToTelegram(fmt.Sprintf("[%s]\n\n%s | %s\n\n %s", entry.Time.UTC().Format("2006-01-02 15:04:05"), entry.Level.CapitalString(), utils.GenerateStringRandomly("market_", 12), entry.Message))
			// count, err := global.MARKET_REDIS.Get(context.Background(), constant.DAILY_REPORT_ERROR).Result()
			// if err == nil || errors.Is(err, redis.Nil) {
			// 	var countInt int64
			// 	if count == "" {
			// 		countInt = 0
			// 	} else {
			// 		countInt, err = strconv.ParseInt(count, 10, 64)
			// 		if err != nil {
			// 			return nil
			// 		}
			// 	}
			// 	global.MARKET_REDIS.Set(context.Background(), constant.DAILY_REPORT_ERROR, countInt+1, 0)
			// }
		}
		return nil
	}

	logger = logger.WithOptions(zap.Hooks(globalHook))

	return logger
}
