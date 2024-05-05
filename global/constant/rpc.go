package constant

import (
	"market/model/market/request"
	"market/model/market/response"
	"market/utils"
	"math/rand"
	"time"
)

var (
	AlchemyMainnetKey = []string{
		"",
	}

	AlchemyTestnetKey = []string{
		"",
	}

	ETHMainnetRPC = []string{
		"https://eth-mainnet.g.alchemy.com/v2/" + getRandomAlchemyKey(true),
		"https://ethereum-rpc.publicnode.com",
	}

	ETHInnerTxMainnetRPC = []string{
		"https://eth.llamarpc.com",
		"https://eth-pokt.nodies.app",
		"https://eth.merkle.io",
	}

	ETHGoerliRPC = []string{
		"https://ethereum-goerli.publicnode.com",
		"https://endpoints.omniatech.io/v1/eth/goerli/public",
		"https://eth-goerli.g.alchemy.com/v2/" + getRandomAlchemyKey(false),
	}

	ETHSepoliaRPC = []string{
		"https://ethereum-sepolia.publicnode.com",
		"https://eth-sepolia.g.alchemy.com/v2/" + getRandomAlchemyKey(false),
	}

	ETHInnerTxSepoliaRPC = []string{
		// "https://eth-sepolia.g.alchemy.com/v2/" + getRandomInnertxAlchemyKey(false),
	}

	BSCTestnetRPC = []string{
		"https://data-seed-prebsc-1-s1.binance.org:8545",
		"https://data-seed-prebsc-2-s1.binance.org:8545",
		"http://data-seed-prebsc-1-s2.binance.org:8545",
		"http://data-seed-prebsc-2-s2.binance.org:8545",
		"https://data-seed-prebsc-1-s3.binance.org:8545",
		"https://data-seed-prebsc-2-s3.binance.org:8545",
	}

	BSCMainnetRPC = []string{
		"https://bsc-dataseed3.binance.org",
		"https://bsc-dataseed4.binance.org",
		"https://bsc-dataseed1.defibit.io",
		"https://bsc-dataseed2.defibit.io",
		"https://bsc-dataseed3.defibit.io",
		"https://bsc-dataseed4.defibit.io",
		"https://bsc-dataseed1.ninicoin.io",
		"https://bsc-dataseed2.ninicoin.io",
		"https://bsc-dataseed3.ninicoin.io",
		"https://bsc-dataseed4.ninicoin.io",
	}

	OPMainnetRPC = []string{
		// "https://opt-mainnet.g.alchemy.com/v2/" + getRandomAlchemyKey(true),
		// "https://mainnet.optimism.io",
		"https://optimism-rpc.publicnode.com",
		"https://op-pokt.nodies.app",
		"https://1rpc.io/op",
	}

	OPGoerliRPC = []string{
		"https://goerli.optimism.io",
		// "https://opt-goerli.g.alchemy.com/v2/" + getRandomAlchemyKey(false),
	}

	OPSepoliaRPC = []string{
		// "https://opt-sepolia.g.alchemy.com/v2/" + getRandomAlchemyKey(false),
		// "https://sepolia.optimism.io",
		"https://optimism-sepolia.blockpi.network/v1/rpc/public",
		"https://endpoints.omniatech.io/v1/op/sepolia/public",
	}

	ArbitrumOneRPC = []string{
		"https://arb-mainnet.g.alchemy.com/v2/" + getRandomAlchemyKey(true),
		"https://arb1.arbitrum.io/rpc",
	}

	ArbitrumNovaRPC = []string{
		"https://nova.arbitrum.io/rpc",
	}

	ArbitrumGoerliRPC = []string{
		"https://goerli-rollup.arbitrum.io/rpc",
		"https://arb-goerli.g.alchemy.com/v2/" + getRandomAlchemyKey(false),
	}

	ArbitrumSepoliaRPC = []string{
		"https://sepolia-rollup.arbitrum.io/rpc",
		"https://arbitrum-sepolia.blockpi.network/v1/rpc/public",
	}
)

