package testnet

import (
	"context"
	"market/global/constant"
	"market/sweep/core"
	"market/sweep/setup"
	MARKET_Client "market/utils/http"
)

var (
	bscTestnetSweepCount = make(map[int64]int)

	bscTestnetClient MARKET_Client.Client
)

func SweepBscTestnetBlockchain() {
	initBscTestnet()

	go func() {
		for {
			SweepBscTestnetBlockchainTransaction()
		}
	}()

	go func() {
		for {
			SweepBscTestnetBlockchainTransactionDetails()
		}
	}()

	go func() {
		for {
			SweepBscTestnetBlockchainPendingBlock()
		}
	}()
}

func initBscTestnet() {
	core.SetupLatestBlockHeight(bscTestnetClient, constant.BSC_TESTNET)

	setup.SetupCacheBlockHeight(context.Background(), constant.BSC_TESTNET)

	setup.SetupSweepBlockHeight(context.Background(), constant.BSC_TESTNET)
}

func SweepBscTestnetBlockchainTransaction() {
	core.SweepBlockchainTransaction(
		bscTestnetClient,
		constant.BSC_TESTNET,
		&setup.BscTestnetPublicKey,
		&bscTestnetSweepCount,
		&setup.BscTestnetSweepBlockHeight,
		&setup.BscTestnetCacheBlockHeight,
		constant.BSC_TESTNET_SWEEP_BLOCK,
		constant.BSC_TESTNET_PENDING_BLOCK,
		constant.BSC_TESTNET_PENDING_TRANSACTION)
}

func SweepBscTestnetBlockchainTransactionDetails() {
	core.SweepBlockchainTransactionDetails(
		bscTestnetClient,
		constant.BSC_TESTNET,
		&setup.BscTestnetPublicKey,
		constant.BSC_TESTNET_PENDING_TRANSACTION)
}

func SweepBscTestnetBlockchainPendingBlock() {
	core.SweepBlockchainPendingBlock(
		bscTestnetClient,
		constant.BSC_TESTNET,
		&setup.BscTestnetPublicKey,
		constant.BSC_TESTNET_PENDING_BLOCK,
		constant.BSC_TESTNET_PENDING_TRANSACTION)
}
