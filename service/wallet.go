package service

import (
	"context"
	"errors"
	"fmt"
	"market/global"
	"market/global/constant"
	"market/model"
	"market/model/market/request"
	"market/model/market/response"
	"market/sweep/setup"

	"gorm.io/gorm"
)

func (n *MService) BulkStorageUserWallets(wallets request.BulkStoreUserWallet) (errWalletResponses response.BulkStoreUserWalletResponse, err error) {
	if len(wallets.BulkStorage) > 0 {

		for _, v := range wallets.BulkStorage {
			if err = n.saveWallet(v.ChainId, v.Address); err != nil {
				var errorWallet response.StoreUserWallet
				errorWallet.ChainId = v.ChainId
				errorWallet.Address = v.Address
				errWalletResponses.BulkStorage = append(errWalletResponses.BulkStorage, errorWallet)
				global.MARKET_LOG.Error(err.Error())

				continue
			}
		}

		if len(errWalletResponses.BulkStorage) > 0 {
			return errWalletResponses, errors.New("some wallets failed to store")
		}
	}

	return
}

func (n *MService) StoreUserWallet(wallet request.StoreUserWallet) (err error) {
	return n.saveWallet(wallet.ChainId, wallet.Address)
}

func (n *MService) HasWalletByChainIdAndAddress(bId int, address string) (hasWallet bool, err error) {
	var findWallet model.Wallet

	err = global.MARKET_DB.Where("chain_id = ? AND address = ?", bId, address).First(&findWallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	if findWallet.ID > 0 {
		return true, nil
	}

	return false, nil
}

func (n *MService) saveWallet(chainId int, address string) (err error) {
	if !constant.IsNetworkSupport(chainId) {
		return errors.New("do not support network")
	}

	address = constant.AddressToLower(chainId, address)

	if !constant.IsAddressSupport(chainId, address) {
		return fmt.Errorf("do not support wallet address: id: %d, address: %s", chainId, address)
	}

	hasWallet, err := n.HasWalletByChainIdAndAddress(chainId, address)
	if err != nil {
		return
	}

	if hasWallet {
		return nil
	}

	var saveWallet model.Wallet
	saveWallet.Address = address
	saveWallet.ChainId = chainId
	saveWallet.NetworkName = constant.GetChainName(chainId)
	saveWallet.Status = 1

	if err = global.MARKET_DB.Create(&saveWallet).Error; err != nil {
		return
	}

	return setup.SavePublicKeyToRedis(context.Background(), chainId, address)
}
