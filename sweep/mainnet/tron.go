package mainnet

import (
	"context"
	"market/global/constant"
	"market/sweep/core"
	"market/sweep/setup"
	MARKET_Client "market/utils/http"
)

var (
	tronSweepCount = make(map[int64]int)

	tronClient MARKET_Client.Client
)

func SweepTronBlockchain() {
	initTron()

	go func() {
		for {
			SweepTronBlockchainTransaction()
		}
	}()

	go func() {
		for {
			SweepTronBlockchainTransactionDetails()
		}
	}()

	go func() {
		for {
			SweepTronBlockchainPendingBlock()
		}
	}()
}

func initTron() {
	core.SetupTronLatestBlockHeight(tronClient, constant.TRON_MAINNET)

	setup.SetupCacheBlockHeight(context.Background(), constant.TRON_MAINNET)

	setup.SetupSweepBlockHeight(context.Background(), constant.TRON_MAINNET)
}

func SweepTronBlockchainTransaction() {
	core.SweepTronBlockchainTransaction(
		tronClient,
		constant.TRON_MAINNET,
		&setup.TronPublicKey,
		&tronSweepCount,
		&setup.TronSweepBlockHeight,
		&setup.TronCacheBlockHeight,
		constant.TRON_SWEEP_BLOCK,
		constant.TRON_PENDING_BLOCK,
		constant.TRON_PENDING_TRANSACTION)
}

func SweepTronBlockchainTransactionDetails() {
	core.SweepTronBlockchainTransactionDetails(
		tronClient,
		constant.TRON_MAINNET,
		&setup.TronPublicKey,
		constant.TRON_PENDING_TRANSACTION)
}

func SweepTronBlockchainPendingBlock() {
	core.SweepTronBlockchainPendingBlock(
		tronClient,
		constant.TRON_MAINNET,
		&setup.TronPublicKey,
		constant.TRON_PENDING_BLOCK,
		constant.TRON_PENDING_TRANSACTION)
}
