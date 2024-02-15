package utils

import (
	"fmt"
	"market/global"
)

func HandlePanic() {
	if r := recover(); r != nil {
		global.MARKET_LOG.Error(fmt.Sprint(r))
		return
	}
}
