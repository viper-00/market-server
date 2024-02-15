package tron

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"market/global"
	"market/global/constant"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func TronDecodeMethod(data string) (string, string, error) {
	if len(data) < 8 {
		return "", "", errors.New("length too short for tron decode method")
	}
	return data[:8], data[8:], nil
}

var (
	Transfer     = "transfer"
	TransferFrom = "transferFrom"

	BatchTransferFrom  = "batchTransferFrom"
	Collect            = "collect"
	SenderTransferFrom = "senderTransferFrom"
	Withdraw           = "withdraw"

	TransferContract     = "TransferContract"
	TriggerSmartContract = "TriggerSmartContract"

	KnownMethods = map[string]string{
		"a9059cbb": Transfer,
		"23b872dd": TransferFrom,

		"b818f9e4": BatchTransferFrom,
		"1e13eee1": Collect,
		"b19385f7": SenderTransferFrom,
		"f7ece0cf": Withdraw,
	}

	joinAllInOneContract = []string{BatchTransferFrom, Collect, SenderTransferFrom, Withdraw}
)

var base58Alphabets = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func hasJoinAllInOneContract(methodName string) bool {
	for _, v := range joinAllInOneContract {
		if v == methodName {
			return true
		}
	}

	return false
}

func ToHexAddress(address string) string {
	return hex.EncodeToString(base58Decode([]byte(address)))
}

func FromHexAddress(hexAddress string) (string, error) {
	addrByte, err := hex.DecodeString(hexAddress)
	if err != nil {
		return "", err
	}

	sha := sha256.New()
	sha.Write(addrByte)
	shaStr := sha.Sum(nil)

	sha2 := sha256.New()
	sha2.Write(shaStr)
	shaStr2 := sha2.Sum(nil)

	addrByte = append(addrByte, shaStr2[:4]...)
	return string(base58Encode(addrByte)), nil
}

func base58Encode(input []byte) []byte {
	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := &big.Int{}
	var result []byte
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, base58Alphabets[mod.Int64()])
	}
	reverseBytes(result)
	return result
}

func base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	for _, b := range input {
		charIndex := bytes.IndexByte(base58Alphabets, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}
	decoded := result.Bytes()
	if input[0] == base58Alphabets[0] {
		decoded = append([]byte{0x00}, decoded...)
	}
	return decoded[:len(decoded)-4]
}

func reverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func IsHandleTransaction(chainId int, hash, contractType, fromAddress, toAddress, monitorAddress, data string) bool {
	if fromAddress == "" || toAddress == "" {
		return false
	}

	switch contractType {
	case TransferContract:
		return handleTransferContract(fromAddress, toAddress, monitorAddress)
	case TriggerSmartContract:
		return handleTriggerSmartContract(chainId, hash, fromAddress, toAddress, monitorAddress, data)
	}

	return false
}

func handleTransferContract(fromAddress, toAddress, monitorAddress string) bool {
	sendAddress, err := FromHexAddress(fromAddress)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return false
	}
	receiveAddress, err := FromHexAddress(toAddress)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return false
	}
	if strings.EqualFold(sendAddress, monitorAddress) || strings.EqualFold(receiveAddress, monitorAddress) {
		return true
	}
	return false
}

func handleTriggerSmartContract(chainId int, hash, fromAddress, toAddress, monitorAddress, data string) bool {

	methodID, _, _ := TronDecodeMethod(data)

	method := KnownMethods[methodID]
	if method == "" {
		return false
	}

	switch method {
	case Transfer:
		sendAddress, err := FromHexAddress(fromAddress)
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return false
		}
		receiveAddress, err := FromHexAddress("41" + data[32:72])
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return false
		}

		if strings.EqualFold(sendAddress, monitorAddress) || strings.EqualFold(receiveAddress, monitorAddress) {
			return true
		}
	case TransferFrom:
		sendAddress, err := FromHexAddress("41" + data[32:72])
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return false
		}
		receiveAddress, err := FromHexAddress("41" + data[96:136])
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			return false
		}

		if strings.EqualFold(sendAddress, monitorAddress) || strings.EqualFold(receiveAddress, monitorAddress) {
			return true
		}
	case BatchTransferFrom, Collect, SenderTransferFrom, Withdraw:
		return handleAllInOneTransaction(chainId, hash, fromAddress, toAddress, monitorAddress, data)
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
			if strings.EqualFold(monitorAddress, from) {
				return true
			}
		}

		for _, to := range toAddresses {
			if strings.EqualFold(monitorAddress, to) {
				return true
			}
		}
	}

	return false
}

