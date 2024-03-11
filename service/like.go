package service

import (
	"market/global"
	"market/model"
	"market/model/market/request"

	"github.com/gin-gonic/gin"
)

func (m *MService) CreateCommentLike(c *gin.Context, req request.CreateCommentLike) (err error) {
	userModel, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	var model model.EventCommentLike
	model.UserId = userModel.ID
	model.CommentId = req.CommentId
	model.IsLike = true
	err = global.MARKET_DB.Save(&model).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return
}

func (m *MService) UpdateCommentLike(c *gin.Context, req request.UpdateCommentLike) (err error) {
	userModel, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	err = global.MARKET_DB.Model(&model.EventCommentLike{}).Where("user_id = ? AND comment_id = ?", userModel.ID, req.CommentId).Update("is_like", req.IsLike).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return
}
