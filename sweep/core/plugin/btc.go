package plugin

import (
	"context"
	"fmt"
	"market/global"
	"market/global/constant"
	"market/model/market/request"
	"market/model/market/response/mempool"
	"market/model/market/response/tatum"
	sweepUtils "market/sweep/utils"
	"market/sweep/utils/btc"
	"market/utils"
	MARKET_Client "market/utils/http"
	"market/utils/notification"
	"math/big"
	"strconv"
	"strings"
)

func GetBtcBlockHeightByTatum(client MARKET_Client.Client, chainId int) int64 {
	var err error
	client.URL = constant.TatumGetBitcoinInfo
	client.Headers = map[string]string{
		"x-api-key": constant.GetTatumRandomKeyByNetwork(chainId),
	}

	var bitcoinInfoResponse tatum.TatumGetBitcoinInfo
	err = client.HTTPGet(&bitcoinInfoResponse)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return 0
	}

	return int64(bitcoinInfoResponse.Blocks)
}

func GetBtcBlockHeightByMempool(client MARKET_Client.Client, chainId int) int64 {
	var err error
	client.URL = constant.MempoolGetBlockHeightByNetwork(chainId)
	var bitcoinHeight int64
	err = client.HTTPGetUnique(&bitcoinHeight)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return 0
	}

	return bitcoinHeight
}

func HandleBtcBlockTransactionsByTatum(
	client MARKET_Client.Client,
	chainId int,
	publicKey *[]string,
	sweepCount *map[int64]int,
	sweepBlockHeight *int64,
	constantSweepBlock, constantPendingTransaction string,
) {
	var err error

	client.URL = constant.TatumGetBitcoinBlockByHashOrHeight + fmt.Sprint(*sweepBlockHeight)
	client.Headers = map[string]string{
		"x-api-key": constant.GetTatumRandomKeyByNetwork(chainId),
	}

	var bitcoinBlockResponse tatum.TatumGetBitcoinBlock
	err = client.HTTPGet(&bitcoinBlockResponse)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	if *sweepBlockHeight == int64(bitcoinBlockResponse.Height) {
		if len(bitcoinBlockResponse.Txs) > 0 {
			for _, transaction := range bitcoinBlockResponse.Txs {

				if len(transaction.Inputs) == 0 || len(transaction.Outputs) == 0 {
					continue
				}

			outerCurrentTxLoop:
				for i := 0; i < len(*publicKey); i++ {
					isMonitorTx := false

					if len(transaction.Inputs) > 0 {
						for _, input := range transaction.Inputs {
							if strings.EqualFold((*publicKey)[i], input.Coin.Address) {
								isMonitorTx = true
								break
							}
						}
					}

					if len(transaction.Outputs) > 0 {
						for _, output := range transaction.Outputs {
							if strings.EqualFold((*publicKey)[i], output.Address) {
								isMonitorTx = true
								break
							}
						}
					}

					if isMonitorTx {
						// Determine duplicate transactions
						redisTxs, err := global.MARKET_REDIS.LRange(context.Background(), constantPendingTransaction, 0, -1).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return
						}

						for _, redisTx := range redisTxs {
							if redisTx == transaction.Hash {
								continue outerCurrentTxLoop
							}
						}

						_, err = global.MARKET_REDIS.RPush(context.Background(), constantPendingTransaction, transaction.Hash).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return
						}

						break
					}
				}
			}
		}

		_, err = global.MARKET_REDIS.Set(context.Background(), constantSweepBlock, *sweepBlockHeight+1, 0).Result()
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}

		delete(*sweepCount, *sweepBlockHeight)

		*sweepBlockHeight += 1
	} else {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), fmt.Sprintf("Not the same sweepBlockHeight and blockHeight: %d - %d", *sweepBlockHeight, int64(bitcoinBlockResponse.Height))))
	}
}

