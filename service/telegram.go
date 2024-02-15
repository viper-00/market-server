package service

import (
	"errors"
	"fmt"
	"market/global"
	"market/model"
	"market/model/market/request"
	"market/model/market/response"
	"market/utils"

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

func (n *MService) GetAllInOneBindingStatus(tgId string) (message string, err error) {
	var findKeys []model.TelegramAllInOneKey
	err = global.MARKET_DB.Where("tg_id =? AND status = 1", tgId).Find(&findKeys).Error
	keyLength := len(findKeys)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
	}

	return fmt.Sprintf(`The number of AllInOne keys you have bound: %d`, keyLength), nil
}

func (n *MService) GetAllInOneKey(tgId, key string) (bool, error) {
	var findKey model.TelegramAllInOneKey
	err := global.MARKET_DB.Where("tg_id = ? AND auth_key = ? AND status = 1", tgId, key).First(&findKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	if findKey.ID > 0 {
		return true, nil
	}

	return false, nil
}

func (n *MService) CreateAllInOneKey(tgId, key string) (err error) {
	hasKey, err := n.GetAllInOneKey(tgId, key)
	if err != nil {
		return
	}

	if hasKey {
		return nil
	}

	var keyModel model.TelegramAllInOneKey
	keyModel.TGID = tgId
	keyModel.AuthKey = key
	keyModel.BotToken = global.MARKET_CONFIG.Telegram.AllInOneNotificationBotToken
	keyModel.Status = 1

	return global.MARKET_DB.Create(&keyModel).Error
}

func (n *MService) SendMessageToTelegram(message request.SendMessageToTelegram) (errorKeys response.IntegrateTelegramKeyResponse, err error) {
	var keys []model.TelegramAllInOneKey
	err = global.MARKET_DB.Where("auth_key = ? AND status = 1", message.AuthKey).Find(&keys).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	if len(keys) == 0 {
		var errorKey response.TelegramKey
		errorKey.AuthKey = message.AuthKey
		errorKey.Message = message.Message
		errorKeys.IntegrateKey = append(errorKeys.IntegrateKey, errorKey)
	} else {
		for _, v := range keys {
			if !utils.NotificationToTelegram(v.BotToken, v.TGID, message.Message) {
				var errorKey response.TelegramKey
				errorKey.TelegramUserId = v.TGID
				errorKey.AuthKey = v.AuthKey
				errorKey.Message = message.Message
				errorKeys.IntegrateKey = append(errorKeys.IntegrateKey, errorKey)
				continue
			}
		}
	}

	if len(errorKeys.IntegrateKey) > 0 {
		return errorKeys, errors.New("some keys failed to be sent to telegram")
	}

	return
}

func (n *MService) RevokeTelegramKey(message request.RevokeTelegramKey) (err error) {
	err = global.MARKET_DB.Model(&model.TelegramAllInOneKey{}).Where("auth_key = ? AND status = 1", message.AuthKey).Update("status", 2).Error
	return
}
