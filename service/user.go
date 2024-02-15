package service

import (
	"errors"
	"fmt"
	"market/global/constant"
	"market/model/market/request"
)

func (n *MService) UserRegister(info request.UserRegister) (err error) {
	if info.Email != "" {

	} else if info.Address != "" && info.ChainId != 0 {
		if !constant.IsNetworkSupport(info.ChainId) {
			return errors.New("do not support network")
		}

		address := constant.AddressToLower(info.ChainId, info.Address)

		if !constant.IsAddressSupport(info.ChainId, address) {
			return fmt.Errorf("do not support wallet address: id: %d, address: %s", info.ChainId, info.Address)
		}

		// generate wallet for contract address
	} else {
		return errors.New("not found")
	}

	return nil
}

func (n *MService) UserVerifyInvitation(info request.UserVerifyInvitation) (err error) {

	return nil
}

func (n *MService) UserLogin(info request.UserLogin) (err error) {

	return nil
}
