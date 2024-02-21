package model

import "market/model/common"

type Setting struct {
	common.MARKET_MODEL
	UserId    uint   `json:"user_id" gorm:"user_id"`
	Email     string `json:"email" gorm:"comment:email"`
	Username  string `json:"username" gorm:"comment:username"`
	AvatarUrl string `json:"avatarUrl" gorm:"comment:avatarUrl"`
	Bio       string `json:"bio" gorm:"comment:bio"`
}

func (Setting) TableName() string {
	return "market_user_settings"
}
