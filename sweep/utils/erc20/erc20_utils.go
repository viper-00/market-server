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

	BatchTransferFrom  = "batchTransferFrom"
	Collect            = "collect"
	SenderTransferFrom = "senderTransferFrom"
	Withdraw           = "withdraw"

	knownMethods = map[string]string{
		"0xa9059cbb": Transfer,
		"0x23b872dd": TransferFrom,

		"0xb818f9e4": BatchTransferFrom,
		"0x1e13eee1": Collect,
		"0xb19385f7": SenderTransferFrom,
		"0xf7ece0cf": Withdraw,
	}

	joinAllInOneContract = []string{BatchTransferFrom, Collect, SenderTransferFrom, Withdraw}
)

func hasJoinAllInOneContract(methodName string) bool {
	for _, v := range joinAllInOneContract {
		if v == methodName {
			return true
		}
	}

	return false
}

func IsHandleTokenTransaction(chainId int, hash, contractName, fromAddress, contractAddress, monitorAddress, data string) bool {
	if !sweepUtils.IsChainJoinSweep(chainId) {
		return false
	}

	switch contractName {
	case constant.ALLINONE:
		return handleAllInOneTransaction(chainId, hash, fromAddress, contractAddress, monitorAddress, data)
	case constant.SWAP:
		break
	default:
		return handleERC20Transaction(chainId, hash, fromAddress, contractAddress, monitorAddress, data)
	}

	return false
}

func handleAllInOneTransaction(chainId int, hash, fromAddress, contractAddress, monitorAddress, data string) bool {
	methodName, fromAddresses, toAddresses, tokens, amounts, err := DecodeAllInOneTransactionInputData(chainId, hash, fromAddress, contractAddress, data)

	if err != nil {
		return false
	}

	if methodName != "" && len(fromAddresses) == len(toAddresses) && len(toAddresses) == len(tokens) && len(tokens) == len(amounts) {
		for _, from := range fromAddresses {
			if utils.HexToAddress(monitorAddress) == from.String() {
				return true
			}
		}

		for _, to := range toAddresses {
			if utils.HexToAddress(monitorAddress) == to.String() {
				return true
			}
		}
	}

	return false
}

func DecodeAllInOneTransactionInputData(chainId int, hash, fromAddress, toAddress, data string) (string, []common.Address, []common.Address, []common.Address, []*big.Int, error) {
	var err error

	methodSigData := data[:10]
	inputSigData := data[10:]

	decodedData, err := hex.DecodeString(inputSigData)
	if err != nil {
		return "", nil, nil, nil, nil, err
	}

	file, err := os.Open("json/AllInOne.json")
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> hash:%s, message:%s", constant.GetChainName(chainId), hash, err.Error()))
		return "", nil, nil, nil, nil, err
	}
	defer file.Close()

	contractABI, err := abi.JSON(file)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> hash:%s, message:%s", constant.GetChainName(chainId), hash, err.Error()))
		return "", nil, nil, nil, nil, err
	}

	methodName, isKnownMethod := knownMethods[methodSigData]
	if !isKnownMethod || !hasJoinAllInOneContract(methodName) {
		return "", nil, nil, nil, nil, errors.New("not a valid method")
	}

	method, isAbiMethod := contractABI.Methods[methodName]
	if !isAbiMethod {
		return "", nil, nil, nil, nil, errors.New("method not found in ABI")
	}

	inputsMap := make(map[string]interface{})

	if err = method.Inputs.UnpackIntoMap(inputsMap, decodedData); err != nil {
		return "", nil, nil, nil, nil, errors.New("can not decode: Unpack Into Map")
	}

	switch method.Name {
	case BatchTransferFrom:
		tokens, ok := inputsMap["_tokens"].([]common.Address)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the token of _tokens")
		}

		var totalLength = len(tokens)

		fromAddresses, ok := inputsMap["_froms"].([]common.Address)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the address of _froms")
		}

		amounts, ok := inputsMap["_amounts"].([]*big.Int)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the amount of _amounts")
		}

		toAddresses := make([]common.Address, totalLength)
		for i := range toAddresses {
			toAddresses[i] = common.HexToAddress(toAddress)
		}

		return method.Name, fromAddresses, toAddresses, tokens, amounts, nil
	case Collect:
		tokens, ok := inputsMap["_tokens"].([]common.Address)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the token of _tokens")
		}

		var totalLength = len(tokens)

		to, ok := inputsMap["to"].(common.Address)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the toAddress of to")
		}

		amounts, ok := inputsMap["_amounts"].([]*big.Int)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the amount of _amounts")
		}

		fromAddresses := make([]common.Address, totalLength)
		toAddresses := make([]common.Address, totalLength)
		for i := range fromAddresses {
			fromAddresses[i] = common.HexToAddress(toAddress)
			toAddresses[i] = to
		}

		return method.Name, fromAddresses, toAddresses, tokens, amounts, nil
	case SenderTransferFrom:
		tokens, ok := inputsMap["_tokens"].([]common.Address)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the token of _tokens")
		}

		var totalLength = len(tokens)

		amounts, ok := inputsMap["_amounts"].([]*big.Int)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the amount of _amounts")
		}

		fromAddresses := make([]common.Address, totalLength)
		toAddresses := make([]common.Address, totalLength)
		for i := range fromAddresses {
			fromAddresses[i] = common.HexToAddress(fromAddress)
			toAddresses[i] = common.HexToAddress(toAddress)
		}

		return method.Name, fromAddresses, toAddresses, tokens, amounts, nil
	case Withdraw:
		tokens, ok := inputsMap["_tokens"].([]common.Address)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the token of _tokens")
		}

		var totalLength = len(tokens)

		toAddresses, ok := inputsMap["_tos"].([]common.Address)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the toAddress of _tos")
		}
		amounts, ok := inputsMap["_amounts"].([]*big.Int)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the amount of _amounts")
		}

		fromAddresses := make([]common.Address, totalLength)
		for i := range fromAddresses {
			fromAddresses[i] = common.HexToAddress(toAddress)
		}

		return method.Name, fromAddresses, toAddresses, tokens, amounts, nil
	}

	return "", nil, nil, nil, nil, errors.New("not found method")
}

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
