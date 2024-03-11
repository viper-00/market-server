package sweep

import (
	"context"
	"market/global"
	"market/sweep/mainnet"
	"market/sweep/setup"
	"market/sweep/testnet"
)

func RunBlockSweep() {
	setup.SetupPublicKey(context.Background())

	if global.MARKET_CONFIG.Blockchain.Ethereum {
		if global.MARKET_CONFIG.Blockchain.SweepMainnet {
			mainnet.SweepEthBlockchain()
		} else {
			testnet.SweepEthGoerliBlockchain()
			// testnet.SweepEthSepoliaBlockchain()
		}
	}

	if global.MARKET_CONFIG.Blockchain.Bsc {
		if global.MARKET_CONFIG.Blockchain.SweepMainnet {
			mainnet.SweepBscBlockchain()
		} else {
			testnet.SweepBscTestnetBlockchain()
		}

	}

	if global.MARKET_CONFIG.Blockchain.Bitcoin {
		if global.MARKET_CONFIG.Blockchain.SweepMainnet {
			mainnet.SweepBtcBlockchain()
		} else {
			testnet.SweepBtcTestnetBlockchain()
		}
	}

	if global.MARKET_CONFIG.Blockchain.Tron {
		if global.MARKET_CONFIG.Blockchain.SweepMainnet {
			mainnet.SweepTronBlockchain()
		} else {
			testnet.SweepTronNileBlockchain()
		}
	}

	if global.MARKET_CONFIG.Blockchain.Litecoin {
		if global.MARKET_CONFIG.Blockchain.SweepMainnet {
			mainnet.SweepLtcBlockchain()
		} else {
			testnet.SweepLtcTestnetBlockchain()
		}
	}

	if global.MARKET_CONFIG.Blockchain.Op {
		if global.MARKET_CONFIG.Blockchain.SweepMainnet {
			mainnet.SweepOpBlockchain()
		} else {
			testnet.SweepOpSepoliaBlockchain()
		}
	}

	if global.MARKET_CONFIG.Blockchain.ArbitrumOne {
		if global.MARKET_CONFIG.Blockchain.SweepMainnet {
			mainnet.SweepArbitrumOneBlockchain()
		} else {
			// testnet.SweepArbitrumGoerliBlockchain()
			testnet.SweepArbitrumSepoliaBlockchain()
		}
	}

	if global.MARKET_CONFIG.Blockchain.ArbitrumNova {
		if global.MARKET_CONFIG.Blockchain.SweepMainnet {
			mainnet.SweepArbitrumNovaBlockchain()
		} else {
		}
	}
}
