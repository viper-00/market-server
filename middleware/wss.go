package middleware

import (
	"market/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Wss() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Sec-WSS-Token")
		if token == "" || token != global.MARKET_CONFIG.Wss.SecWssToken {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
