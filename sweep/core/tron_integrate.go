package core

import (
	"context"
	"errors"
	"fmt"
	"market/global"
	"market/global/constant"
	"market/model/market/request"
	"market/model/market/response"
	"market/sweep/setup"
	sweepUtils "market/sweep/utils"
	"market/sweep/utils/tron"
	"market/utils"
	MARKET_Client "market/utils/http"
	"market/utils/notification"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func SetupTronLatestBlockHeight(client MARKET_Client.Client, chainId int) {
	var err error
	client.URL = constant.TronGetBlockByNetwork(chainId)
	client.Headers = map[string]string{
		"TRON-PRO-API-KEY": constant.GetRandomHTTPKeyByNetwork(chainId),
	}

	var blockRequest request.TronGetBlockRequest
	blockRequest.Detail = false
	var blockResponse response.TronGetBlockResponse
	err = client.HTTPPost(blockRequest, &blockResponse)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	setup.SetupLatestBlockHeight(context.Background(), chainId, int64(blockResponse.BlockHeader.RawData.Number))
}

func SweepTronBlockchainTransaction(
	client MARKET_Client.Client,
	chainId int,
	publicKey *[]string,
	sweepCount *map[int64]int,
	sweepBlockHeight, cacheBlockHeight *int64,
	constantSweepBlock, constantPendingBlock, constantPendingTransaction string) {
	defer utils.HandlePanic()

	if len(*publicKey) <= 0 {
		SetupTronLatestBlockHeight(client, chainId)
		setup.UpdateCacheBlockHeight(context.Background(), chainId)
		setup.UpdateSweepBlockHeight(context.Background(), chainId)
		setup.UpdatePublicKey(context.Background(), chainId)
		time.Sleep(time.Second * 3)
		return
	}

	if *sweepBlockHeight >= *cacheBlockHeight {
		SetupTronLatestBlockHeight(client, chainId)
		setup.UpdateCacheBlockHeight(context.Background(), chainId)
		setup.UpdatePublicKey(context.Background(), chainId)
		time.Sleep(time.Second * 3)
		return
	}

	var err error

	blockN, ok := (*sweepCount)[*sweepBlockHeight]
	if !ok {
		(*sweepCount)[*sweepBlockHeight] = 1
	} else if blockN >= setup.SweepThreshold {
		// skip current block
		_, err = global.MARKET_REDIS.Set(context.Background(), constantSweepBlock, *sweepBlockHeight+1, 0).Result()
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}

		// current block to pending queue
		_, err = global.MARKET_REDIS.RPush(context.Background(), constantPendingBlock, *sweepBlockHeight).Result()
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}

		delete(*sweepCount, *sweepBlockHeight)

		*sweepBlockHeight += 1
		return
	} else {
		(*sweepCount)[*sweepBlockHeight]++
	}

	client.URL = constant.TronGetBlockByNumByNetwork(chainId)
	client.Headers = map[string]string{
		"TRON-PRO-API-KEY": constant.GetRandomHTTPKeyByNetwork(chainId),
	}

	var blockByNumRequest request.TronGetBlockByNumRequest
	blockByNumRequest.Num = int(*sweepBlockHeight)
	var blockByNumResponse response.TronGetBlockByNumResponse
	err = client.HTTPPost(blockByNumRequest, &blockByNumResponse)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	if *sweepBlockHeight == int64(blockByNumResponse.BlockHeader.RawData.Number) {
		if len(blockByNumResponse.Transactions) > 0 {
			for _, transaction := range blockByNumResponse.Transactions {
				if transaction.Ret[0].ContractRet != "SUCCESS" {
					continue
				}

				contractType := transaction.RawData.Contract[0].Type
				var toAddress string

				if contractType == tron.TransferContract {
					toAddress = transaction.RawData.Contract[0].Parameter.Value.ToAddress
				} else if contractType == tron.TriggerSmartContract {
					contractData := transaction.RawData.Contract[0].Parameter.Value.Data
					methodID, _, _ := tron.TronDecodeMethod(contractData)

					method := tron.KnownMethods[methodID]
					if method == "" {
						continue
					}

					toAddress = transaction.RawData.Contract[0].Parameter.Value.ContractAddress
				} else {
					continue
				}

				isMonitorTx := false

			outerCurrentTxLoop:

				for i := 0; i < len(*publicKey); i++ {
					isMonitorTx = tron.IsHandleTransaction(
						chainId,
						transaction.TxID,
						contractType,
						transaction.RawData.Contract[0].Parameter.Value.OwnerAddress,
						toAddress,
						(*publicKey)[i],
						transaction.RawData.Contract[0].Parameter.Value.Data,
					)

					if isMonitorTx {
						// Determine duplicate transactions
						redisTxs, err := global.MARKET_REDIS.LRange(context.Background(), constantPendingTransaction, 0, -1).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return
						}

						for _, redisTx := range redisTxs {
							if redisTx == transaction.TxID {
								continue outerCurrentTxLoop
							}
						}

						_, err = global.MARKET_REDIS.RPush(context.Background(), constantPendingTransaction, transaction.TxID).Result()
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
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), fmt.Sprintf("Not the same height of block: %d - %d", *sweepBlockHeight, int64(blockByNumResponse.BlockHeader.RawData.Number))))
	}
}

