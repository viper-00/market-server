package testnet

import (
	"context"
	"market/global/constant"
	"market/sweep/core"
	"market/sweep/setup"
	MARKET_Client "market/utils/http"
)

var (
	ltcTestnetSweepCount = make(map[int64]int)

	ltcTestnetClient MARKET_Client.Client
)

func SweepLtcTestnetBlockchain() {
	initLtcTestnet()

	go func() {
		for {
			SweepLtcTestnetBlockchainTransaction()
		}
	}()

	go func() {
		for {
			SweepLtcTestnetBlockchainTransactionDetails()
		}
	}()

	go func() {
		for {
			SweepLtcTestnetBlockchainPendingBlock()
		}
	}()
}

func initLtcTestnet() {
	core.SetupLtcLatestBlockHeight(ltcTestnetClient, constant.LTC_TESTNET)

	setup.SetupCacheBlockHeight(context.Background(), constant.LTC_TESTNET)

	setup.SetupSweepBlockHeight(context.Background(), constant.LTC_TESTNET)
}

func SweepLtcTestnetBlockchainTransaction() {
	core.SweepLtcBlockchainTransaction(
		ltcTestnetClient,
		constant.LTC_TESTNET,
		&setup.LtcTestnetPublicKey,
		&ltcTestnetSweepCount,
		&setup.LtcTestnetSweepBlockHeight,
		&setup.LtcTestnetCacheBlockHeight,
		constant.LTC_TESTNET_SWEEP_BLOCK,
		constant.LTC_TESTNET_PENDING_BLOCK,
		constant.LTC_TESTNET_PENDING_TRANSACTION)
}

func SweepLtcTestnetBlockchainTransactionDetails() {
	core.SweepLtcBlockchainTransactionDetails(
		ltcTestnetClient,
		constant.LTC_TESTNET,
		&setup.LtcTestnetPublicKey,
		constant.LTC_TESTNET_PENDING_TRANSACTION)
}

func SweepLtcTestnetBlockchainPendingBlock() {
	core.SweepLtcBlockchainPendingBlock(
		ltcTestnetClient,
		constant.LTC_TESTNET,
		&setup.LtcTestnetPublicKey,
		constant.LTC_TESTNET_PENDING_BLOCK,
		constant.LTC_TESTNET_PENDING_TRANSACTION)
}
