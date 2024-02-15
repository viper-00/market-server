package middleware

import (
	"fmt"
	"market/model"
	"market/service"

	"gopkg.in/telebot.v3"
)

func HandleTelegramUser(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {

		if c.Sender().IsBot {
			return nil
		}

		if c.Sender().ID >= 0 {
			var user model.TelegramUser
			user.TGID = fmt.Sprint(c.Sender().ID)
			user.FirstName = c.Sender().FirstName
			user.LastName = c.Sender().LastName
			user.Username = c.Sender().Username
			user.LanguageCode = c.Sender().LanguageCode

			if err := service.MarketService.CreateUser(user); err != nil {
				return nil
			}
		}
		return next(c)
	}
}

func HandlePrivateBot(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if c.Chat().Type == telebot.ChatPrivate {
			return next(c)
		}

		return nil
	}
}
