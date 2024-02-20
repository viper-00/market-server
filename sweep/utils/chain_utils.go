package utils

import (
	"market/global/constant"
	"market/model"
	"market/utils"
	"strings"
)

func IsChainJoinSweep(chainId int) bool {
	if chainId == 0 {
		return false
	}

	for _, v := range constant.JoinSweep {
		if v == chainId {
			return true
		}
	}

	return false
}

// isContract, symbol, contractAddress, decimals
func GetContractInfoByChainIdAndContractAddress(chainId int, contractAddress string) (bool, string, string, int) {
	if !IsChainJoinSweep(chainId) {
		return false, "", "", 0
	}

	for _, element := range model.ChainList {
		if element.ChainId != chainId {
			continue
		}

		for _, coin := range element.Coins {
			switch chainId {
			case constant.ETH_MAINNET,
				constant.ETH_GOERLI,
				constant.ETH_SEPOLIA,
				constant.BSC_MAINNET,
				constant.BSC_TESTNET,
				constant.OP_MAINNET,
				constant.OP_SEPOLIA,
				constant.ARBITRUM_ONE,
				constant.ARBITRUM_NOVA,
				constant.ARBITRUM_GOERLI,
				constant.ARBITRUM_SEPOLIA:
				if utils.HexToAddress(coin.Contract) == utils.HexToAddress(contractAddress) {
					return true, coin.Symbol, coin.Contract, coin.Decimals
				}
			case constant.TRON_NILE, constant.TRON_MAINNET:
				if strings.EqualFold(coin.Contract, contractAddress) {
					return true, coin.Symbol, coin.Contract, coin.Decimals
				}
			case constant.BTC_MAINNET, constant.BTC_TESTNET:
				if coin.IsMainCoin {
					return true, coin.Symbol, coin.Contract, coin.Decimals
				}
			case constant.LTC_MAINNET, constant.LTC_TESTNET:
				if coin.IsMainCoin {
					return true, coin.Symbol, coin.Contract, coin.Decimals
				}
			}
		}

		return false, "", "", 0

	}
	return false, "", "", 0
}

// isContract, symbol, contractAddress, decimals
func GetContractInfoByChainIdAndSymbol(chainId int, symbol string) (bool, string, string, int) {
	if !IsChainJoinSweep(chainId) {
		return false, "", "", 0
	}

	for _, element := range model.ChainList {
		if element.ChainId != chainId {
			continue
		}

		for _, coin := range element.Coins {
			if coin.Symbol == symbol {
				return true, coin.Symbol, coin.Contract, coin.Decimals
			}
		}
	}
	return false, "", "", 0
}