func HandleBtcBlockTransactionsByMempool(
	client MARKET_Client.Client,
	chainId int,
	publicKey *[]string,
	sweepCount *map[int64]int,
	sweepBlockHeight *int64,
	constantSweepBlock, constantPendingTransaction string,
) {
	var err error

	var blockHash string
	client.URL = fmt.Sprintf(constant.MempoolGetBlockHashByNetwork(chainId), *sweepBlockHeight)
	err = client.HTTPGetUnique(&blockHash)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	var block mempool.MempoolBlock
	client.URL = fmt.Sprintf(constant.MempoolGetBlockByNetwork(chainId), blockHash)
	err = client.HTTPGet(&block)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	var bitcoinTxsResponses []mempool.MempoolTx

	for i := 0; i < block.TxCount; i += 25 {
		client.URL = fmt.Sprintf(constant.MempoolGetBlockTransactionByNetwork(chainId), blockHash, i)
		var bitcoinTxsResponse []mempool.MempoolTx
		err = client.HTTPGet(&bitcoinTxsResponse)
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}

		bitcoinTxsResponses = append(bitcoinTxsResponses, bitcoinTxsResponse...)
	}

	if len(bitcoinTxsResponses) == 0 {
		return
	}

	if *sweepBlockHeight == int64(block.Height) {
		if len(bitcoinTxsResponses) > 0 {
			for _, transaction := range bitcoinTxsResponses {

				if len(transaction.Vin) == 0 || len(transaction.Vout) == 0 {
					continue
				}

			outerCurrentTxLoop:
				for i := 0; i < len(*publicKey); i++ {
					isMonitorTx := false

					if len(transaction.Vin) > 0 {
						for _, input := range transaction.Vin {
							if strings.EqualFold((*publicKey)[i], input.Prevout.Scriptpubkey_address) {
								isMonitorTx = true
								break
							}
						}
					}

					if len(transaction.Vout) > 0 {
						for _, output := range transaction.Vout {
							if strings.EqualFold((*publicKey)[i], output.Scriptpubkey_address) {
								isMonitorTx = true
								break
							}
						}
					}

					if isMonitorTx {

						// Determine duplicate transactions
						redisTxs, err := global.MARKET_REDIS.LRange(context.Background(), constantPendingTransaction, 0, -1).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return
						}

						for _, redisTx := range redisTxs {
							if redisTx == transaction.TxId {
								continue outerCurrentTxLoop
							}
						}

						_, err = global.MARKET_REDIS.RPush(context.Background(), constantPendingTransaction, transaction.TxId).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return
						}

						break
					}
				}
			}
		}

		_, err = global.MARKET_REDIS.Set(context.Background(), constantSweepBlock, *sweepBlockHeight+1, 0).Result()
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}

		delete(*sweepCount, *sweepBlockHeight)

		*sweepBlockHeight += 1
	} else {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), fmt.Sprintf("Not the same sweepBlockHeight and blockHeight: %d - %d", *sweepBlockHeight, int64(block.Height))))
	}
}

func HandleBtcTransactionDetailsByTatum(
	client MARKET_Client.Client,
	chainId int,
	publicKey *[]string,
	constantPendingTransaction string,
	txHash string,
) {
	var err error

	client.URL = constant.TatumGetBitcoinTxByHash + txHash
	client.Headers = map[string]string{
		"x-api-key": constant.GetTatumRandomKeyByNetwork(chainId),
	}

	var bitcoinTxResponse tatum.TatumBitcoinTx
	err = client.HTTPGet(&bitcoinTxResponse)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	var notifyRequest request.NotificationRequest

	notifyRequest.Hash = bitcoinTxResponse.Hash
	notifyRequest.Chain = chainId
	notifyRequest.BlockTimestamp = bitcoinTxResponse.Time * 1000

	if len(bitcoinTxResponse.Inputs) == 0 || len(bitcoinTxResponse.Outputs) == 0 {
		return
	}

	_, _, _, decimals := sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, "")
	if decimals == 0 {
		return
	}

	for _, input := range bitcoinTxResponse.Inputs {
		if input.Coin.Address != "" {
			notifyRequest.FromAddress = input.Coin.Address
			continue
		}
	}

	var isOnmiUSDT bool
	var onmiData map[string]int

	for _, output := range bitcoinTxResponse.Outputs {
		if output.Value == 0 && output.Address == "" {
			onmiData, isOnmiUSDT = btc.ParseOmniUSDTData(output.Script)
			break
		} else {
			isOnmiUSDT = false
		}
	}

	if isOnmiUSDT {
		// omni
		notifyRequest.Token = "USDT"
		notifyRequest.Amount = strconv.Itoa(onmiData["token_amount"])

		for _, omniOutput := range bitcoinTxResponse.Outputs {
			if strings.EqualFold(omniOutput.Address, notifyRequest.FromAddress) || omniOutput.Value == 0 || omniOutput.Address == "" {
				continue
			}
			notifyRequest.ToAddress = omniOutput.Address
		}

		for _, v := range *publicKey {
			notifyRequest.Address = v

			if strings.EqualFold(notifyRequest.FromAddress, v) {
				notifyRequest.TransactType = "send"

				err = notification.NotificationRequest(notifyRequest)
				if err != nil {
					global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
					// return
				}
			}

			if strings.EqualFold(notifyRequest.ToAddress, v) {
				notifyRequest.TransactType = "receive"

				err = notification.NotificationRequest(notifyRequest)
				if err != nil {
					global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
					// return
				}
			}
		}

	} else {
		// btc
		notifyRequest.Token = "BTC"
		for _, output := range bitcoinTxResponse.Outputs {
			if strings.EqualFold(output.Address, notifyRequest.FromAddress) {
				continue
			}

			notifyRequest.Amount = utils.CalculateBalance(big.NewInt(int64(output.Value)), decimals)
			for _, v := range *publicKey {
				notifyRequest.Address = v
				notifyRequest.ToAddress = output.Address

				if strings.EqualFold(notifyRequest.FromAddress, v) {
					notifyRequest.TransactType = "send"

					err = notification.NotificationRequest(notifyRequest)
					if err != nil {
						global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
						// return
					}
				}

				if strings.EqualFold(output.Address, v) {
					notifyRequest.TransactType = "receive"

					err = notification.NotificationRequest(notifyRequest)
					if err != nil {
						global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
						// return
					}
				}
			}
		}
	}

	_, err = global.MARKET_REDIS.LPop(context.Background(), constantPendingTransaction).Result()
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}
}

