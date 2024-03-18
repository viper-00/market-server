package erc20

import (
	"encoding/hex"
	"errors"
	"fmt"
	"market/global"
	"market/global/constant"
	sweepUtils "market/sweep/utils"
	"market/utils"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	Transfer     = "transfer"
	TransferFrom = "transferFrom"
	Approve      = "approve"

	knownMethods = map[string]string{
		"0xa9059cbb": Transfer,
		"0x23b872dd": TransferFrom,
	}
)

func IsHandleTokenTransaction(chainId int, hash, contractName, fromAddress, contractAddress, monitorAddress, data string) bool {
	if !sweepUtils.IsChainJoinSweep(chainId) {
		return false
	}

	switch contractName {
	// case constant.PREDICTMARKET:
	// 	return handlePredictMarketTransaction(chainId, hash, fromAddress, contractAddress, monitorAddress, data)
	default:
		return handleERC20Transaction(chainId, hash, fromAddress, contractAddress, monitorAddress, data)
	}

	return false
}

// func handlePredictMarketTransaction(chainId int, hash, fromAddress, contractAddress, monitorAddress, data string) bool {

// }

func handleERC20Transaction(chainId int, hash, fromAddress, contractAddress, monitorAddress, data string) bool {
	methodName, decodeFromAddress, decodeToAddress, _, err := DecodeERC20TransactionInputData(chainId, hash, data)

	if err != nil {
		return false
	}

	switch methodName {
	case Transfer:
		if utils.HexToAddress(monitorAddress) == utils.HexToAddress(fromAddress) {
			return true
		}

		if utils.HexToAddress(monitorAddress) == utils.HexToAddress(decodeToAddress) {
			return true
		}
	case TransferFrom:
		if utils.HexToAddress(monitorAddress) == utils.HexToAddress(decodeFromAddress) {
			return true
		}

		if utils.HexToAddress(monitorAddress) == utils.HexToAddress(decodeToAddress) {
			return true
		}
	}

	return false
}

func DecodeERC20TransactionInputData(chainId int, hash, data string) (methodName, decodeFromAddress, decodeToAddress string, amount *big.Int, err error) {

	if len(data) < 138 {
		err = errors.New("insufficient data length")
		return
	}

	methodSigData := data[:10]
	inputSigData := data[10:]

	decodedData, err := hex.DecodeString(inputSigData)
	if err != nil {
		err = errors.New("can not decode input sig data")
		return
	}

	file, err := os.Open("json/ERC20.json")
	if err != nil {
		err = errors.New("can not open erc20 file")
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> hash:%s, message:%s", constant.GetChainName(chainId), hash, err.Error()))
		return
	}
	defer file.Close()

	contractABI, err := abi.JSON(file)
	if err != nil {
		err = errors.New("can not from file to json of abi")
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> hash:%s, message:%s", constant.GetChainName(chainId), hash, err.Error()))
		return
	}

	methodName, isKnownMethod := knownMethods[methodSigData]
	if !isKnownMethod || (methodName != Transfer && methodName != TransferFrom) {
		err = errors.New("not a valid transfer method")
		return
	}

	method, isAbiMethod := contractABI.Methods[methodName]
	if !isAbiMethod {
		err = errors.New("transfer method not found in ABI")
		return
	}

	inputsMap := make(map[string]interface{})

	if err = method.Inputs.UnpackIntoMap(inputsMap, decodedData); err != nil {
		err = errors.New("can not decode: Unpack Into Map")
		return
	}

	switch method.Name {
	case Transfer:
		address, ok := inputsMap["_to"].(common.Address)
		if !ok {
			err = errors.New("can not get the value of _to")
			return
		}

		value, ok := inputsMap["_value"].(*big.Int)
		if !ok {
			err = errors.New("can not get the value of _value")
			return
		}

		return method.Name, "", address.Hex(), value, nil
	case TransferFrom:
		fromAddress, ok := inputsMap["_from"].(common.Address)
		if !ok {
			err = errors.New("can not get the value of _from")
			return
		}
		toAddress, ok := inputsMap["_to"].(common.Address)
		if !ok {
			err = errors.New("can not get the value of _to")
			return
		}
		value, ok := inputsMap["_value"].(*big.Int)
		if !ok {
			err = errors.New("can not get the value of _value")
			return
		}

		return method.Name, fromAddress.Hex(), toAddress.Hex(), value, nil
	}

	err = errors.New("not found method")
	return
}
