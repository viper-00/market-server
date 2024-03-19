package mail

import (
	"errors"
	"fmt"
	"market/global"
	"net/smtp"
)

type emailAuth struct {
	username, password string
}

func EmailAuth(username, password string) smtp.Auth {
	return &emailAuth{username, password}
}

func (e *emailAuth) Start(server *smtp.ServerInfo) (proto string, toServer []byte, err error) {
	return "LOGIN", []byte{}, nil
}

func (e *emailAuth) Next(fromServer []byte, more bool) (toServer []byte, err error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(e.username), nil
		case "Password:":
			return []byte(e.password), nil
		default:
			return nil, errors.New("unkown fromServer")
		}
	}
	return nil, nil
}

func SendMail(to string, message []byte) error {
	return sendCoreMail(
		global.MARKET_CONFIG.Smtp.Host,
		global.MARKET_CONFIG.Smtp.Port,
		global.MARKET_CONFIG.Smtp.Username,
		global.MARKET_CONFIG.Smtp.Password,
		to,
		message,
	)
}

func sendCoreMail(smtpHost string, smtpPort int, smtpUsername string, smtpPassword string, toEmail string, message []byte) error {

	auth := EmailAuth(smtpUsername, smtpPassword)

	serverAddr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	to := []string{toEmail}
	err := smtp.SendMail(serverAddr, auth, smtpUsername, to, message)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
	}

	return nil
}