func GetAllRPCUrlByNetwork(id int) []string {
	switch id {
	case ETH_MAINNET:
		return ETHMainnetRPC
	case ETH_GOERLI:
		return ETHGoerliRPC
	case ETH_SEPOLIA:
		return ETHSepoliaRPC
	case OP_MAINNET:
		return OPMainnetRPC
	case OP_GOERLI:
		return OPGoerliRPC
	case OP_SEPOLIA:
		return OPSepoliaRPC
	case BSC_MAINNET:
		return BSCMainnetRPC
	case BSC_TESTNET:
		return BSCTestnetRPC
	case ARBITRUM_ONE:
		return ArbitrumOneRPC
	case ARBITRUM_NOVA:
		return ArbitrumNovaRPC
	case ARBITRUM_GOERLI:
		return ArbitrumGoerliRPC
	case ARBITRUM_SEPOLIA:
		return ArbitrumSepoliaRPC
	}

	return nil
}

func GetRealRpcByArray(rpcs []string) string {
	for _, rpc := range rpcs {
		client.URL = rpc
		var rpcBlockInfo response.RPCBlockInfo
		var jsonRpcRequest request.JsonRpcRequest
		jsonRpcRequest.Id = 1
		jsonRpcRequest.Jsonrpc = "2.0"
		jsonRpcRequest.Method = "eth_getBlockByNumber"
		jsonRpcRequest.Params = []interface{}{"latest", false}
		err := client.HTTPPost(jsonRpcRequest, &rpcBlockInfo)
		if err != nil {
			continue
		}

		height, err := utils.HexStringToInt64(rpcBlockInfo.Result.Number)
		if err != nil || !(height > 0) {
			continue
		}
		return rpc
	}
	return ""
}

// get real rpc url
func GetRPCUrlByNetwork(id int) string {
	switch id {
	case ETH_MAINNET:
		return GetRealRpcByArray(ETHMainnetRPC)
	case ETH_GOERLI:
		return GetRealRpcByArray(ETHGoerliRPC)
	case ETH_SEPOLIA:
		return GetRealRpcByArray(ETHSepoliaRPC)
	case OP_MAINNET:
		return GetRealRpcByArray(OPMainnetRPC)
	case OP_GOERLI:
		return GetRealRpcByArray(OPGoerliRPC)
	case OP_SEPOLIA:
		return GetRealRpcByArray(OPSepoliaRPC)
	case BSC_MAINNET:
		return GetRealRpcByArray(BSCMainnetRPC)
	case BSC_TESTNET:
		return GetRealRpcByArray(BSCTestnetRPC)
	case ARBITRUM_ONE:
		return GetRealRpcByArray(ArbitrumOneRPC)
	case ARBITRUM_NOVA:
		return GetRealRpcByArray(ArbitrumNovaRPC)
	case ARBITRUM_GOERLI:
		return GetRealRpcByArray(ArbitrumGoerliRPC)
	case ARBITRUM_SEPOLIA:
		return GetRealRpcByArray(ArbitrumSepoliaRPC)
	}

	return ""
}

// get real inner tx(trace_debug) rpc url
func GetInnerTxRPCUrlByNetwork(id int) string {
	rand.Seed(time.Now().UnixNano())

	switch id {
	case ETH_MAINNET:
		index := rand.Intn(len(ETHInnerTxMainnetRPC))
		return ETHInnerTxMainnetRPC[index]
	case ETH_SEPOLIA:
		index := rand.Intn(len(ETHInnerTxSepoliaRPC))
		return ETHInnerTxSepoliaRPC[index]
	}

	return ""
}

func getRandomAlchemyKey(isMainnet bool) string {
	rand.Seed(time.Now().UnixMilli())

	if isMainnet {
		index := rand.Intn(len(AlchemyMainnetKey))
		return AlchemyMainnetKey[index]
	} else {
		index := rand.Intn(len(AlchemyTestnetKey))
		return AlchemyTestnetKey[index]
	}
}
