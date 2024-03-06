package service

import (
	"market/global"
	"market/global/constant"
	"market/model"
	"market/model/market/request"
	"market/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func (m *MService) CreateMarketEvent(c *gin.Context, req request.CreateMarketEvent) (err error) {
	var event model.Event
	event.Title = req.Title
	event.UniqueWebsiteLink = utils.GenerateStringRandomly("event_", 12)
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

	return
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