func SweepTronBlockchainTransactionDetails(
	client MARKET_Client.Client,
	chainId int,
	publicKey *[]string,
	constantPendingTransaction string,
) {
	defer utils.HandlePanic()

	txHash, err := global.MARKET_REDIS.LIndex(context.Background(), constantPendingTransaction, 0).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			time.Sleep(5 * time.Second)
			return
		}
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	client.URL = constant.TronGetTxByIdByNetwork(chainId)
	client.Headers = map[string]string{
		"TRON-PRO-API-KEY": constant.GetRandomHTTPKeyByNetwork(chainId),
	}

	var txRequest request.TronGetBlockTxByIdRequest
	txRequest.Value = txHash
	var txResponse response.TronGetTxResponse
	err = client.HTTPPost(txRequest, &txResponse)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	if txResponse.Ret[0].ContractRet != "SUCCESS" {
		return
	}

	var notifyRequest request.NotificationRequest
	notifyRequest.Hash = txResponse.TxID
	notifyRequest.Chain = chainId
	notifyRequest.BlockTimestamp = txResponse.RawData.Timestamp

	contractType := txResponse.RawData.Contract[0].Type

	switch contractType {
	case tron.TransferContract:
		handleTransferContractTx(chainId, publicKey, notifyRequest, txResponse)

	case tron.TriggerSmartContract:
		handleTriggerSmartContract(chainId, publicKey, notifyRequest, txResponse)
	}

	_, err = global.MARKET_REDIS.LPop(context.Background(), constantPendingTransaction).Result()
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}
}

func handleTransferContractTx(chainId int, publicKey *[]string, notifyRequest request.NotificationRequest, txResponse response.TronGetTxResponse) {
	var err error

	isSupportContract, contractName, _, decimals := sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb")
	if !isSupportContract {
		return
	}
	if decimals == 0 {
		return
	}

	fromAddress, err := tron.FromHexAddress(txResponse.RawData.Contract[0].Parameter.Value.OwnerAddress)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}
	toAddress, err := tron.FromHexAddress(txResponse.RawData.Contract[0].Parameter.Value.ToAddress)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	notifyRequest.Amount = utils.CalculateBalance(big.NewInt(int64(txResponse.RawData.Contract[0].Parameter.Value.Amount)), decimals)
	notifyRequest.Token = contractName

	handleNotification(chainId, publicKey, notifyRequest, fromAddress, toAddress)
}

func handleTriggerSmartContract(chainId int, publicKey *[]string, notifyRequest request.NotificationRequest, txResponse response.TronGetTxResponse) {
	contractData := txResponse.RawData.Contract[0].Parameter.Value.Data
	methodID, _, _ := tron.TronDecodeMethod(contractData)

	method := tron.KnownMethods[methodID]
	if method == "" {
		return
	}

	switch method {
	case tron.Transfer:
		fromAddress, err := tron.FromHexAddress(txResponse.RawData.Contract[0].Parameter.Value.OwnerAddress)
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			break
		}
		toAddress, err := tron.FromHexAddress("41" + contractData[32:72])
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			break
		}
		amountInt, err := strconv.ParseInt(contractData[len(contractData)-64:], 16, 64)
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}

		contractAddress, err := tron.FromHexAddress(txResponse.RawData.Contract[0].Parameter.Value.ContractAddress)
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}

		_, contractName, _, decimals := sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, contractAddress)

		notifyRequest.Token = contractName
		notifyRequest.Amount = utils.CalculateBalance(big.NewInt(int64(amountInt)), decimals)

		handleNotification(chainId, publicKey, notifyRequest, fromAddress, toAddress)

	case tron.TransferFrom:
		fromAddress, err := tron.FromHexAddress("41" + contractData[32:72])
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		}
		toAddress, err := tron.FromHexAddress("41" + contractData[96:136])
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		}
		amountInt, err := strconv.ParseInt(contractData[len(contractData)-64:], 16, 64)
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}

		contractAddress, err := tron.FromHexAddress(txResponse.RawData.Contract[0].Parameter.Value.ContractAddress)
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}

		_, contractName, _, decimals := sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, contractAddress)

		notifyRequest.Token = contractName
		notifyRequest.Amount = utils.CalculateBalance(big.NewInt(int64(amountInt)), decimals)

		handleNotification(chainId, publicKey, notifyRequest, fromAddress, toAddress)
	}
}

