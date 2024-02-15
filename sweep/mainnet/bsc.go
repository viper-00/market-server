package mainnet

import (
	"context"
	"market/global/constant"
	"market/sweep/core"
	"market/sweep/setup"
	MARKET_Client "market/utils/http"
)

var (
	bscSweepCount = make(map[int64]int)

	bscClient MARKET_Client.Client
)

func SweepBscBlockchain() {
	initBsc()

	go func() {
		for {
			SweepBscBlockchainTransaction()
		}
	}()

	go func() {
		for {
			SweepBscBlockchainTransactionDetails()
		}
	}()

	go func() {
		for {
			SweepBscBlockchainPendingBlock()
		}
	}()
}

func initBsc() {
	core.SetupLatestBlockHeight(bscClient, constant.BSC_MAINNET)

	setup.SetupCacheBlockHeight(context.Background(), constant.BSC_MAINNET)

	setup.SetupSweepBlockHeight(context.Background(), constant.BSC_MAINNET)
}

func SweepBscBlockchainTransaction() {
	core.SweepBlockchainTransaction(
		bscClient,
		constant.BSC_MAINNET,
		&setup.BscPublicKey,
		&bscSweepCount,
		&setup.BscSweepBlockHeight,
		&setup.BscCacheBlockHeight,
		constant.BSC_SWEEP_BLOCK,
		constant.BSC_PENDING_BLOCK,
		constant.BSC_PENDING_TRANSACTION)
}

func SweepBscBlockchainTransactionDetails() {
	core.SweepBlockchainTransactionDetails(
		bscClient,
		constant.BSC_MAINNET,
		&setup.BscPublicKey,
		constant.BSC_PENDING_TRANSACTION)
}

func SweepBscBlockchainPendingBlock() {
	core.SweepBlockchainPendingBlock(
		bscClient,
		constant.BSC_MAINNET,
		&setup.BscPublicKey,
		constant.BSC_PENDING_BLOCK,
		constant.BSC_PENDING_TRANSACTION)
}
