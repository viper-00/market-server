package service

import (
	"market/global"
	"market/model"
	"market/model/market/request"

	"github.com/gin-gonic/gin"
)

func (m *MService) CreateEventComment(c *gin.Context, req request.CreateEventComment) (result interface{}, err error) {
	userModel, err := m.GetUserInfo(c)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	var eventModel model.Event
	err = global.MARKET_DB.Where("unique_website_code = ? AND event_status = 1 AND status = 1", req.Code).First(&eventModel).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	var model model.EventComment
	model.Content = req.Content
	model.UserId = userModel.ID
	model.ReplyId = req.ReplyId
	model.EventId = eventModel.ID
	err = global.MARKET_DB.Save(&model).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	return model, nil
}

func (m *MService) GetEventComment(c *gin.Context, req request.GetEventComment) (result interface{}, err error) {
	var eventModel model.Event
	err = global.MARKET_DB.Where("unique_website_code = ? AND event_status = 1 AND status = 1", req.Code).First(&eventModel).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	var eventComments []model.EventComment
	err = global.MARKET_DB.Where("event_id = ? AND status = 1", eventModel.ID).Find(&eventComments).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}
	return eventComments, nil
}

func (m *MService) RemoveEventComment(c *gin.Context, req request.RemoveEventComment) (err error) {
	err = global.MARKET_DB.Where("id = ?", req.CommentId).Delete(&model.EventComment{}).Error
	return
}
