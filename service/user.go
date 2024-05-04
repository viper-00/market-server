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
				"time":            time.Now().UTC().UnixMilli(),
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

			if err = mail.SendMail(info.Email, mail.UserRegisterTemplate(info.Email, confirmUrl)); err != nil {
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

func (n *MService) UserLoginByCode(info request.UserLogin) (user response.User, err error) {
	email, err := global.MARKET_REDIS.Get(context.Background(), info.Code).Result()
	if err != nil {
		return user, errors.New("not found the code")
	}

	if email != info.Email {
		return user, errors.New("permission denied")
	}

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
		"chain_id":        modelUser.ChainId,
		"email":           modelUser.Email,
		"address":         modelUser.Address,
		"contractAddress": modelUser.ContractAddress,
		"time":            time.Now().UTC().UnixMilli(),
	})

	if jwtErr != nil {
		global.MARKET_LOG.Error(jwtErr.Error())
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
	user.JoinedDate = modelUser.CreatedAt.UnixMilli()

	if _, err = global.MARKET_REDIS.Set(context.Background(), user.Auth, info.Email, time.Hour*24).Result(); err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return user, nil
}

func (n *MService) UserLogin(info request.UserLogin) (err error) {
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

		var randomString = utils.GenerateStringRandomly("", 6)

		if err = mail.SendMail(info.Email, mail.UserLoginTemplate(modelUserSetting.Username, info.Email, randomString)); err != nil {
			global.MARKET_LOG.Error(err.Error())
			return err
		}

		if _, err = global.MARKET_REDIS.Set(context.Background(), randomString, info.Email, time.Minute*10).Result(); err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}

		return nil

	} else if info.Address != "" && info.ChainId != 0 {
		// todo:
	}

	return errors.New("no support")
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
	userNotificationSetting.EmailUpdate = 2
	userNotificationSetting.DailyUpdate = 2
	userNotificationSetting.IncomingUpdate = 2
	userNotificationSetting.OutgoingUpdate = 2
	userNotificationSetting.EventUpdate = 2
	userNotificationSetting.OrderUpdate = 2
	userNotificationSetting.CryptoPriceUpdate = 2
	userNotificationSetting.Status = 1
	err = global.MARKET_DB.Save(&userNotificationSetting).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return setup.SavePublicKeyToRedis(context.Background(), chainId, contractAddress)
}

// value, _ := c.Get("chainId")
// c.Set("chainId", chainId)
// c.Set("email", email)
// c.Set("address", address)
// c.Set("contractAddress", contractAddress)
// c.Set("time", time)

func (m *MService) GetUserInfo(c *gin.Context) (model model.User, err error) {
	chainId, _ := c.Get("chainId")
	contractAddress, _ := c.Get("contractAddress")

	err = global.MARKET_DB.Where("contract_address = ? AND chain_id = ? AND status = 1", contractAddress, chainId).First(&model).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return
}

func (m *MService) GetUserProfile(c *gin.Context, profile request.GetUserProfile) (interface{}, error) {
	var user model.User
	err := global.MARKET_DB.Where("contract_address = ? AND status = 1", profile.Address).First(&user).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	var setting model.UserSetting
	err = global.MARKET_DB.Where("user_id = ? AND status = 1", user.ID).First(&setting).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	var events []model.Event
	err = global.MARKET_DB.Where("user_id = ? AND status = 1", user.ID).Order("id desc").Find(&events).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	var orders []model.EventOrder
	err = global.MARKET_DB.Where("user_id = ? AND status = 1", user.ID).Order("id desc").Find(&orders).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	var filterOrders []response.EventOrderResponse
	for _, v := range orders {
		var filterOrder response.EventOrderResponse
		filterOrder.Amount = v.Amount
		filterOrder.OrderType = constant.AllOrderTypes[v.OrderType]
		filterOrder.CreatedTime = int(v.CreatedAt.UnixMilli())
		filterOrder.Hash = v.Hash

		var eventModel model.Event
		err = global.MARKET_DB.Where("id = ? AND status = 1", v.EventId).First(&eventModel).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return nil, err
		}

		var eventPlay model.EventPlay
		err = global.MARKET_DB.Where("id = ? AND status = 1", eventModel.PlayId).First(&eventPlay).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return nil, err
		}

		filterOrder.Coin = eventPlay.Coin
		filterOrders = append(filterOrders, filterOrder)
	}

	var comments []model.EventComment
	err = global.MARKET_DB.Where("user_id = ? AND status = 1", user.ID).Order("id desc").Find(&comments).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	return map[string]interface{}{
		"profile": map[string]interface{}{
			"contract_address": user.ContractAddress,
			"email":            user.Email,
			"invitation_code":  user.InvitationCode,
			"username":         setting.Username,
			"avatar_url":       setting.AvatarUrl,
			"bio":              setting.Bio,
			"created_time":     user.CreatedAt.UnixMilli(),
		},
		"event":   events,
		"order":   filterOrders,
		"comment": comments,
	}, nil
}

