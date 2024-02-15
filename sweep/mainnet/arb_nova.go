package mainnet

import (
	"context"
	"market/global/constant"
	"market/sweep/core"
	"market/sweep/setup"
	MARKET_Client "market/utils/http"
)

var (
	arbitrumNovaSweepCount = make(map[int64]int)

	arbitrumNovaClient MARKET_Client.Client
)

func SweepArbitrumNovaBlockchain() {
	initArbitrumNova()

	go func() {
		for {
			SweepArbitrumNovaBlockchainTransaction()
		}
	}()

	go func() {
		for {
			SweepArbitrumNovaBlockchainTransactionDetails()
		}
	}()

	go func() {
		for {
			SweepArbitrumNovaBlockchainPendingBlock()
		}
	}()
}

func initArbitrumNova() {
	core.SetupLatestBlockHeight(arbitrumNovaClient, constant.ARBITRUM_NOVA)

	setup.SetupCacheBlockHeight(context.Background(), constant.ARBITRUM_NOVA)

	setup.SetupSweepBlockHeight(context.Background(), constant.ARBITRUM_NOVA)
}

func SweepArbitrumNovaBlockchainTransaction() {
	core.SweepBlockchainTransaction(
		arbitrumNovaClient,
		constant.ARBITRUM_NOVA,
		&setup.ArbitrumNovaPublicKey,
		&arbitrumNovaSweepCount,
		&setup.ArbitrumNovaSweepBlockHeight,
		&setup.ArbitrumNovaCacheBlockHeight,
		constant.ARBITRUM_NOVA_SWEEP_BLOCK,
		constant.ARBITRUM_NOVA_PENDING_BLOCK,
		constant.ARBITRUM_NOVA_PENDING_TRANSACTION)
}

func SweepArbitrumNovaBlockchainTransactionDetails() {
	core.SweepBlockchainTransactionDetails(
		arbitrumNovaClient,
		constant.ARBITRUM_NOVA,
		&setup.ArbitrumNovaPublicKey,
		constant.ARBITRUM_NOVA_PENDING_TRANSACTION)
}

func SweepArbitrumNovaBlockchainPendingBlock() {
	core.SweepBlockchainPendingBlock(
		arbitrumNovaClient,
		constant.ARBITRUM_NOVA,
		&setup.ArbitrumNovaPublicKey,
		constant.ARBITRUM_NOVA_PENDING_BLOCK,
		constant.ARBITRUM_NOVA_PENDING_TRANSACTION)
}
