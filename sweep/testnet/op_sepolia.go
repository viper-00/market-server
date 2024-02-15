package testnet

import (
	"context"
	"market/global/constant"
	"market/sweep/core"
	"market/sweep/setup"
	MARKET_Client "market/utils/http"
)

var (
	opSepoliaSweepCount = make(map[int64]int)

	opSepoliaClient MARKET_Client.Client
)

func SweepOpSepoliaBlockchain() {

	initOpSepolia()

	go func() {
		for {
			SweepOpSepoliaBlockchainTransaction()
		}
	}()

	go func() {
		for {
			SweepOpSepoliaBlockchainTransactionDetails()
		}
	}()

	go func() {
		for {
			SweepOpSepoliaBlockchainPendingBlock()
		}
	}()
}

func initOpSepolia() {
	core.SetupLatestBlockHeight(opSepoliaClient, constant.OP_SEPOLIA)

	setup.SetupCacheBlockHeight(context.Background(), constant.OP_SEPOLIA)

	setup.SetupSweepBlockHeight(context.Background(), constant.OP_SEPOLIA)
}

func SweepOpSepoliaBlockchainTransaction() {
	core.SweepBlockchainTransaction(
		opSepoliaClient,
		constant.OP_SEPOLIA,
		&setup.OpSepoliaPublicKey,
		&opSepoliaSweepCount,
		&setup.OpSepoliaSweepBlockHeight,
		&setup.OpSepoliaCacheBlockHeight,
		constant.OP_SEPOLIA_SWEEP_BLOCK,
		constant.OP_SEPOLIA_PENDING_BLOCK,
		constant.OP_SEPOLIA_PENDING_TRANSACTION)
}

func SweepOpSepoliaBlockchainTransactionDetails() {
	core.SweepBlockchainTransactionDetails(
		opSepoliaClient,
		constant.OP_SEPOLIA,
		&setup.OpSepoliaPublicKey,
		constant.OP_SEPOLIA_PENDING_TRANSACTION)
}

func SweepOpSepoliaBlockchainPendingBlock() {
	core.SweepBlockchainPendingBlock(
		opSepoliaClient,
		constant.OP_SEPOLIA,
		&setup.OpSepoliaPublicKey,
		constant.OP_SEPOLIA_PENDING_BLOCK,
		constant.OP_SEPOLIA_PENDING_TRANSACTION)
}
