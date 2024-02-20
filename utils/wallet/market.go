package wallet

import (
	"errors"
	"fmt"
	"market/global"
	"market/global/constant"
	"market/sweep/utils"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func GenerateEthereumCollectionContract(chainId int, ownerPublicKey string) (contractAddress string, err error) {
	var gasLimit uint64 = 0

	rpc := constant.GetRPCUrlByNetwork(chainId)
	if rpc == "" {
		return "", errors.New("chain not support")
	}

	isSupport, gasLimit := GetCallCreateContractGasLimitFromChainId(chainId)
	if !isSupport {
		return "", errors.New("chain not support")
	}

	isSupport, _, contractAddress, _ = utils.GetContractInfoByChainIdAndSymbol(chainId, constant.PREDICTMARKET)
	if !isSupport {
		return "", errors.New("contract address not found")
	}

	generalPriAccountm, generalPubAccount, err := GetGeneralAccountByChainId(chainId)
	if err != nil {
		return "", errors.New("contract address not found")
	}

	bindAddresses := []string{ownerPublicKey, generalPubAccount}

	hash, err := SendEthereumCollectionContract(rpc, generalPriAccountm, generalPubAccount, contractAddress, bindAddresses, gasLimit)

	if err != nil {
		return "", err
	}

	newContractAddress, err := GetNewContractAddressByTxHash(chainId, hash, contractAddress)
	if err != nil {
		return "", err
	}

	return newContractAddress, nil
}

func GetGeneralAccountByChainId(chainId int) (string, string, error) {
	switch chainId {
	case constant.OP_MAINNET, constant.OP_SEPOLIA, constant.OP_GOERLI:
		return global.MARKET_CONFIG.GeneralAccount.Op.PrivateKey, global.MARKET_CONFIG.GeneralAccount.Op.PublicKey, nil
	}

	return "", "", errors.New("not found the account")
}

func GetCallCreateContractGasLimitFromChainId(chainId int) (bool, uint64) {
	switch chainId {
	case constant.OP_SEPOLIA:
		return true, 1000000
	}

	return false, 0
}

func GetNewContractAddressByTxHash(chainId int, hash, callContractAddress string) (contractAddress string, err error) {
	rpc := constant.GetRPCUrlByNetwork(chainId)
	if rpc == "" {
		return "", errors.New("chain not support")
	}

	var receipt *types.Receipt

	for {
		receipt, err = GetTransactionByHash(rpc, hash)
		if err == nil {
			break
		}
		time.Sleep(1)
		global.MARKET_LOG.Info(fmt.Sprintf("retry the GetTransactionByHash, hash: %s, chainId: %d", hash, chainId))
	}

	if receipt.Status != 1 {
		return "", errors.New("Transaction not included in block")
	}

	for _, v := range receipt.Logs {
		if common.HexToAddress(v.Topics[0].Hex()) == common.HexToAddress("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0") &&
			common.HexToAddress(v.Topics[1].Hex()) == common.HexToAddress("0x0000000000000000000000000000000000000000") &&
			common.HexToAddress(v.Topics[2].Hex()) == common.HexToAddress(callContractAddress) {
			return v.Address.String(), nil
		}
	}

	return "", errors.New("tx not support")
}
