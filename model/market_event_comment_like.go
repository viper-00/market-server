package model

import (
	"market/model/common"
)

type EventCommentLike struct {
	common.MARKET_MODEL
	CommentId uint `json:"comment_id" gorm:"comment:comment_id"`
	IsLike    bool `json:"is_like" gorm:"comment:is_like"`
	UserId    uint `json:"user_id" gorm:"comment:user_id"`
}

func (EventCommentLike) TableName() string {
	return "market_event_comment_likes"
}
