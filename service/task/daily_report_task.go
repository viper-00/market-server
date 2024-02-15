package task

import (
	"context"
	"errors"
	"fmt"
	"market/global"
	"market/global/constant"
	"market/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

func RunDailyReportTask() {
	for {
		now := time.Now()

		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		durationUntilMidnight := nextMidnight.Sub(now)

		ticker := time.NewTicker(durationUntilMidnight)

		<-ticker.C

		RunDailyReportCore()
	}
}

func RunDailyReportCore() {
	global.MARKET_LOG.Info("---------- Run Daily Report Task ----------")

	count, err := global.MARKET_REDIS.Get(context.Background(), constant.DAILY_REPORT_ERROR).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			count = "0"
		} else {
			RunDailyReportCore()
		}
	}

	message := fmt.Sprintf("[Daily Report] %s\n\nNumber of failures today: %s", time.Now().UTC().Format("2006-01-02 15:04:05"), count)

	if utils.InformToTelegram(message) {
		global.MARKET_REDIS.Del(context.Background(), constant.DAILY_REPORT_ERROR)
	}
}
