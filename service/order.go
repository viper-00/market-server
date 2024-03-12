package service

import (
	"errors"
	"market/global"
	"market/model"
	"market/model/market/request"
	"market/utils/wallet"

	"github.com/gin-gonic/gin"
)

func (m *MService) CreateMarketEventOrder(c *gin.Context, req request.CreateMarketEventOrder) (err error) {
	userModel, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	chainId, _ := c.Get("chainId")
	intChainId := int(chainId.(float64))

	eventModel, codeErr := m.GetMarketEventByUniqueCode(req.EventUniqueCode)
	if codeErr != nil {
		global.MARKET_LOG.Error(codeErr.Error())
		return codeErr
	}

	if req.Type == 1 {

		if req.Amount == 0 {
			return errors.New("amount must be greater than 0")
		}
		// send usdt to receive account
		hash, walletErr := wallet.TransferAssetToReceiveAddress(intChainId, userModel.ContractAddress, req.Amount)
		if walletErr != nil {
			global.MARKET_LOG.Error(walletErr.Error())
			return walletErr
		}

		var order model.EventOrder
		order.UserId = userModel.ID

		order.EventId = eventModel.ID
		order.Amount = req.Amount
		order.PlayValue = req.PlayValue
		order.OrderStatus = 1
		order.OrderType = req.Type
		order.Hash = hash

		err = global.MARKET_DB.Save(&order).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	} else if req.Type == 2 {
		// no send usdt, just record the sell order

		// get amount by fine the order
		var buyOrder model.EventOrder
		err = global.MARKET_DB.Where("user_id = ? AND event_id = ? AND play_value = ? AND order_status = 1 AND order_type = 1 AND hash != ?", userModel.ID, eventModel.ID, req.PlayValue, "").Order("id desc").First(&buyOrder).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}

		var order model.EventOrder
		order.UserId = userModel.ID
		order.EventId = eventModel.ID
		order.PlayValue = req.PlayValue
		order.OrderType = req.Type
		order.OrderStatus = 1
		order.Amount = buyOrder.Amount

		err = global.MARKET_DB.Save(&order).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	}

	return
}
