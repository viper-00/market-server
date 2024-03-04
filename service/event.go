package service

import (
	"market/global"
	"market/global/constant"
	"market/model"
	"market/model/market/request"

	"github.com/gin-gonic/gin"
)

func (m *MService) CreateMarketEvent(c *gin.Context, req request.CreateMarketEvent) (result interface{}, err error) {
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
	var model model.EventPlay
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
