package service

import (
	"market/global/constant"
	"market/model/market/request"
	"market/model/market/response"

	"github.com/gin-gonic/gin"
)

func (m *MService) GetMarketEventTypeForHome(c *gin.Context) (result interface{}, err error) {
	return []string{string(constant.EVENT_ALL), string(constant.EVENT_FOR_YOU), string(constant.EVENT_CRYPTO), string(constant.EVENT_BUSINESS), string(constant.EVENT_SCIENCE), string(constant.EVENT_GAME)}, nil
}

func (m *MService) GetMarketEventForHome(c *gin.Context, req request.GetMarketEventForHome) (result []response.EventForHomeResponse, err error) {
	switch req.Type {
	case string(constant.EVENT_ALL):
		// global.MARKET_DB.Where("status = 1")
	}

	return
}