func HandleBtcTransactionDetailsByMempool(
	client MARKET_Client.Client,
	chainId int,
	publicKey *[]string,
	constantPendingTransaction string,
	txHash string,
) {
	var err error

	client.URL = fmt.Sprintf(constant.MempoolGetTransctionByNetwork(chainId), txHash)

	var bitcoinTxResponse mempool.MempoolTx
	err = client.HTTPGet(&bitcoinTxResponse)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	var notifyRequest request.NotificationRequest

	notifyRequest.Hash = bitcoinTxResponse.TxId
	notifyRequest.Chain = chainId
	notifyRequest.BlockTimestamp = bitcoinTxResponse.Status.BlockTime * 1000

	if len(bitcoinTxResponse.Vin) == 0 || len(bitcoinTxResponse.Vout) == 0 {
		return
	}

	_, _, _, decimals := sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, "")
	if decimals == 0 {
		return
	}

	for _, input := range bitcoinTxResponse.Vin {
		if input.Prevout.Scriptpubkey_address != "" {
			notifyRequest.FromAddress = input.Prevout.Scriptpubkey_address
			continue
		}
	}

	var isOnmiUSDT bool
	var onmiData map[string]int

	for _, output := range bitcoinTxResponse.Vout {
		if output.Value == 0 && output.Scriptpubkey_address == "" {
			onmiData, isOnmiUSDT = btc.ParseOmniUSDTData(output.Scriptpubkey)
			break
		} else {
			isOnmiUSDT = false
		}
	}

	if isOnmiUSDT {
		// omni
		notifyRequest.Token = "USDT"
		notifyRequest.Amount = strconv.Itoa(onmiData["token_amount"])

		for _, omniOutput := range bitcoinTxResponse.Vout {
			if strings.EqualFold(omniOutput.Scriptpubkey_address, notifyRequest.FromAddress) || omniOutput.Value == 0 || omniOutput.Scriptpubkey_address == "" {
				continue
			}
			notifyRequest.ToAddress = omniOutput.Scriptpubkey_address
		}

		for _, v := range *publicKey {
			notifyRequest.Address = v

			if strings.EqualFold(notifyRequest.FromAddress, v) {
				notifyRequest.TransactType = "send"

				err = notification.NotificationRequest(notifyRequest)
				if err != nil {
					global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
					// return
				}
			}

			if strings.EqualFold(notifyRequest.ToAddress, v) {
				notifyRequest.TransactType = "receive"

				err = notification.NotificationRequest(notifyRequest)
				if err != nil {
					global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
					// return
				}
			}
		}

	} else {
		// btc
		notifyRequest.Token = "BTC"
		for _, output := range bitcoinTxResponse.Vout {
			if strings.EqualFold(output.Scriptpubkey_address, notifyRequest.FromAddress) {
				continue
			}

			notifyRequest.Amount = utils.CalculateBalance(big.NewInt(int64(output.Value)), decimals)
			for _, v := range *publicKey {
				notifyRequest.Address = v
				notifyRequest.ToAddress = output.Scriptpubkey_address

				if strings.EqualFold(notifyRequest.FromAddress, v) {
					notifyRequest.TransactType = "send"

					err = notification.NotificationRequest(notifyRequest)
					if err != nil {
						global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
						// return
					}
				}

				if strings.EqualFold(output.Scriptpubkey_address, v) {
					notifyRequest.TransactType = "receive"

					err = notification.NotificationRequest(notifyRequest)
					if err != nil {
						global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
						// return
					}
				}
			}
		}
	}

	_, err = global.MARKET_REDIS.LPop(context.Background(), constantPendingTransaction).Result()
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}
}

