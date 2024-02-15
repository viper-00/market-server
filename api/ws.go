package api

import (
	"context"
	"encoding/json"
	"market/global"
	"market/global/constant"
	"market/service"
	"market/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	// ReadBufferSize:  0,
	// WriteBufferSize: 0,
}

func (n *MarketApi) WsForTxInfo(c *gin.Context) {
	defer utils.HandlePanic()

	c.Request.Header.Add("Connection", "upgrade")
	c.Request.Header.Add("Upgrade", "websocket")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
	}
	defer conn.Close()

	done := make(chan struct{})

	go func(conn *websocket.Conn) {
		defer utils.HandlePanic()
		defer close(done)

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return
			}

			err = conn.WriteMessage(messageType, message)
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				return
			}
		}
	}(conn)

	for {
		select {
		case <-done:
			return
		default:
			crycleTask(conn)
		}
	}
}

func crycleTask(conn *websocket.Conn) {
	global.MARKET_MUTEX.Lock()
	defer global.MARKET_MUTEX.Unlock()

	time.Sleep(500 * time.Millisecond)

	id, err := global.MARKET_REDIS.LIndex(context.Background(), constant.WS_NOTIFICATION, 0).Result()
	if err != nil {
		return
	}

	tx, err := service.MarketService.GetOwnTxById(id)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	encodeData, err := json.Marshal(tx)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	err = conn.WriteMessage(1, encodeData)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}

	_, err = global.MARKET_REDIS.LPop(context.Background(), constant.WS_NOTIFICATION).Result()
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}
}
