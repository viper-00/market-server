package service

import (
	"errors"
	"market/global"
	"market/model"
	"market/model/market/request"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (m *MService) CreateCommentLike(c *gin.Context, req request.CreateCommentLike) (err error) {
	userModel, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	var sourceModel model.EventCommentLike
	err = global.MARKET_DB.Where("comment_id = ? AND user_id = ? AND status = 1", req.CommentId, userModel.ID).First(&sourceModel).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// save
		var saveModel model.EventCommentLike
		saveModel.CommentId = req.CommentId
		saveModel.IsLike = 1
		saveModel.UserId = userModel.ID
		saveModel.Status = 1
		err = global.MARKET_DB.Save(&saveModel).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	} else {
		// update
		var isLike int
		if sourceModel.IsLike == 1 {
			isLike = 2
		} else {
			isLike = 1
		}

		err = global.MARKET_DB.Model(&model.EventCommentLike{}).Where("id = ?", sourceModel.ID).Updates(map[string]interface{}{
			"is_like": isLike,
		}).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
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
