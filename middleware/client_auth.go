package middleware

import (
	"context"
	"market/global"
	"market/utils/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ClientAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		emailorAddressResult, err := global.MARKET_REDIS.Get(context.Background(), authorization).Result()
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		uniqueBearer := strings.Split(authorization, " ")[1]

		claims, err := jwt.ValidateJWT(uniqueBearer)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		chainId := claims["chain_id"].(float64)
		email := claims["email"].(string)
		address := claims["address"].(string)
		contractAddress := claims["contractAddress"].(string)
		time := claims["time"].(float64)

		if emailorAddressResult != email && emailorAddressResult != address {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if chainId == 0 || address == "" || contractAddress == "" || time == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("chainId", chainId)
		c.Set("email", email)
		c.Set("address", address)
		c.Set("contractAddress", contractAddress)
		c.Set("time", time)

		c.Next()
	}
}