func HandleBtcPendingBlockByTatum(
	client MARKET_Client.Client,
	chainId int,
	publicKey *[]string,
	constantPendingBlock, constantPendingTransaction string,
	blockHeight string,
	blockHeightInt int64,
) {
	var err error

	client.URL = constant.TatumGetBitcoinBlockByHashOrHeight + blockHeight
	client.Headers = map[string]string{
		"x-api-key": constant.GetTatumRandomKeyByNetwork(chainId),
	}

	var bitcoinBlockResponse tatum.TatumGetBitcoinBlock
	err = client.HTTPGet(&bitcoinBlockResponse)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	if int(blockHeightInt) == bitcoinBlockResponse.Height {
		if len(bitcoinBlockResponse.Txs) > 0 {
			for _, transaction := range bitcoinBlockResponse.Txs {

				if len(transaction.Inputs) == 0 || len(transaction.Outputs) == 0 {
					continue
				}

			outerCurrentTxLoop:
				for i := 0; i < len(*publicKey); i++ {
					isMonitorTx := false

					if len(transaction.Inputs) > 0 {
						for _, input := range transaction.Inputs {
							if strings.EqualFold((*publicKey)[i], input.Coin.Address) {
								isMonitorTx = true
								break
							}
						}
					}

					if len(transaction.Outputs) > 0 {
						for _, output := range transaction.Outputs {
							if strings.EqualFold((*publicKey)[i], output.Address) {
								isMonitorTx = true
								break
							}
						}
					}

					if isMonitorTx {

						// Determine duplicate transactions
						redisTxs, err := global.MARKET_REDIS.LRange(context.Background(), constantPendingTransaction, 0, -1).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return
						}

						for _, redisTx := range redisTxs {
							if redisTx == transaction.Hash {
								continue outerCurrentTxLoop
							}
						}

						_, err = global.MARKET_REDIS.RPush(context.Background(), constantPendingTransaction, transaction.Hash).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return
						}

						break
					}
				}
			}
		}

		_, err = global.MARKET_REDIS.LPop(context.Background(), constantPendingBlock).Result()
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		}

	} else {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), fmt.Sprintf("Not the same height of block: %d - %d", blockHeightInt, int64(bitcoinBlockResponse.Height))))
	}
}

func HandleBtcPendingBlockByMempool(
	client MARKET_Client.Client,
	chainId int,
	publicKey *[]string,
	constantPendingBlock, constantPendingTransaction string,
	blockHeight string,
	blockHeightInt int64,
) {
	var err error

	var blockHash string
	client.URL = fmt.Sprintf(constant.MempoolGetBlockHashByNetwork(chainId), blockHeightInt)
	err = client.HTTPGetUnique(&blockHash)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	var block mempool.MempoolBlock
	client.URL = fmt.Sprintf(constant.MempoolGetBlockByNetwork(chainId), blockHash)
	err = client.HTTPGet(&block)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	var bitcoinTxsResponses []mempool.MempoolTx

	for i := 0; i < block.TxCount; i += 25 {
		client.URL = fmt.Sprintf(constant.MempoolGetBlockTransactionByNetwork(chainId), blockHash, i)
		var bitcoinTxsResponse []mempool.MempoolTx
		err = client.HTTPGet(&bitcoinTxsResponse)
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}

		bitcoinTxsResponses = append(bitcoinTxsResponses, bitcoinTxsResponse...)
	}

	if len(bitcoinTxsResponses) == 0 {
		return
	}

	if blockHeightInt == int64(block.Height) {
		if len(bitcoinTxsResponses) > 0 {
			for _, transaction := range bitcoinTxsResponses {

				if len(transaction.Vin) == 0 || len(transaction.Vout) == 0 {
					continue
				}

			outerCurrentTxLoop:
				for i := 0; i < len(*publicKey); i++ {
					isMonitorTx := false

					if len(transaction.Vin) > 0 {
						for _, input := range transaction.Vin {
							if strings.EqualFold((*publicKey)[i], input.Prevout.Scriptpubkey_address) {
								isMonitorTx = true
								break
							}
						}
					}

					if len(transaction.Vout) > 0 {
						for _, output := range transaction.Vout {
							if strings.EqualFold((*publicKey)[i], output.Scriptpubkey_address) {
								isMonitorTx = true
								break
							}
						}
					}

					if isMonitorTx {

						// Determine duplicate transactions
						redisTxs, err := global.MARKET_REDIS.LRange(context.Background(), constantPendingTransaction, 0, -1).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return
						}

						for _, redisTx := range redisTxs {
							if redisTx == transaction.TxId {
								continue outerCurrentTxLoop
							}
						}

						_, err = global.MARKET_REDIS.RPush(context.Background(), constantPendingTransaction, transaction.TxId).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return
						}

						break
					}
				}
			}
		}

		_, err = global.MARKET_REDIS.LPop(context.Background(), constantPendingBlock).Result()
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		}

	} else {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), fmt.Sprintf("Not the same sweepBlockHeight and blockHeight: %d - %d", blockHeightInt, int64(block.Height))))
	}
}