func (m *MService) GetUserSetting(c *gin.Context) (interface{}, error) {
	user, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	var setting model.UserSetting
	err = global.MARKET_DB.Where("user_id = ? AND status = 1", user.ID).First(&setting).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	var notification model.UserNotificationSetting
	err = global.MARKET_DB.Where("user_id = ? AND status = 1", user.ID).First(&notification).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	return map[string]interface{}{
		"setting": map[string]interface{}{
			"contract_address":    user.ContractAddress,
			"email":               user.Email,
			"invitation_code":     user.InvitationCode,
			"username":            setting.Username,
			"avatar_url":          setting.AvatarUrl,
			"bio":                 setting.Bio,
			"created_time":        user.CreatedAt.UnixMilli(),
			"email_update":        notification.EmailUpdate,
			"daily_update":        notification.DailyUpdate,
			"incoming_update":     notification.IncomingUpdate,
			"outgoing_update":     notification.OutgoingUpdate,
			"event_update":        notification.EventUpdate,
			"order_update":        notification.OrderUpdate,
			"crypto_price_update": notification.CryptoPriceUpdate,
		},
	}, nil
}

func (m *MService) UpdateUserSetting(c *gin.Context, req request.UpdateUserSetting) (err error) {
	user, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	maps := make(map[string]interface{})
	if req.Username != "" {
		maps["username"] = req.Username
	}
	if req.AvatarUrl != "" {
		maps["avatar_url"] = req.AvatarUrl
	}
	if req.Bio != "" {
		maps["bio"] = req.Bio
	}

	err = global.MARKET_DB.Model(&model.UserSetting{}).Where("user_id = ?", user.ID).Updates(maps).Error

	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}
	return nil
}

func (m *MService) UpdateUserNotificationSetting(c *gin.Context, req request.UpdateUserNotificationSetting) (err error) {
	user, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	if req.Status != 1 && req.Status != 2 {
		return errors.New("no operation authority")
	}

	switch req.Type {
	case "email":
		err = global.MARKET_DB.Model(&model.UserNotificationSetting{}).Where("user_id = ?", user.ID).Updates(map[string]interface{}{
			"email_update": req.Status,
		}).Error
	case "daily":
		err = global.MARKET_DB.Model(&model.UserNotificationSetting{}).Where("user_id = ?", user.ID).Updates(map[string]interface{}{
			"daily_update": req.Status,
		}).Error
	case "incoming":
		err = global.MARKET_DB.Model(&model.UserNotificationSetting{}).Where("user_id = ?", user.ID).Updates(map[string]interface{}{
			"incoming_update": req.Status,
		}).Error
	case "outgoing":
		err = global.MARKET_DB.Model(&model.UserNotificationSetting{}).Where("user_id = ?", user.ID).Updates(map[string]interface{}{
			"outgoing_update": req.Status,
		}).Error
	case "event":
		err = global.MARKET_DB.Model(&model.UserNotificationSetting{}).Where("user_id = ?", user.ID).Updates(map[string]interface{}{
			"event_update": req.Status,
		}).Error
	case "order":
		err = global.MARKET_DB.Model(&model.UserNotificationSetting{}).Where("user_id = ?", user.ID).Updates(map[string]interface{}{
			"order_update": req.Status,
		}).Error
	case "crypto":
		err = global.MARKET_DB.Model(&model.UserNotificationSetting{}).Where("user_id = ?", user.ID).Updates(map[string]interface{}{
			"crypto_price_update": req.Status,
		}).Error
	default:
		return
	}

	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}
	return nil

	// err = global.MARKET_DB.Model(&model.UserNotificationSetting{}).Where("user_id = ?", user.ID).Updates(map[string]interface{}{
	// 	"email_update":        req.EmailUpdate,
	// 	"daily_update":        req.DailyUpdate,
	// 	"incoming_update":     req.IncomingUpdate,
	// 	"outgoing_update":     req.OutgoingUpdate,
	// 	"event_update":        req.EventUpdate,
	// 	"order_update":        req.OrderUpdate,
	// 	"crypto_price_update": req.CryptoPriceUpdate,
	// }).Error
	// if err != nil {
	// 	global.MARKET_LOG.Error(err.Error())
	// 	return
	// }
	// return nil
}

func (m *MService) CreateUserAffiliate(c *gin.Context, req request.CreateUserAffiliate) (result interface{}, err error) {
	return
}

func (m *MService) GetUserNotification(c *gin.Context) (result interface{}, err error) {
	var nts []model.Notification
	chainId, _ := c.Get("chainId")
	contractAddress, _ := c.Get("contractAddress")

	var user model.User
	err = global.MARKET_DB.Where("chain_id = ? AND contract_address = ? AND status = 1", chainId, contractAddress).First(&user).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	err = global.MARKET_DB.Where("chain_id = ? AND user_id = ? AND is_read = ?", chainId, user.ID, constant.UNREAD).Order("id desc").Find(&nts).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return nts, nil
}

func (m *MService) GetUserBalance(c *gin.Context) (result interface{}, err error) {
	chainId, _ := c.Get("chainId")
	contractAddress, _ := c.Get("contractAddress")

	intChainId := int(chainId.(float64))

	result, err = wallet.GetAllTokenBalance(intChainId, contractAddress.(string))
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return result, nil
}
