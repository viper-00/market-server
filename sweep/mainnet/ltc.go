package mainnet

import (
	"context"
	"market/global/constant"
	"market/sweep/core"
	"market/sweep/setup"
	MARKET_Client "market/utils/http"
)

var (
	ltcSweepCount = make(map[int64]int)

	ltcClient MARKET_Client.Client
)

func SweepLtcBlockchain() {
	initLtc()

	go func() {
		for {
			SweepLtcBlockchainTransaction()
		}
	}()

	go func() {
		for {
			SweepLtcBlockchainTransactionDetails()
		}
	}()

	go func() {
		for {
			SweepLtcBlockchainPendingBlock()
		}
	}()
}

func initLtc() {
	core.SetupLtcLatestBlockHeight(ltcClient, constant.LTC_MAINNET)

	setup.SetupCacheBlockHeight(context.Background(), constant.LTC_MAINNET)

	setup.SetupSweepBlockHeight(context.Background(), constant.LTC_MAINNET)
}

func SweepLtcBlockchainTransaction() {
	core.SweepLtcBlockchainTransaction(
		ltcClient,
		constant.LTC_MAINNET,
		&setup.LtcPublicKey,
		&ltcSweepCount,
		&setup.LtcSweepBlockHeight,
		&setup.LtcCacheBlockHeight,
		constant.LTC_SWEEP_BLOCK,
		constant.LTC_PENDING_BLOCK,
		constant.LTC_PENDING_TRANSACTION)
}

func SweepLtcBlockchainTransactionDetails() {
	core.SweepLtcBlockchainTransactionDetails(
		ltcClient,
		constant.LTC_MAINNET,
		&setup.LtcPublicKey,
		constant.LTC_PENDING_TRANSACTION)
}

func SweepLtcBlockchainPendingBlock() {
	core.SweepLtcBlockchainPendingBlock(
		ltcClient,
		constant.LTC_MAINNET,
		&setup.LtcPublicKey,
		constant.LTC_PENDING_BLOCK,
		constant.LTC_PENDING_TRANSACTION)
}
