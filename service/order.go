package service

import (
	"errors"
	"market/global"
	"market/global/constant"
	"market/model"
	"market/model/market/request"
	"market/utils"
	"market/utils/wallet"
	"time"

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

func (m *MService) SettleMarketEventOrder(c *gin.Context, req request.SettleMarketOrder) (err error) {
	userModel, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	// chainId, _ := c.Get("chainId")
	// intChainId := int(chainId.(float64))

	eventModel, codeErr := m.GetMarketEventByUniqueCode(req.EventUniqueCode)
	if codeErr != nil {
		global.MARKET_LOG.Error(codeErr.Error())
		return codeErr
	}

	if eventModel.UserId != userModel.ID || eventModel.Password != utils.EncryptoThroughMd5([]byte(req.Password)) {
		err = errors.New("You don't have permission to perform this operation")
		return
	}

	if eventModel.ExpireTime.After(time.Now()) {
		return errors.New("The time has not come yet")
	}

	var eventPlay model.EventPlay
	err = global.MARKET_DB.Where("id = ? AND status = 1", eventModel.PlayId).First(&eventPlay).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	var eventOrders []model.EventOrder
	err = global.MARKET_DB.Where("event_id = ?", eventModel.ID).Order("id desc").Find(&eventOrders).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	allPlays := constant.AllPlays[eventPlay.Title]
	var (
		buyPlays int = 0
	)

	for _, v := range allPlays {
		for _, o := range eventOrders {
			if v == o.PlayValue {
				if o.OrderType == 1 {
					buyPlays += 1
					break
				} else if o.OrderType == 2 {
					return errors.New("Some orders are not completed")
				}
			}
		}
	}

	if buyPlays != len(allPlays) {
		return errors.New("Some orders are not completed")
	}

	var (
		// totalAmount     float64 = 0
		totalBuyAmount  float64 = 0
		totalSellAmount float64 = 0
	)

	for _, o := range eventOrders {
		if o.OrderType == 1 {
			totalBuyAmount += o.Amount
		} else if o.OrderType == 2 {
			totalSellAmount += o.Amount
		}
	}

	return
}
