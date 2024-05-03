package service

import (
	"errors"
	"market/global"
	"market/global/constant"
	"market/model/market/request"
	"market/utils/wallet"

	"github.com/gin-gonic/gin"
)

func (m *MService) GetFreeCoin(c *gin.Context, req request.GetFreeCoin) (err error) {
	user, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	if !constant.IsSupportFreeChain(user.ChainId) {
		return errors.New("do not support the chain for free")
	}

	if !constant.IsSupportFreeCoin(req.Coin) {
		return errors.New("do not support the coin for free")
	}

	var hash string

	switch req.Coin {
	case constant.ETH:
		hash, err = wallet.TransferEthToReceiveAddress(user.ChainId, user.ContractAddress, 0.0017)
	case constant.USDT, constant.USDC:
		hash, err = wallet.TransferTokenToReceiveAddress(user.ChainId, user.ContractAddress, req.Coin, 10)
	default:
		return nil
	}

	if err != nil {
		return
	}

	if hash == "" {
		return errors.New("the transaction fails, please try it out")
	}

	return nil
}
