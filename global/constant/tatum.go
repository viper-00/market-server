package constant

import (
	"math/rand"
	"time"
)

var (
	TatumAPI        = "https://api.tatum.io/v3"
	TatumMainnetKey = []string{}

	TatumTestnetKey = []string{}

	// Bitcoin
	TatumGetBitcoinInfo                = TatumAPI + "/bitcoin/info"
	TatumGetBitcoinBlockByHashOrHeight = TatumAPI + "/bitcoin/block/"
	TatumGetBitcoinTxByHash            = TatumAPI + "/bitcoin/transaction/"

	// Litecoin
	TatumGetLitecoinInfo                = TatumAPI + "/litecoin/info"
	TatumGetLitecoinBlockByHashOrHeight = TatumAPI + "/litecoin/block/"
	TatumGetLitecoinTxByHash            = TatumAPI + "/litecoin/transaction/"
)

var TatumSupportChain = []int{
	BTC_MAINNET,
	BTC_TESTNET,
	ETH_MAINNET,
	ETH_SEPOLIA,
	LTC_MAINNET,
	LTC_TESTNET,
	BSC_MAINNET,
	BSC_TESTNET,
	TRON_MAINNET,
}

func IsNetworkSupportTatum(id int) bool {
	for _, v := range TatumSupportChain {
		if v == id {
			return true
		}
	}

	return false
}

func GetTatumRandomKeyByNetwork(id int) string {
	rand.Seed(time.Now().UnixNano())

	switch id {
	case BTC_MAINNET, ETH_MAINNET, LTC_MAINNET, BSC_MAINNET, TRON_MAINNET:
		index := rand.Intn(len(TatumMainnetKey))
		return TatumMainnetKey[index]
	case BTC_TESTNET, ETH_SEPOLIA, LTC_TESTNET, BSC_TESTNET:
		index := rand.Intn(len(TatumTestnetKey))
		return TatumTestnetKey[index]
	}

	return ""
}

func GetAllTatumAPiKey(id int) []string {
	switch id {
	case BTC_MAINNET, ETH_MAINNET, LTC_MAINNET, BSC_MAINNET, TRON_MAINNET:
		return TatumMainnetKey
	case BTC_TESTNET, ETH_SEPOLIA, LTC_TESTNET, BSC_TESTNET:
		return TatumTestnetKey

	}
	return nil
}
