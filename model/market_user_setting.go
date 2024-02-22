package model

import "market/model/common"

type UserSetting struct {
	common.MARKET_MODEL
	UserId    uint   `json:"user_id" gorm:"user_id"`
	Email     string `json:"email" gorm:"comment:email"`
	Username  string `json:"username" gorm:"comment:username"`
	AvatarUrl string `json:"avatarUrl" gorm:"comment:avatarUrl"`
	Bio       string `json:"bio" gorm:"comment:bio"`
}

func (UserSetting) TableName() string {
	return "market_user_settings"
}
