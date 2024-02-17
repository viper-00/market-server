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
	"market/utils"
	"market/utils/jwt"
	"market/utils/mail"
	"market/utils/wallet"
	"time"

	"gorm.io/gorm"
)

func (m *MService) IsUserExist(email string) (error, bool) {
	var user model.User
	err := global.MARKET_DB.Where("email = ? AND status = 1", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false
	}

	if user.ID > 0 {
		return nil, true
	}

	return err, false
}

const (
	ConstantAccountString = "account-%s"
)

func (m *MService) UserRegister(info request.UserRegister) (err error) {
	if info.Email != "" {

		errExist, isExist := m.IsUserExist(info.Email)
		if isExist {
			return errors.New("User already exists")
		} else if errExist == nil && !isExist {
			var randomString = utils.GenerateStringRandomly("", 6)

			randomkey := fmt.Sprintf(ConstantAccountString, info.Email)

			confirmUrl := fmt.Sprintf("%s/confirm?email=%s&invitation_code=%s", global.MARKET_CONFIG.Client.Url, info.Email, randomString)

			if _, err = global.MARKET_REDIS.Set(context.Background(), randomkey, randomString, time.Minute*20).Result(); err != nil {
				global.MARKET_LOG.Error(err.Error())
				return
			}

			if err = mail.SendMail(info.Email, mail.UserLoginTemplate(info.Email, confirmUrl)); err != nil {
				global.MARKET_LOG.Error(err.Error())
				return
			}

			return nil
		} else {
			return errors.New("Service error")
		}

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

func (m *MService) UserVerifyInvitation(info request.UserVerifyInvitation) (err error) {
	if info.Email != "" && info.InvitationCode != "" {

		randomkey := fmt.Sprintf(ConstantAccountString, info.Email)

		code, err := global.MARKET_REDIS.Get(context.Background(), randomkey).Result()
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return err
		}

		if code == info.InvitationCode {
			// initialze account
			err = m.InitializeAccount(info.Email)
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return err
			}

			// delete code for redis
			_, err = global.MARKET_REDIS.Del(context.Background(), randomkey).Result()
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return err
			}

			return nil
		} else {
			return errors.New("Account is invalid")
		}

	}
	return errors.New(fmt.Sprintf("not found email or code: %s, %s", info.Email, info.InvitationCode))
}

func (n *MService) UserLogin(info request.UserLogin) (user response.User, err error) {
	if info.Email != "" {
		err = global.MARKET_DB.Where("email = ? AND status = 1", info.Email).First(&user).Error
		if err != nil {
			return
		}

		if err != nil {
			return
		}

		user.Auth, err = jwt.CreateJWT(map[string]interface{}{
			"email":           user.Email,
			"address":         user.Address,
			"contractAddress": user.ContractAddress,
			"time":            time.Now().UTC().Unix(),
		})

		if err != nil {
			return
		}

		return user, nil

	} else if info.Address != "" && info.ChainId != 0 {
		// todo:
	}

	return user, errors.New("No support")
}

func (m *MService) InitializeAccount(email string) (err error) {
	err, privateKey, address := wallet.GenerateEthereumWallet()
	if err != nil {
		return
	}

	var user model.User
	user.PrivateKey = privateKey
	user.Address = address
	user.Email = email
	user.Status = 1

	err = global.MARKET_DB.Save(&user).Error
	if err != nil {
		return
	}

	return nil
}
