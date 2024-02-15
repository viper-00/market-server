package constant

import (
	"math/rand"
	"time"
)

var (
	TrongridMainnetAPI = "https://api.trongrid.io"
	TrongridNileAPI    = "https://nile.trongrid.io"

	TrongridMainnetKey = []string{}

	TrongridNileKey = []string{
		"709b2085-98df-41ba-a680-5eda56806430",
	}
)

func GetHTTPUrlByNetwork(network int) string {
	rand.Seed(time.Now().UnixNano())

	switch network {
	case TRON_MAINNET:
		return TrongridMainnetAPI
	case TRON_NILE:
		return TrongridNileAPI
	}

	return ""
}

func GetRandomHTTPKeyByNetwork(network int) string {
	rand.Seed(time.Now().UnixNano())

	switch network {
	case TRON_MAINNET:
		index := rand.Intn(len(TrongridMainnetKey))
		return TrongridMainnetKey[index]
	case TRON_NILE:
		index := rand.Intn(len(TrongridNileKey))
		return TrongridNileKey[index]
	}

	return ""
}

func GetAllHTTPKeyByNetwork(network int) []string {
	switch network {
	case TRON_MAINNET:
		return TrongridMainnetKey
	case TRON_NILE:
		return TrongridNileKey
	}

	return nil
}

func TronGetBlockByNetwork(network int) string {
	switch network {
	case TRON_MAINNET:
		return TrongridMainnetAPI + "/walletsolidity/getblock"
	case TRON_NILE:
		return TrongridNileAPI + "/walletsolidity/getblock"
	}

	return ""
}

func TronGetBlockByNumByNetwork(network int) string {
	switch network {
	case TRON_MAINNET:
		return TrongridMainnetAPI + "/walletsolidity/getblockbynum"
	case TRON_NILE:
		return TrongridNileAPI + "/walletsolidity/getblockbynum"
	}

	return ""
}

func TronGetTxByIdByNetwork(network int) string {
	switch network {
	case TRON_MAINNET:
		return TrongridMainnetAPI + "/walletsolidity/gettransactionbyid"
	case TRON_NILE:
		return TrongridNileAPI + "/walletsolidity/gettransactionbyid"
	}

	return ""
}

func TronValidateAddressByNetwork(network int) string {
	switch network {
	case TRON_MAINNET:
		return TrongridMainnetAPI + "/wallet/validateaddress"
	case TRON_NILE:
		return TrongridNileAPI + "/wallet/validateaddress"
	}

	return ""
}

func TronValidateContractAddressByNetwork(network int) string {
	switch network {
	case TRON_MAINNET:
		return TrongridMainnetAPI + "/wallet/getcontract"
	case TRON_NILE:
		return TrongridNileAPI + "/wallet/getcontract"
	}

	return ""
}
