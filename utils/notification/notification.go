package notification

import (
	"context"
	"fmt"
	"market/global"
	"market/global/constant"
	"market/model/market/request"
	"market/service"
	MARKET_Client "market/utils/http"
)

var (
	client MARKET_Client.Client
)

func NotificationRequest(request request.NotificationRequest) (err error) {

	err = handleNotification(request)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
	}

	return nil
}

func handleNotification(request request.NotificationRequest) (err error) {
	err = service.MarketService.SaveTx(request.Chain, request.Hash)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	ownId, err := service.MarketService.SaveOwnTx(request)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	if ownId == 0 {
		global.MARKET_LOG.Info(fmt.Sprintf("OwnId already existed, hash: %s", request.Hash))
		return
	}

	_, err = global.MARKET_REDIS.RPush(context.Background(), constant.WS_NOTIFICATION, ownId).Result()
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return nil
}
