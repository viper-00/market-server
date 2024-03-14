package service

import (
	"errors"
	"market/global"
	"market/global/constant"
	"market/model"
	"market/model/market/request"
	sweepUtils "market/sweep/utils"
	"market/utils"
	"market/utils/wallet"
	"math/big"
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

func (m *MService) SettleMarketEventOrder(c *gin.Context, req request.SettleMarketOrder) (interface{}, error) {
	userModel, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	chainId, _ := c.Get("chainId")
	intChainId := int(chainId.(float64))

	eventModel, err := m.GetMarketEventByUniqueCode(req.EventUniqueCode)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	if eventModel.UserId != userModel.ID || eventModel.Password != utils.EncryptoThroughMd5([]byte(req.Password)) {
		return nil, errors.New("you don't have permission to perform this operation")
	}

	if eventModel.ExpireTime.After(time.Now()) {
		return nil, errors.New("the time has not come yet")
	}

	var eventPlay model.EventPlay
	err = global.MARKET_DB.Where("id = ? AND status = 1", eventModel.PlayId).First(&eventPlay).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	var eventOrders []model.EventOrder
	err = global.MARKET_DB.Where("event_id = ?", eventModel.ID).Order("id desc").Find(&eventOrders).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
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
					return nil, errors.New("some orders are not completed")
				}
			}
		}
	}

	if buyPlays != len(allPlays) {
		return nil, errors.New("some orders are not completed")
	}

	// Get random results
	randomResult := utils.GetRandomValueFromStringArray(allPlays)

	var winnerOrder model.EventOrder
	err = global.MARKET_DB.Where("order_type = 1 AND hash != ? AND order_status = 1 AND play_value = ?", "", randomResult).Order("id desc").First(&winnerOrder).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	// winnerWallet
	var winnerWallet model.User
	err = global.MARKET_DB.Where("id = ? AND status = 1", winnerOrder.UserId).First(&winnerWallet).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	// bankerWallet
	var bankerWallet model.User
	err = global.MARKET_DB.Where("id = ? AND status = 1", eventModel.UserId).First(&bankerWallet).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	// 资金池：押金+总买钱（A+O）+总卖钱
	// 获取总收益：资金池-买钱（A）-押金-卖钱（O）
	// A：买钱（A）
	// B：押金
	// D：总卖钱（O）
	// 收益率分配（百分比）：
	// A：买钱(A)/总买钱（O）
	// B：（1-买钱(A)/总买钱（O））/2
	// C：（1-买钱(A)/总买钱（O））/2

	var (
		depositAmount   float64 = 0
		totalBuyAmount  float64 = 0
		totalSellAmount float64 = 0
		victoryAmount   float64 = 0
		splitAmount     float64 = 0

		victoryIncomeRate float64 = 0
		bankerIncomeRate  float64 = 0
		// platformIncomeRate float64 = 0

		tokenAddresses  []string
		sendToAddresses []string
		sendValues      []big.Int
	)

	depositAmount = eventPlay.PledgeAmount

	isSupport, _, tokenContractAddress, decimals := sweepUtils.GetContractInfoByChainIdAndSymbol(intChainId, constant.USDT)
	if !isSupport {
		return nil, errors.New("contract address not found")
	}

	// 1. 发送资金 - 押金
	tokenAddresses = append(tokenAddresses, tokenContractAddress)
	sendToAddresses = append(sendToAddresses, bankerWallet.ContractAddress)
	sendValues = append(sendValues, *big.NewInt(utils.FormatToOriginalValue(depositAmount, decimals)))

	// 2. 发送资金 - W买钱
	tokenAddresses = append(tokenAddresses, tokenContractAddress)
	sendToAddresses = append(sendToAddresses, winnerWallet.ContractAddress)
	sendValues = append(sendValues, *big.NewInt(utils.FormatToOriginalValue(winnerOrder.Amount, decimals)))

	// 3 发送资金 - 所有卖钱
	for _, o := range eventOrders {
		if o.OrderType == 1 {
			totalBuyAmount += o.Amount
		} else if o.OrderType == 2 {
			totalSellAmount += o.Amount

			var sellUser model.User
			err = global.MARKET_DB.Where("id = ? AND status = 1", o.UserId).First(&sellUser).Error
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return nil, err
			}
			tokenAddresses = append(tokenAddresses, tokenContractAddress)
			sendToAddresses = append(sendToAddresses, sellUser.ContractAddress)
			sendValues = append(sendValues, *big.NewInt(utils.FormatToOriginalValue(o.Amount, decimals)))
		}
	}

	victoryAmount = winnerOrder.Amount
	splitAmount = totalBuyAmount - victoryAmount

	victoryIncomeRate = victoryAmount / totalBuyAmount
	bankerIncomeRate = (1 - (victoryIncomeRate)) / 2
	// platformIncomeRate = (1 - (victoryIncomeRate)) / 2

	// 发送资金 - 收益 - 买家
	tokenAddresses = append(tokenAddresses, tokenContractAddress)
	sendToAddresses = append(sendToAddresses, winnerWallet.ContractAddress)
	sendValues = append(sendValues, *big.NewInt(utils.FormatToOriginalValue(splitAmount*victoryIncomeRate, decimals)))

	// 发送资金 - 收益 - 庄家
	tokenAddresses = append(tokenAddresses, tokenContractAddress)
	sendToAddresses = append(sendToAddresses, bankerWallet.ContractAddress)
	sendValues = append(sendValues, *big.NewInt(utils.FormatToOriginalValue(splitAmount*bankerIncomeRate, decimals)))

	// 发送资金 - 收益 - 平台（不发）
	// tokenAddresses = append(tokenAddresses, tokenContractAddress)
	// sendToAddresses = append(sendToAddresses, "")
	// sendValues = append(sendValues, *big.NewInt(utils.FormatToOriginalValue(splitAmount * platformIncomeRate, decimals)))

	hash, err := wallet.TransferAssetToMoreReceiveAddres(intChainId, tokenAddresses, sendToAddresses, sendValues)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	return map[string]interface{}{
		"hash": hash,
	}, nil
}
