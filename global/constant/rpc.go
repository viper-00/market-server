package constant

import (
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

	BSCTestnetRPC = []string{
		"https://bsc-testnet-dataseed.bnbchain.org",
		"https://bsc-testnet.bnbchain.org",
		"https://bsc-prebsc-dataseed.bnbchain.org",
	}

	BSCMainnetRPC = []string{
		"https://bsc-dataseed.bnbchain.org",
		"https://bsc-dataseed.nariox.org",
		"https://bsc-dataseed.defibit.io",
		"https://bsc-dataseed.ninicoin.io",
		"https://bsc.nodereal.io",
		"https://bsc-dataseed-public.bnbchain.org",
		"https://bscrpc.com",
	}

	OPMainnetRPC = []string{
		"https://mainnet.optimism.io",
		"https://opt-mainnet.g.alchemy.com/v2/" + getRandomAlchemyKey(true),
	}

	OPGoerliRPC = []string{
		"https://goerli.optimism.io",
		"https://opt-goerli.g.alchemy.com/v2/" + getRandomAlchemyKey(false),
	}

	OPSepoliaRPC = []string{
		"https://sepolia.optimism.io",
		"https://optimism-sepolia.blockpi.network/v1/rpc/public",
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

func GetRPCUrlByNetwork(id int) string {
	rand.Seed(time.Now().UnixNano())

	switch id {
	case ETH_MAINNET:
		index := rand.Intn(len(ETHMainnetRPC))
		return ETHMainnetRPC[index]
	case ETH_GOERLI:
		index := rand.Intn(len(ETHGoerliRPC))
		return ETHGoerliRPC[index]
	case ETH_SEPOLIA:
		index := rand.Intn(len(ETHSepoliaRPC))
		return ETHSepoliaRPC[index]
	case OP_MAINNET:
		index := rand.Intn(len(OPMainnetRPC))
		return OPMainnetRPC[index]
	case OP_GOERLI:
		index := rand.Intn(len(OPGoerliRPC))
		return OPGoerliRPC[index]
	case OP_SEPOLIA:
		index := rand.Intn(len(OPSepoliaRPC))
		return OPSepoliaRPC[index]
	case BSC_MAINNET:
		index := rand.Intn(len(BSCMainnetRPC))
		return BSCMainnetRPC[index]
	case BSC_TESTNET:
		index := rand.Intn(len(BSCTestnetRPC))
		return BSCTestnetRPC[index]
	case ARBITRUM_ONE:
		index := rand.Intn(len(ArbitrumOneRPC))
		return ArbitrumOneRPC[index]
	case ARBITRUM_NOVA:
		index := rand.Intn(len(ArbitrumNovaRPC))
		return ArbitrumNovaRPC[index]
	case ARBITRUM_GOERLI:
		index := rand.Intn(len(ArbitrumGoerliRPC))
		return ArbitrumGoerliRPC[index]
	case ARBITRUM_SEPOLIA:
		index := rand.Intn(len(ArbitrumSepoliaRPC))
		return ArbitrumSepoliaRPC[index]
	}

	return ""
}

func getRandomAlchemyKey(isMainnet bool) string {
	rand.Seed(time.Now().UnixNano())

	if isMainnet {
		index := rand.Intn(len(AlchemyMainnetKey))
		return AlchemyMainnetKey[index]
	} else {
		index := rand.Intn(len(AlchemyTestnetKey))
		return AlchemyTestnetKey[index]
	}
}
