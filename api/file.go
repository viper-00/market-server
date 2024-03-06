package api

import (
	"fmt"
	"market/global"
	"market/model/common"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (n *MarketApi) UploadFile(c *gin.Context) {
	var res common.Response

	file, err := c.FormFile("file")
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	fileElement := strings.Split(filepath.Base(file.Filename), ".")

	// name := fileElement[0]
	extension := fileElement[len(fileElement)-1]

	if extension == "" {
		res = common.FailWithMessage("File extension not found")
		c.JSON(http.StatusOK, res)
		return
	}

	currentTime := time.Now()

	fileName := fmt.Sprintf("market_%d%02d%02d%02d%02d%02d",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second())

	saveName := fileName + "." + extension

	if err := c.SaveUploadedFile(file, "./resource/images/"+saveName); err != nil {
		global.MARKET_LOG.Error(err.Error())
		res = common.FailWithMessage(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res = common.OKWithData(gin.H{
		"file_url": global.MARKET_CONFIG.File.ImageUrl + saveName,
	})

	c.JSON(http.StatusOK, res)
}
