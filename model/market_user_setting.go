package model

import "market/model/common"

type UserSetting struct {
	common.MARKET_MODEL
	UserId    uint   `json:"-" gorm:"user_id"`
	Email     string `json:"email" gorm:"comment:email"`
	Username  string `json:"username" gorm:"comment:username"`
	AvatarUrl string `json:"avatar_url" gorm:"comment:avatar_url"`
	Bio       string `json:"bio" gorm:"comment:bio"`
}

func (UserSetting) TableName() string {
	return "market_user_settings"
}
