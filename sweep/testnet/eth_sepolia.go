package testnet

import (
	"context"
	"market/global/constant"
	"market/sweep/core"
	"market/sweep/setup"
	MARKET_Client "market/utils/http"
)

var (
	ethSepoliaSweepCount = make(map[int64]int)

	ethSepoliaClient MARKET_Client.Client
)

func SweepEthSepoliaBlockchain() {
	initEthSepolia()

	go func() {
		for {
			SweepEthSepoliaBlockchainTransaction()
		}
	}()

	go func() {
		for {
			SweepEthSepoliaBlockchainTransactionDetails()
		}
	}()

	go func() {
		for {
			SweepEthSepoliaBlockchainPendingBlock()
		}
	}()
}

func initEthSepolia() {
	core.SetupLatestBlockHeight(ethSepoliaClient, constant.ETH_SEPOLIA)

	setup.SetupCacheBlockHeight(context.Background(), constant.ETH_SEPOLIA)

	setup.SetupSweepBlockHeight(context.Background(), constant.ETH_SEPOLIA)
}

func SweepEthSepoliaBlockchainTransaction() {
	core.SweepBlockchainTransaction(
		ethSepoliaClient,
		constant.ETH_SEPOLIA,
		&setup.EthSepoliaPublicKey,
		&ethSepoliaSweepCount,
		&setup.EthSepoliaSweepBlockHeight,
		&setup.EthSepoliaCacheBlockHeight,
		constant.ETH_SEPOLIA_SWEEP_BLOCK,
		constant.ETH_SEPOLIA_PENDING_BLOCK,
		constant.ETH_SEPOLIA_PENDING_TRANSACTION)
}

func SweepEthSepoliaBlockchainTransactionDetails() {
	core.SweepBlockchainTransactionDetails(
		ethSepoliaClient,
		constant.ETH_SEPOLIA,
		&setup.EthSepoliaPublicKey,
		constant.ETH_SEPOLIA_PENDING_TRANSACTION)
}

func SweepEthSepoliaBlockchainPendingBlock() {

	core.SweepBlockchainPendingBlock(
		ethSepoliaClient,
		constant.ETH_SEPOLIA,
		&setup.EthSepoliaPublicKey,
		constant.ETH_SEPOLIA_PENDING_BLOCK,
		constant.ETH_SEPOLIA_PENDING_TRANSACTION)
}
