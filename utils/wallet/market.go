package wallet

import (
	"errors"
	"fmt"
	"market/global"
	"market/global/constant"
	"market/model/market/response"
	sweepUtils "market/sweep/utils"
	"market/utils"
	"math/big"
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

	isSupport, _, contractAddress, _ = sweepUtils.GetContractInfoByChainIdAndSymbol(chainId, constant.PREDICTMARKET)
	if !isSupport {
		return "", errors.New("contract address not found")
	}

	generalPriAccountm, generalPubAccount, err := GetGeneralAccountByChainId(chainId)
	if err != nil {
		return "", errors.New("contract address not found")
	}

	bindAddresses := []string{ownerPublicKey, generalPubAccount}

	hash, err := CreateNewCollectionContract(rpc, generalPriAccountm, generalPubAccount, contractAddress, bindAddresses, gasLimit)

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
		return global.MARKET_CONFIG.GeneralAccount.OpSepolia.PrivateKey, global.MARKET_CONFIG.GeneralAccount.OpSepolia.PublicKey, nil
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

// func GetCallWithdrawContractGasLimitFromChainId(chainId int) (bool, uint64) {
// 	switch chainId {
// 	case constant.OP_SEPOLIA:
// 		return true, 60000
// 	}

// 	return false, 0
// }

func GetNewContractAddressByTxHash(chainId int, hash, callContractAddress string) (contractAddress string, err error) {
	rpc := constant.GetRPCUrlByNetwork(chainId)
	if rpc == "" {
		return "", errors.New("chain not support")
	}

	var receipt *types.Receipt

	for {
		receipt, err = GetTransactionReceiptByHash(rpc, hash)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 1)
		global.MARKET_LOG.Info(fmt.Sprintf("retry the GetTransactionByHash, hash: %s, chainId: %d", hash, chainId))
	}

	if receipt.Status != 1 {
		return "", errors.New("transaction not included in block")
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

func MonitorTxStatus(chainId int, hash string) (err error) {
	rpc := constant.GetRPCUrlByNetwork(chainId)
	if rpc == "" {
		return errors.New("chain not support")
	}

	var receipt *types.Receipt

	for {
		receipt, err = GetTransactionReceiptByHash(rpc, hash)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 1)
		// global.MARKET_LOG.Info(fmt.Sprintf("retry the MonitorTxStatus, hash: %s, chainId: %d", hash, chainId))
	}

	if receipt.Status == 1 {
		return nil
	} else {
		return errors.New("transaction failed")
	}
}

func GetUsdtTokenBalance(chainId int, address string) (balance string, err error) {
	isSupport, _, contractAddress, decimals := sweepUtils.GetContractInfoByChainIdAndSymbol(chainId, constant.USDT)
	if !isSupport {
		return "", errors.New("contract address not found")
	}

	return GetSingleTokenBalance(chainId, decimals, contractAddress, address)
}

func GetSingleTokenBalance(chainId, decimals int, contractAddress, address string) (balance string, err error) {
	rpc := constant.GetRPCUrlByNetwork(chainId)
	if rpc == "" {
		return "", errors.New("chain not support")
	}

	intBalance, err := CallTokenBalanceOf(rpc, address, contractAddress)
	if err != nil {
		return "", err
	}

	balance = utils.CalculateBalance(intBalance, decimals)

	return
}

func GetAllTokenBalance(chainId int, address string) (balance response.UserBalance, err error) {
	rpc := constant.GetRPCUrlByNetwork(chainId)
	if rpc == "" {
		return balance, errors.New("chain not support")
	}

	// eth, usdt, usdc
	ethBalance, err := GetEthBalanceByAddress(rpc, address)
	if err != nil {
		return
	}

	chainInfo := sweepUtils.GetOneChainInfoByChainId(chainId)

	for _, v := range chainInfo.Coins {
		if v.Symbol == constant.ETH {
			balance.ETH = utils.CalculateBalance(ethBalance, v.Decimals)
		} else if v.Symbol == constant.USDT {
			usdtBalance, usdtErr := CallTokenBalanceOf(rpc, address, v.Contract)
			if usdtErr != nil {
				return balance, usdtErr
			}
			balance.USDT = utils.CalculateBalance(usdtBalance, v.Decimals)
		} else if v.Symbol == constant.USDC {
			usdcBalance, usdcErr := CallTokenBalanceOf(rpc, address, v.Contract)
			if usdcErr != nil {
				return balance, usdcErr
			}
			balance.USDC = utils.CalculateBalance(usdcBalance, v.Decimals)
		}
	}

	return balance, nil
}

func TransferAssetToReceiveAddress(chainId int, callContractAddress string, sendVal float64) (hash string, err error) {
	rpc := constant.GetRPCUrlByNetwork(chainId)
	if rpc == "" {
		return "", errors.New("chain not support")
	}

	isSupport, _, contractAddress, decimals := sweepUtils.GetContractInfoByChainIdAndSymbol(chainId, constant.USDT)
	if !isSupport {
		return "", errors.New("contract address not found")
	}

	tokenAddresses := []string{contractAddress}
	sendToAddresses := []string{global.MARKET_CONFIG.GeneralAccount.OpSepolia.ReceiveAccount}
	sendValues := []big.Int{*big.NewInt(utils.FormatToOriginalValue(sendVal, decimals))}

	var gasLimit uint64 = 100000

	hash, err = CallWithdrawByCollectionContract(rpc, global.MARKET_CONFIG.GeneralAccount.OpSepolia.PrivateKey, global.MARKET_CONFIG.GeneralAccount.OpSepolia.PublicKey, callContractAddress, tokenAddresses, sendToAddresses, sendValues, gasLimit)
	if err != nil {
		return "", err
	}

	err = MonitorTxStatus(chainId, hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func TransferAssetToMoreReceiveAddres(chainId int, tokenAddresses, sendToAddresses []string, sendValues []big.Int) (hash string, err error) {

	rpc := constant.GetRPCUrlByNetwork(chainId)
	if rpc == "" {
		return "", errors.New("chain not support")
	}

	var gasLimit uint64 = 200000

	hash, err = CallWithdrawByCollectionContract(rpc, global.MARKET_CONFIG.GeneralAccount.OpSepolia.PrivateKey, global.MARKET_CONFIG.GeneralAccount.OpSepolia.PublicKey, global.MARKET_CONFIG.GeneralAccount.OpSepolia.ReceiveAccount, tokenAddresses, sendToAddresses, sendValues, gasLimit)
	if err != nil {
		return "", err
	}

	err = MonitorTxStatus(chainId, hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func TransferEthToReceiveAddress(chainId int, sendToAddress string, sendVal float64) (hash string, err error) {
	rpc := constant.GetRPCUrlByNetwork(chainId)
	if rpc == "" {
		return "", errors.New("chain not support")
	}

	var gasLimit uint64 = 96000
	sendValue := big.NewInt(utils.FormatToOriginalValue(sendVal, 18))

	hash, err = CallEthTransfer(rpc, global.MARKET_CONFIG.GeneralAccount.OpSepolia.PrivateKey, global.MARKET_CONFIG.GeneralAccount.OpSepolia.PublicKey, sendToAddress, sendValue, gasLimit)
	if err != nil {
		return "", err
	}

	err = MonitorTxStatus(chainId, hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func TransferTokenToReceiveAddress(chainId int, sendToAddress, coin string, sendVal float64) (hash string, err error) {
	rpc := constant.GetRPCUrlByNetwork(chainId)
	if rpc == "" {
		return "", errors.New("chain not support")
	}

	isSupport, _, contractAddress, decimals := sweepUtils.GetContractInfoByChainIdAndSymbol(chainId, coin)
	if !isSupport {
		return "", errors.New("contract address not found")
	}

	var gasLimit uint64 = 96000
	sendValue := big.NewInt(utils.FormatToOriginalValue(sendVal, decimals))

	hash, err = CallTokenTransfer(rpc, global.MARKET_CONFIG.GeneralAccount.OpSepolia.PrivateKey, global.MARKET_CONFIG.GeneralAccount.OpSepolia.PublicKey, sendToAddress, contractAddress, sendValue, gasLimit)
	if err != nil {
		return "", err
	}

	err = MonitorTxStatus(chainId, hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}
