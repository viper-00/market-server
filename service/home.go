package service

import (
	"errors"
	"market/global"
	"market/global/constant"
	"market/model"
	"market/model/market/request"
	"market/model/market/response"
	"market/utils/wallet"

	"github.com/gin-gonic/gin"
)

func (m *MService) GetMarketEventTypeForHome(c *gin.Context) (result interface{}, err error) {
	return []string{string(constant.EVENT_ALL), string(constant.EVENT_FOR_YOU), string(constant.EVENT_CRYPTO), string(constant.EVENT_BUSINESS), string(constant.EVENT_SCIENCE), string(constant.EVENT_GAME)}, nil
}

func (m *MService) GetMarketEventForHome(c *gin.Context, req request.GetMarketEventForHome) (result []response.EventForHomeResponse, err error) {
	var searchEvent []model.Event

	var limit = 10
	switch req.Type {
	case string(constant.EVENT_ALL):
		err = global.MARKET_DB.Where("event_status = 1 AND status = 1").Order("id desc").Limit(limit).Find(&searchEvent).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	case string(constant.EVENT_FOR_YOU):
		err = global.MARKET_DB.Where("event_status = 1 AND status = 1").Order("RAND()").Limit(limit).Find(&searchEvent).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	case string(constant.EVENT_CRYPTO):
		err = global.MARKET_DB.Where("event_status = 1 AND status = 1 AND type = ?", string(constant.EVENT_CRYPTO)).Order("id desc").Limit(limit).Find(&searchEvent).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	case string(constant.EVENT_BUSINESS):
		err = global.MARKET_DB.Where("event_status = 1 AND status = 1 AND type = ?", string(constant.EVENT_BUSINESS)).Order("id desc").Limit(limit).Find(&searchEvent).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	case string(constant.EVENT_SCIENCE):
		err = global.MARKET_DB.Where("event_status = 1 AND status = 1 AND type = ?", string(constant.EVENT_SCIENCE)).Order("id desc").Limit(limit).Find(&searchEvent).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	case string(constant.EVENT_GAME):
		err = global.MARKET_DB.Where("event_status = 1 AND status = 1 AND type = ?", string(constant.EVENT_GAME)).Order("id desc").Limit(limit).Find(&searchEvent).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	default:
		return nil, errors.New("no support")
	}

	if len(searchEvent) > 0 {
		for _, v := range searchEvent {
			var r response.EventForHomeResponse
			r.EventLogo = v.EventLogo
			r.Title = v.Title
			r.ExpireTime = int(v.ExpireTime.UnixMilli())
			r.SettlementTime = int(v.SettlementTime.UnixMilli())
			r.Type = v.Type
			r.UniqueCode = v.UniqueWebsiteCode

			var play model.EventPlay
			err = global.MARKET_DB.Where("id = ? AND status = 1", v.PlayId).First(&play).Error
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return
			}

			r.Coin = play.Coin

			var orders []model.EventOrder
			err = global.MARKET_DB.Where("event_id = ? AND status = 1 AND order_status = 1", v.ID).Find(&orders).Error
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return
			}

			for _, o := range orders {
				r.TotalOrderAmount += o.Amount
			}

			var comments []model.EventComment
			err = global.MARKET_DB.Where("event_id = ? AND status = 1", v.ID).Find(&comments).Error
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return
			}
			r.CommentCount = len(comments)

			result = append(result, r)
		}
	}

	return
}

func (m *MService) GetTopVolumnForHome(c *gin.Context) (results []response.TopVolumnForHomeResponse, err error) {
	// chainId, _ := c.Get("chainId")
	// intChainId := int(chainId.(float64))
	intChainId := constant.OP_SEPOLIA

	var limit = 10

	var users []model.User
	err = global.MARKET_DB.Where("status = 1").Order("id desc").Limit(limit).Find(&users).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	for _, v := range users {
		var result response.TopVolumnForHomeResponse
		result.UserContractAddress = v.ContractAddress

		balance, balanceErr := wallet.GetAllTokenBalance(intChainId, v.ContractAddress)
		if balanceErr != nil {
			global.MARKET_LOG.Error(balanceErr.Error())
			return
		}

		result.EthBalance = balance.ETH
		result.UsdtBalance = balance.USDT
		result.UsdcBalance = balance.USDC

		var setting model.UserSetting
		err = global.MARKET_DB.Where("user_id = ? AND status = 1", v.ID).First(&setting).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}

		result.AvatarUrl = setting.AvatarUrl
		result.Username = setting.Username

		results = append(results, result)
	}

	return
}

func (m *MService) GetRecentActivityForHome(c *gin.Context) (results []response.RecentActivityForHomeResponse, err error) {
	var limit = 5

	var orders []model.EventOrder
	err = global.MARKET_DB.Where("status = 1").Order("id desc").Limit(limit).Find(&orders).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	for _, v := range orders {
		var result response.RecentActivityForHomeResponse
		result.Amount = v.Amount
		result.OrderType = constant.AllOrderTypes[v.OrderType]
		result.CreatedTime = int(v.CreatedAt.UnixMilli())
		result.PlayValue = v.PlayValue

		var event model.Event
		err = global.MARKET_DB.Where("id = ? AND status = 1", v.EventId).First(&event).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}

		var play model.EventPlay
		err = global.MARKET_DB.Where("id = ? AND status = 1", event.PlayId).First(&play).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}

		result.Coin = play.Coin

		result.EventLogo = event.EventLogo
		result.Title = event.Title
		result.UniqueCode = event.UniqueWebsiteCode

		var setting model.UserSetting
		err = global.MARKET_DB.Where("user_id = ? AND status = 1", v.UserId).First(&setting).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
		result.Username = setting.Username
		result.AvatarUrl = setting.AvatarUrl

		results = append(results, result)
	}

	return
}
