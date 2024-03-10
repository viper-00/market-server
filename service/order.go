package service

import (
	"market/global"
	"market/model"
	"market/model/market/request"

	"github.com/gin-gonic/gin"
)

func (m *MService) CreateMarketEventOrder(c *gin.Context, req request.CreateMarketEventOrder) (err error) {
	userModel, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	eventModel, err := m.GetMarketEventByUniqueCode(req.EventUniqueCode)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	var order model.EventOrder
	order.UserId = userModel.ID

	order.EventId = eventModel.ID
	order.Amount = req.Amount
	order.PlayValue = req.PlayValue
	order.OrderStatus = 1
	order.OrderType = req.Type

	err = global.MARKET_DB.Save(&order).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}
	return
}
