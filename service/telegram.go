package service

import (
	"errors"
	"market/global"
	"market/model"
	"market/model/market/request"
	"market/model/market/response"

	"gorm.io/gorm"
)

func (n *MService) GetUser(tgId string) (bool, error) {
	var findUser model.TelegramUser
	err := global.MARKET_DB.Where("tg_id = ? AND status = 1", tgId).First(&findUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	if findUser.ID > 0 {
		return true, nil
	}

	return false, nil
}

func (n *MService) CreateUser(user model.TelegramUser) (err error) {
	hasUser, err := n.GetUser(user.TGID)
	if err != nil {
		return
	}

	if hasUser {
		return nil
	}

	var tgUser model.TelegramUser
	tgUser.TGID = user.TGID
	tgUser.FirstName = user.FirstName
	tgUser.LastName = user.LastName
	tgUser.Username = user.Username
	tgUser.LanguageCode = user.LanguageCode
	tgUser.Status = 1

	return global.MARKET_DB.Create(&tgUser).Error
}

func (n *MService) SendMessageToTelegram(message request.SendMessageToTelegram) (errorKeys response.IntegrateTelegramKeyResponse, err error) {
	// err = global.MARKET_DB.Where("auth_key = ? AND status = 1", message.AuthKey).Find(&keys).Error
	// if err != nil {
	// 	global.MARKET_LOG.Error(err.Error())
	// 	return
	// }

	// if len(keys) == 0 {
	// 	var errorKey response.TelegramKey
	// 	errorKey.AuthKey = message.AuthKey
	// 	errorKey.Message = message.Message
	// 	errorKeys.IntegrateKey = append(errorKeys.IntegrateKey, errorKey)
	// } else {
	// 	for _, v := range keys {
	// 		if !utils.NotificationToTelegram(v.BotToken, v.TGID, message.Message) {
	// 			var errorKey response.TelegramKey
	// 			errorKey.TelegramUserId = v.TGID
	// 			errorKey.AuthKey = v.AuthKey
	// 			errorKey.Message = message.Message
	// 			errorKeys.IntegrateKey = append(errorKeys.IntegrateKey, errorKey)
	// 			continue
	// 		}
	// 	}
	// }

	// if len(errorKeys.IntegrateKey) > 0 {
	// 	return errorKeys, errors.New("some keys failed to be sent to telegram")
	// }

	return
}