func handleNotification(chainId int, publicKey *[]string, notifyRequest request.NotificationRequest, fromAddress, toAddress string) {
	var err error

	if fromAddress == "" || toAddress == "" {
		return
	}

	notifyRequest.FromAddress = fromAddress
	notifyRequest.ToAddress = toAddress

	for _, v := range *publicKey {
		if strings.EqualFold(v, fromAddress) {
			notifyRequest.TransactType = "send"
			notifyRequest.Address = v

			err = notification.NotificationRequest(notifyRequest)
			if err != nil {
				global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
				// return
			}
		}

		if strings.EqualFold(v, toAddress) {
			notifyRequest.TransactType = "receive"
			notifyRequest.Address = v

			err = notification.NotificationRequest(notifyRequest)
			if err != nil {
				global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
				// return
			}
		}
	}
}

func SweepTronBlockchainPendingBlock(
	client MARKET_Client.Client,
	chainId int,
	publicKey *[]string,
	constantPendingBlock, constantPendingTransaction string,
) {
	defer utils.HandlePanic()

	blockHeight, err := global.MARKET_REDIS.LIndex(context.Background(), constantPendingBlock, 0).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			time.Sleep(10 * time.Second)
			return
		}
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	blockHeightInt, err := strconv.ParseInt(blockHeight, 10, 64)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	client.URL = constant.TronGetBlockByNumByNetwork(chainId)
	client.Headers = map[string]string{
		"TRON-PRO-API-KEY": constant.GetRandomHTTPKeyByNetwork(chainId),
	}

	var blockByNumRequest request.TronGetBlockByNumRequest
	blockByNumRequest.Num = int(blockHeightInt)
	var blockByNumResponse response.TronGetBlockByNumResponse
	err = client.HTTPPost(blockByNumRequest, &blockByNumResponse)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	if blockHeightInt == int64(blockByNumResponse.BlockHeader.RawData.Number) {
		if len(blockByNumResponse.Transactions) > 0 {
			for _, transaction := range blockByNumResponse.Transactions {
				if transaction.Ret[0].ContractRet != "SUCCESS" {
					continue
				}

				contractType := transaction.RawData.Contract[0].Type
				var toAddress string

				if contractType == tron.TransferContract {
					toAddress = transaction.RawData.Contract[0].Parameter.Value.ToAddress
				} else if contractType == tron.TriggerSmartContract {
					contractData := transaction.RawData.Contract[0].Parameter.Value.Data
					methodID, _, _ := tron.TronDecodeMethod(contractData)

					method := tron.KnownMethods[methodID]
					if method == "" {
						continue
					}

					toAddress = transaction.RawData.Contract[0].Parameter.Value.ContractAddress
				} else {
					continue
				}

				isMonitorTx := false

			outerCurrentTxLoop:

				for i := 0; i < len(*publicKey); i++ {
					isMonitorTx = tron.IsHandleTransaction(
						chainId,
						transaction.TxID,
						contractType,
						transaction.RawData.Contract[0].Parameter.Value.OwnerAddress,
						toAddress,
						(*publicKey)[i],
						transaction.RawData.Contract[0].Parameter.Value.Data,
					)

					if isMonitorTx {
						// Determine duplicate transactions
						redisTxs, err := global.MARKET_REDIS.LRange(context.Background(), constantPendingTransaction, 0, -1).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return
						}

						for _, redisTx := range redisTxs {
							if redisTx == transaction.TxID {
								continue outerCurrentTxLoop
							}
						}

						_, err = global.MARKET_REDIS.RPush(context.Background(), constantPendingTransaction, transaction.TxID).Result()
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
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), fmt.Sprintf("Not the same height of block: %d - %d", blockHeightInt, int64(blockByNumResponse.BlockHeader.RawData.Number))))
	}
}
