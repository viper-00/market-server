package constant

import (
	"encoding/json"
	"errors"
	"market/global"
	"market/model"
	"os"
)

func UpdateChainListFromFile() (err error) {

	file, err := os.Open("json/ChainList.json")
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		err = errors.New("can not open chainlist file")
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&model.ChainList)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		err = errors.New("can not from file to json")
		return
	}

	return nil
}