func DecodeAllInOneTransactionInputData(chainId int, hash, fromAddress, toAddress, data string) (string, []string, []string, []string, []*big.Int, error) {
	var err error

	methodSigData := data[:8]
	inputSigData := data[8:]

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

	methodName, isKnownMethod := KnownMethods[methodSigData]
	if !isKnownMethod || !hasJoinAllInOneContract(methodName) {
		return "", nil, nil, nil, nil, errors.New("not a valid method")
	}

	method, isAbiMethod := contractABI.Methods[methodName]
	if !isAbiMethod {
		return "", nil, nil, nil, nil, errors.New("method not found in ABI")
	}

	inputsMap := make(map[string]interface{})

	if err = method.Inputs.UnpackIntoMap(inputsMap, decodedData); err != nil {
		return "", nil, nil, nil, nil, err
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

		var (
			decodeFromAddresses  = make([]string, totalLength)
			decodeToAddresses    = make([]string, totalLength)
			decodeTokenAddresses = make([]string, totalLength)
		)

		for i := 0; i < totalLength; i++ {
			decodeFromAddresses[i], err = FromHexAddress("41" + fromAddresses[i].Hex()[2:])
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
			}
			decodeToAddresses[i], err = FromHexAddress(toAddress)
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
			}

			decodeTokenAddresses[i], err = FromHexAddress("41" + tokens[i].Hex()[2:])
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
			}
		}

		return methodName, decodeFromAddresses, decodeToAddresses, decodeTokenAddresses, amounts, nil
	case Collect:
		tokens, ok := inputsMap["_tokens"].([]common.Address)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the token of _tokens")
		}

		to, ok := inputsMap["to"].(common.Address)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the toAddress of to")
		}

		amounts, ok := inputsMap["_amounts"].([]*big.Int)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the amount of _amounts")
		}

		var totalLength = len(tokens)
		var (
			decodeFromAddresses  = make([]string, totalLength)
			decodeToAddresses    = make([]string, totalLength)
			decodeTokenAddresses = make([]string, totalLength)
		)

		for i := 0; i < totalLength; i++ {
			decodeFromAddresses[i], err = FromHexAddress(toAddress)
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
			}
			decodeToAddresses[i], err = FromHexAddress("41" + to.Hex()[2:])
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
			}
			decodeTokenAddresses[i], err = FromHexAddress("41" + tokens[i].Hex()[2:])
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
			}
		}

		return methodName, decodeFromAddresses, decodeToAddresses, decodeTokenAddresses, amounts, nil
	case SenderTransferFrom:
		tokens, ok := inputsMap["_tokens"].([]common.Address)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the token of _tokens")
		}

		amounts, ok := inputsMap["_amounts"].([]*big.Int)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the amount of _amounts")
		}

		var totalLength = len(tokens)
		var (
			decodeFromAddresses  = make([]string, totalLength)
			decodeToAddresses    = make([]string, totalLength)
			decodeTokenAddresses = make([]string, totalLength)
		)

		for i := 0; i < totalLength; i++ {
			decodeFromAddresses[i], err = FromHexAddress(fromAddress)
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
			}

			decodeToAddresses[i], err = FromHexAddress(toAddress)
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
			}

			decodeTokenAddresses[i], err = FromHexAddress("41" + tokens[i].Hex()[2:])
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
			}
		}

		return methodName, decodeFromAddresses, decodeToAddresses, decodeTokenAddresses, amounts, nil
	case Withdraw:
		tokens, ok := inputsMap["_tokens"].([]common.Address)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the token of _tokens")
		}

		toAddresses, ok := inputsMap["_tos"].([]common.Address)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the toAddress of _tos")
		}

		amounts, ok := inputsMap["_amounts"].([]*big.Int)
		if !ok {
			return "", nil, nil, nil, nil, errors.New("can not get the amount of _amounts")
		}

		var totalLength = len(tokens)
		var (
			decodeFromAddresses  = make([]string, totalLength)
			decodeToAddresses    = make([]string, totalLength)
			decodeTokenAddresses = make([]string, totalLength)
		)

		for i := 0; i < totalLength; i++ {
			decodeFromAddresses[i], err = FromHexAddress(toAddress)
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
			}

			decodeToAddresses[i], err = FromHexAddress("41" + toAddresses[i].Hex()[2:])
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
			}

			decodeTokenAddresses[i], err = FromHexAddress("41" + tokens[i].Hex()[2:])
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
			}
		}

		return methodName, decodeFromAddresses, decodeToAddresses, decodeTokenAddresses, amounts, nil
	}

	return "", nil, nil, nil, nil, errors.New("not found method")
}
