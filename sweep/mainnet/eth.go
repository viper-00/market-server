package mainnet

import (
	"context"
	"market/global/constant"
	"market/sweep/core"
	"market/sweep/setup"
	MARKET_Client "market/utils/http"
)

var (
	ethSweepCount = make(map[int64]int)

	ethClient MARKET_Client.Client
)

func SweepEthBlockchain() {

	initEth()

	go func() {
		for {
			SweepEthBlockchainTransaction()
		}
	}()

	go func() {
		for {
			SweepEthBlockchainTransactionDetails()
		}
	}()

	go func() {
		for {
			SweepEthBlockchainPendingBlock()
		}
	}()
}

func initEth() {
	core.SetupLatestBlockHeight(ethClient, constant.ETH_MAINNET)

	setup.SetupCacheBlockHeight(context.Background(), constant.ETH_MAINNET)

	setup.SetupSweepBlockHeight(context.Background(), constant.ETH_MAINNET)
}

func SweepEthBlockchainTransaction() {

	core.SweepBlockchainTransaction(
		ethClient,
		constant.ETH_MAINNET,
		&setup.EthPublicKey,
		&ethSweepCount,
		&setup.EthSweepBlockHeight,
		&setup.EthCacheBlockHeight,
		constant.ETH_SWEEP_BLOCK,
		constant.ETH_PENDING_BLOCK,
		constant.ETH_PENDING_TRANSACTION)
}

func SweepEthBlockchainTransactionDetails() {
	core.SweepBlockchainTransactionDetails(
		ethClient,
		constant.ETH_MAINNET,
		&setup.EthPublicKey,
		constant.ETH_PENDING_TRANSACTION)
}

func SweepEthBlockchainPendingBlock() {
	core.SweepBlockchainPendingBlock(
		ethClient,
		constant.ETH_MAINNET,
		&setup.EthPublicKey,
		constant.ETH_PENDING_BLOCK,
		constant.ETH_PENDING_TRANSACTION)
}
