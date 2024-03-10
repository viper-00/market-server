package service

import (
	"market/global"
	"market/global/constant"
	"market/model"
	"market/model/market/request"
	"market/model/market/response"
	"market/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func (m *MService) GetMarketEventByUniqueCode(code string) (model model.Event, err error) {
	err = global.MARKET_DB.Where("unique_website_code = ? AND event_status = 1", code).First(&model).Error
	return
}

func (m *MService) GetMarketEvent(c *gin.Context, req request.GetMarketEvent) (result interface{}, err error) {
	var eventModel model.Event
	err = global.MARKET_DB.Where("unique_website_code = ? AND event_status = 1 AND status = 1", req.Code).First(&eventModel).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	var eventPlay model.EventPlay
	err = global.MARKET_DB.Where("id = ? AND status = 1", eventModel.PlayId).First(&eventPlay).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	var eventPlayResponse response.EventPlayResponse
	eventPlayResponse.Title = eventPlay.Title
	eventPlayResponse.Introduce = eventPlay.Introduce
	eventPlayResponse.GuessNumber = eventPlay.GuessNumber
	eventPlayResponse.MinimumCapitalPool = eventPlay.MinimumCapitalPool
	eventPlayResponse.MaximumCapitalPool = eventPlay.MaximumCapitalPool
	eventPlayResponse.Coin = eventPlay.Coin
	eventPlayResponse.PledgeAmount = eventPlay.PledgeAmount

	var orderModel []model.EventOrder
	err = global.MARKET_DB.Where("event_id = ?", eventModel.ID).Find(&orderModel).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	var values []response.EventPlayValueResponse

	for _, v := range constant.AllPlays[eventPlay.Title] {
		var value response.EventPlayValueResponse
		value.Value = v
		if len(orderModel) > 0 {
			for _, o := range orderModel {
				if v == o.PlayValue {
					var singleOrder response.EventOrderResponse
					singleOrder.Amount = o.Amount
					singleOrder.OrderType = constant.AllOrderTypes[o.OrderType]

					var user model.User
					err = global.MARKET_DB.Where("id = ? AND status = 1", o.UserId).First(&user).Error
					if err != nil {
						global.MARKET_LOG.Error(err.Error())
						return
					}
					singleOrder.UserContractAddress = user.ContractAddress

					var userSetting model.UserSetting
					err = global.MARKET_DB.Where("user_id = ? AND status = 1", o.UserId).First(&userSetting).Error
					if err != nil {
						global.MARKET_LOG.Error(err.Error())
						return
					}
					singleOrder.Username = userSetting.Username

					value.Orders = append(value.Orders, singleOrder)
				}
			}
		}

		values = append(values, value)
	}

	eventPlayResponse.Values = values

	var eventComments []model.EventComment
	err = global.MARKET_DB.Where("event_id = ? AND status = 1", eventModel.ID).Find(&eventComments).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return map[string]interface{}{
		"event":   eventModel,
		"play":    eventPlayResponse,
		"comment": eventComments,
	}, nil
}

func (m *MService) CreateMarketEvent(c *gin.Context, req request.CreateMarketEvent) (result interface{}, err error) {
	userModel, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	var event model.Event
	event.UserId = userModel.ID
	event.Title = req.Title
	event.UniqueWebsiteCode = utils.GenerateStringRandomly("event_", 12)
	event.ExpireTime = time.Unix(req.ExpireTime/1000, (req.ExpireTime%1000)*int64(time.Millisecond))
	event.Type = req.Type

	var play model.EventPlay
	err = global.MARKET_DB.Where("title = ? AND status = 1", req.PlayType).First(&play).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	event.PlayId = play.ID
	event.EventLogo = req.EventLogo
	event.SettlementAddress = req.SettlementAddress
	event.ResolverAddress = req.ResolverAddress
	event.EventStatus = 1
	event.Password = utils.EncryptoThroughMd5([]byte(req.Password))
	event.Status = 1

	err = global.MARKET_DB.Save(&event).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return map[string]interface{}{
		"unique_code": event.UniqueWebsiteCode,
	}, nil
}

func (m *MService) UpdateMarketEvent(c *gin.Context, req request.UpdateMarketEvent) (result interface{}, err error) {
	return
}

func (m *MService) CreateMarketEventPlay(c *gin.Context, req request.CreateMarketEventPlay) (result interface{}, err error) {
	return
}

func (m *MService) UpdateMarketEventPlay(c *gin.Context, req request.UpdateMarketEventPlay) (result interface{}, err error) {
	return
}

func (m *MService) GetMarketEventPlay(c *gin.Context) (result interface{}, err error) {
	var model []model.EventPlay
	err = global.MARKET_DB.Where("status = 1").Find(&model).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return model, nil
}

func (m *MService) GetMarketEventType(c *gin.Context) (result interface{}, err error) {
	return []string{string(constant.EVENT_CRYPTO), string(constant.EVENT_BUSINESS), string(constant.EVENT_SCIENCE), string(constant.EVENT_GAME)}, nil
}
