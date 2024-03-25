package service

import (
	"errors"
	"market/global"
	"market/global/constant"
	"market/model"
	"market/model/market/request"
	"market/model/market/response"

	"github.com/gin-gonic/gin"
)

func (m *MService) GetMarketEventTypeForHome(c *gin.Context) (result interface{}, err error) {
	return []string{string(constant.EVENT_ALL), string(constant.EVENT_FOR_YOU), string(constant.EVENT_CRYPTO), string(constant.EVENT_BUSINESS), string(constant.EVENT_SCIENCE), string(constant.EVENT_GAME)}, nil
}

func (m *MService) GetMarketEventForHome(c *gin.Context, req request.GetMarketEventForHome) (result []response.EventForHomeResponse, err error) {
	var searchEvent []model.Event

	var limit = 10
	switch req.Type {
	case string(constant.EVENT_ALL):
		err = global.MARKET_DB.Where("event_status = 1 AND status = 1").Order("id desc").Limit(limit).Find(&searchEvent).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	case string(constant.EVENT_FOR_YOU):
		err = global.MARKET_DB.Where("event_status = 1 AND status = 1").Order("RAND()").Limit(limit).Find(&searchEvent).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	case string(constant.EVENT_CRYPTO):
		err = global.MARKET_DB.Where("event_status = 1 AND status = 1 AND type = ?", string(constant.EVENT_CRYPTO)).Order("id desc").Limit(limit).Find(&searchEvent).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	case string(constant.EVENT_BUSINESS):
		err = global.MARKET_DB.Where("event_status = 1 AND status = 1 AND type = ?", string(constant.EVENT_BUSINESS)).Order("id desc").Limit(limit).Find(&searchEvent).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	case string(constant.EVENT_SCIENCE):
		err = global.MARKET_DB.Where("event_status = 1 AND status = 1 AND type = ?", string(constant.EVENT_SCIENCE)).Order("id desc").Limit(limit).Find(&searchEvent).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	case string(constant.EVENT_GAME):
		err = global.MARKET_DB.Where("event_status = 1 AND status = 1 AND type = ?", string(constant.EVENT_GAME)).Order("id desc").Limit(limit).Find(&searchEvent).Error
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return
		}
	default:
		return nil, errors.New("no support")
	}

	if len(searchEvent) > 0 {
		for _, v := range searchEvent {
			var r response.EventForHomeResponse
			r.EventLogo = v.EventLogo
			r.Title = v.Title
			r.ExpireTime = int(v.ExpireTime.UnixMilli())
			r.SettlementTime = int(v.SettlementTime.UnixMilli())
			r.Type = v.Type
			r.UniqueCode = v.UniqueWebsiteCode

			var orders []model.EventOrder
			err = global.MARKET_DB.Where("event_id = ? AND status = 1 AND order_status = 1", v.ID).Find(&orders).Error
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return
			}

			for _, o := range orders {
				r.TotalOrderAmount += o.Amount
			}

			var comments []model.EventComment
			err = global.MARKET_DB.Where("event_id = ? AND status = 1", v.ID).Find(&comments).Error
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return
			}
			r.CommentCount = len(comments)

			result = append(result, r)
		}
	}

	return
}
