package mainnet

import (
	"context"
	"market/global/constant"
	"market/sweep/core"
	"market/sweep/setup"
	MARKET_Client "market/utils/http"
)

var (
	opSweepCount = make(map[int64]int)

	opClient MARKET_Client.Client
)

func SweepOpBlockchain() {

	initOp()

	go func() {
		for {
			SweepOpBlockchainTransaction()
		}
	}()

	go func() {
		for {
			SweepOpBlockchainTransactionDetails()
		}
	}()

	go func() {
		for {
			SweepOpBlockchainPendingBlock()
		}
	}()
}

func initOp() {
	core.SetupLatestBlockHeight(opClient, constant.OP_MAINNET)

	setup.SetupCacheBlockHeight(context.Background(), constant.OP_MAINNET)

	setup.SetupSweepBlockHeight(context.Background(), constant.OP_MAINNET)
}

func SweepOpBlockchainTransaction() {
	core.SweepBlockchainTransaction(
		opClient,
		constant.OP_MAINNET,
		&setup.OpPublicKey,
		&opSweepCount,
		&setup.OpSweepBlockHeight,
		&setup.OpCacheBlockHeight,
		constant.OP_SWEEP_BLOCK,
		constant.OP_PENDING_BLOCK,
		constant.OP_PENDING_TRANSACTION)
}

func SweepOpBlockchainTransactionDetails() {
	core.SweepBlockchainTransactionDetails(
		opClient,
		constant.OP_MAINNET,
		&setup.OpPublicKey,
		constant.OP_PENDING_TRANSACTION)
}

func SweepOpBlockchainPendingBlock() {
	core.SweepBlockchainPendingBlock(
		opClient,
		constant.OP_MAINNET,
		&setup.OpPublicKey,
		constant.OP_PENDING_BLOCK,
		constant.OP_PENDING_TRANSACTION)
}
