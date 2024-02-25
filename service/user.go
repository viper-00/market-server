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

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (m *MService) IsUserExist(email string) (bool, error) {
	var user model.User
	err := global.MARKET_DB.Where("email = ? AND status = 1", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if user.ID > 0 {
		return true, err
	}

	return false, err
}

func (m *MService) UserRegister(info request.UserRegister) error {
	if info.Email != "" {

		isExist, errExist := m.IsUserExist(info.Email)
		if isExist {
			return errors.New("user already exists")
		} else if errExist == nil && !isExist {
			var randomString = utils.GenerateStringRandomly("", 6)

			randomJWT, err := jwt.CreateJWT(map[string]interface{}{
				"email":           info.Email,
				"invitation_code": randomString,
				"chain_id":        info.ChainId,
				"time":            time.Now().UTC().Unix(),
			})

			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return err
			}

			confirmUrl := fmt.Sprintf("%s/confirm?code=%s", global.MARKET_CONFIG.Client.Url, randomJWT)

			if _, err = global.MARKET_REDIS.Set(context.Background(), randomJWT, randomString, time.Minute*20).Result(); err != nil {
				global.MARKET_LOG.Error(err.Error())
				return err
			}

			if err = mail.SendMail(info.Email, mail.UserLoginTemplate(info.Email, confirmUrl)); err != nil {
				global.MARKET_LOG.Error(err.Error())
				return err
			}

			return nil
		} else {
			return errors.New("service error")
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
	if info.Code != "" {
		claims, err := jwt.ValidateJWT(info.Code)
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return err
		}

		chainId := claims["chain_id"].(float64)
		email := claims["email"].(string)
		invitation_code := claims["invitation_code"].(string)
		// time := claims["time"]

		invitation_code_for_redis, err := global.MARKET_REDIS.Get(context.Background(), info.Code).Result()
		if err != nil {
			return errors.New("not found the code")
		}

		if invitation_code == invitation_code_for_redis {
			// initialze account
			err = m.InitializeAccount(int(chainId), email)
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return err
			}

			// delete code for redis
			_, err = global.MARKET_REDIS.Del(context.Background(), info.Code).Result()
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return err
			}

			return nil
		} else {
			return errors.New("code is invalid")
		}

	}

	return fmt.Errorf("not found the code: %s", info.Code)
}

func (n *MService) UserLogin(info request.UserLogin) (user response.User, err error) {
	if info.Email != "" {
		var modelUser model.User
		err = global.MARKET_DB.Where("email = ? AND status = 1", info.Email).First(&modelUser).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}

		var modelUserSetting model.UserSetting
		err = global.MARKET_DB.Where("user_id = ? AND status = 1", modelUser.ID).First(&modelUserSetting).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}

		jwtString, jwtErr := jwt.CreateJWT(map[string]interface{}{
			"chain_id":        info.ChainId,
			"email":           modelUser.Email,
			"address":         modelUser.Address,
			"contractAddress": modelUser.ContractAddress,
			"time":            time.Now().UTC().Unix(),
		})

		if jwtErr != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}

		user.ChainId = modelUser.ChainId
		user.Address = modelUser.Address
		user.ContractAddress = modelUser.ContractAddress
		user.Email = modelUser.Email
		user.Auth = fmt.Sprintf("Bearer %s", jwtString)
		user.InviteCode = modelUser.InvitationCode
		user.UserName = modelUserSetting.Username
		user.AvatarUrl = modelUserSetting.AvatarUrl
		user.Bio = modelUserSetting.Bio
		user.JoinedDate = modelUser.CreatedAt.Unix()

		if _, err = global.MARKET_REDIS.Set(context.Background(), user.Auth, info.Email, time.Minute*20).Result(); err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}

		return user, nil

	} else if info.Address != "" && info.ChainId != 0 {
		// todo:
	}

	return user, errors.New("no support")
}

func (m *MService) InitializeAccount(chainId int, email string) (err error) {
	privateKey, address, err := wallet.GenerateEthereumWallet()
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	contractAddress, err := wallet.GenerateEthereumCollectionContract(chainId, address)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	// default wallet for user
	var user model.User
	user.PrivateKey = privateKey
	user.Address = address
	user.Email = email
	user.ContractAddress = contractAddress
	user.Status = 1
	user.ChainId = chainId
	user.SuperiorId = 0
	user.InvitationCode = utils.GenerateStringRandomly("market_invite_code_", 8)
	err = global.MARKET_DB.Save(&user).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	// user info of setting
	var userSetting model.UserSetting
	userSetting.UserId = user.ID
	userSetting.Email = user.Email
	userSetting.Username = utils.GenerateStringRandomly("", 8)
	userSetting.AvatarUrl = ""
	userSetting.Bio = ""
	userSetting.Status = 1
	err = global.MARKET_DB.Save(&userSetting).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	// user notification of setting
	var userNotificationSetting model.UserNotificationSetting
	userNotificationSetting.UserId = user.ID
	userNotificationSetting.Email = user.Email
	userNotificationSetting.MarketUpdate = false
	err = global.MARKET_DB.Save(&userNotificationSetting).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return nil
}

// value, _ := c.Get("chainId")
// c.Set("chainId", chainId)
// c.Set("email", email)
// c.Set("address", address)
// c.Set("contractAddress", contractAddress)
// c.Set("time", time)

func (m *MService) GetUserInfo(c *gin.Context) (result interface{}, err error) {
	return
}

func (m *MService) UpdateUserInfo(c *gin.Context, req request.UpdateUserInfo) (result interface{}, err error) {
	return
}

func (m *MService) UpdateUserSetting(c *gin.Context, req request.UpdateUserSetting) (result interface{}, err error) {
	return
}

func (m *MService) UpdateUserNotificationSetting(c *gin.Context, req request.UpdateUserNotificationSetting) (result interface{}, err error) {
	return
}

func (m *MService) CreateUserAffiliate(c *gin.Context, req request.CreateUserAffiliate) (result interface{}, err error) {
	return
}
