package constant

func MempoolGetBlockHeightByNetwork(network int) string {
	switch network {
	case BTC_MAINNET:
		return "https://mempool.space/api/blocks/tip/height"
	case BTC_TESTNET:
		return "https://mempool.space/testnet/api/blocks/tip/height"
	case LTC_MAINNET:
		return "https://litecoinspace.org/api/blocks/tip/height"
	}

	return ""
}

func MempoolGetBlockTransactionByNetwork(network int) string {
	switch network {
	case BTC_MAINNET:
		return "https://mempool.space/api/block/%s/txs/%d"
	case BTC_TESTNET:
		return "https://mempool.space/testnet/api/block/%s/txs/%d"
	case LTC_MAINNET:
		return "https://litecoinspace.org/api/block/%s/txs/%d"
	}

	return ""
}

func MempoolGetBlockHashByNetwork(network int) string {
	switch network {
	case BTC_MAINNET:
		return "https://mempool.space/api/block-height/%d"
	case BTC_TESTNET:
		return "https://mempool.space/testnet/api/block-height/%d"
	case LTC_MAINNET:
		return "https://litecoinspace.org/api/block-height/%d"
	}

	return ""
}

func MempoolGetBlockByNetwork(network int) string {
	switch network {
	case BTC_MAINNET:
		return "https://mempool.space/api/block/%s"
	case BTC_TESTNET:
		return "https://mempool.space/testnet/api/block/%s"
	case LTC_MAINNET:
		return "https://litecoinspace.org/api/block/%s"
	}

	return ""
}

func MempoolGetTransctionByNetwork(network int) string {
	switch network {
	case BTC_MAINNET:
		return "https://mempool.space/api/tx/%s"
	case BTC_TESTNET:
		return "https://mempool.space/testnet/api/tx/%s"
	case LTC_MAINNET:
		return "https://litecoinspace.org/api/tx/%s"
	}

	return ""
}
