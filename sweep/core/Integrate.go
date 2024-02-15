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
	"market/sweep/utils/erc20"
	"market/utils"
	"math/big"
	"strconv"
	"time"

	sweepUtils "market/sweep/utils"
	MARKET_Client "market/utils/http"
	"market/utils/notification"

	"github.com/redis/go-redis/v9"
)

func SetupLatestBlockHeight(client MARKET_Client.Client, chainId int) {
	var err error
	client.URL = constant.GetRPCUrlByNetwork(chainId)
	var rpcBlockInfo response.RPCBlockInfo
	var jsonRpcRequest request.JsonRpcRequest
	jsonRpcRequest.Id = 1
	jsonRpcRequest.Jsonrpc = "2.0"
	jsonRpcRequest.Method = "eth_getBlockByNumber"
	jsonRpcRequest.Params = []interface{}{"latest", false}

	err = client.HTTPPost(jsonRpcRequest, &rpcBlockInfo)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	height, err := utils.HexStringToInt64(rpcBlockInfo.Result.Number)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	if height > 0 {
		setup.SetupLatestBlockHeight(context.Background(), chainId, height)
	}
}

func SweepBlockchainTransaction(
	client MARKET_Client.Client,
	chainId int,
	publicKey *[]string,
	sweepCount *map[int64]int,
	sweepBlockHeight, cacheBlockHeight *int64,
	constantSweepBlock, constantPendingBlock, constantPendingTransaction string) {
	defer utils.HandlePanic()

	if len(*publicKey) <= 0 {
		SetupLatestBlockHeight(client, chainId)
		setup.UpdateCacheBlockHeight(context.Background(), chainId)
		setup.UpdateSweepBlockHeight(context.Background(), chainId)
		setup.UpdatePublicKey(context.Background(), chainId)
		time.Sleep(time.Second * 3)
		return
	}

	if *sweepBlockHeight >= *cacheBlockHeight {
		SetupLatestBlockHeight(client, chainId)
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

	client.URL = constant.GetRPCUrlByNetwork(chainId)
	var rpcBlockDetail response.RPCBlockDetail
	var jsonRpcRequest request.JsonRpcRequest
	jsonRpcRequest.Id = 1
	jsonRpcRequest.Jsonrpc = "2.0"
	jsonRpcRequest.Method = "eth_getBlockByNumber"
	jsonRpcRequest.Params = []interface{}{"0x" + strconv.FormatInt(*sweepBlockHeight, 16), true}

	err = client.HTTPPost(jsonRpcRequest, &rpcBlockDetail)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	height, err := utils.HexStringToInt64(rpcBlockDetail.Result.Number)
	if err != nil {
		return
	}

	if *sweepBlockHeight == height {
		if len(rpcBlockDetail.Result.Transactions) > 0 {
			for _, transaction := range rpcBlockDetail.Result.Transactions {

				isMonitorTx := false

			outerCurrentTxLoop:
				for i := 0; i < len(*publicKey); i++ {
					if transaction.Input == "0x" {
						if utils.HexToAddress(transaction.From) == utils.HexToAddress((*publicKey)[i]) || utils.HexToAddress(transaction.To) == utils.HexToAddress((*publicKey)[i]) {
							isMonitorTx = true
						}
					} else {
						// Ordinary erc20 transactions
						isSupportContract, contractName, _, _ := sweepUtils.GetContractInfo(chainId, transaction.To)
						if isSupportContract {
							if erc20.IsHandleTokenTransaction(chainId, transaction.Hash, contractName, transaction.From, transaction.To, (*publicKey)[i], transaction.Input) {
								isMonitorTx = true
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
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), fmt.Sprintf("Not the same height of block: %d - %d", *sweepBlockHeight, height)))
	}
}

func SweepBlockchainTransactionDetails(
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

	client.URL = constant.GetRPCUrlByNetwork(chainId)
	var rpcDetail response.RPCTransactionDetail
	var jsonRpcRequest request.JsonRpcRequest
	jsonRpcRequest.Id = 1
	jsonRpcRequest.Jsonrpc = "2.0"
	jsonRpcRequest.Method = "eth_getTransactionByHash"
	jsonRpcRequest.Params = []interface{}{txHash}

	err = client.HTTPPost(jsonRpcRequest, &rpcDetail)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	var rpcBlockInfo response.RPCBlockInfo
	jsonRpcRequest.Id = 1
	jsonRpcRequest.Jsonrpc = "2.0"
	jsonRpcRequest.Method = "eth_getBlockByNumber"
	jsonRpcRequest.Params = []interface{}{rpcDetail.Result.BlockNumber, false}

	err = client.HTTPPost(jsonRpcRequest, &rpcBlockInfo)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	blockTimeStamp, err := utils.HexStringToInt64(rpcBlockInfo.Result.Timestamp)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	var notifyRequest request.NotificationRequest

	notifyRequest.Hash = rpcDetail.Result.Hash
	notifyRequest.Chain = chainId
	notifyRequest.BlockTimestamp = int(blockTimeStamp) * 1000

	if rpcDetail.Result.Input == "0x" {
		handleERC20(chainId, publicKey, notifyRequest, rpcDetail)
	} else {
		_, contractName, _, _ := sweepUtils.GetContractInfo(chainId, rpcDetail.Result.To)

		switch contractName {
		case constant.ALLINONE:
			handleAllInOne(chainId, publicKey, notifyRequest, rpcDetail)
		case constant.SWAP:
			break
		default:
			handleERC20(chainId, publicKey, notifyRequest, rpcDetail)
		}
	}

	_, err = global.MARKET_REDIS.LPop(context.Background(), constantPendingTransaction).Result()
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}
}

func handleAllInOne(chainId int, publicKey *[]string, notifyRequest request.NotificationRequest, rpcDetail response.RPCTransactionDetail) (err error) {
	_, fromAddresses, toAddresses, tokens, amounts, err := erc20.DecodeAllInOneTransactionInputData(chainId, rpcDetail.Result.Hash, rpcDetail.Result.From, rpcDetail.Result.To, rpcDetail.Result.Input)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	var totalLen = len(fromAddresses)

	for i := 0; i < totalLen; i++ {
		for _, v := range *publicKey {
			if utils.HexToAddress(fromAddresses[i].String()) == utils.HexToAddress(v) {
				isSupportContract, contractName, _, decimals := sweepUtils.GetContractInfo(chainId, tokens[i].String())
				if !isSupportContract {
					continue
				}
				notifyRequest.Token = contractName
				notifyRequest.Amount = utils.CalculateBalance(amounts[i], decimals)
				notifyRequest.TransactType = "send"
				notifyRequest.Address = v
				notifyRequest.FromAddress = v
				notifyRequest.ToAddress = toAddresses[i].String()

				err = notification.NotificationRequest(notifyRequest)
				if err != nil {
					global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
					// return
				}
			}

			if utils.HexToAddress(toAddresses[i].String()) == utils.HexToAddress(v) {
				isSupportContract, contractName, _, decimals := sweepUtils.GetContractInfo(chainId, tokens[i].String())
				if !isSupportContract {
					continue
				}
				notifyRequest.Token = contractName
				notifyRequest.Amount = utils.CalculateBalance(amounts[i], decimals)
				notifyRequest.TransactType = "receive"
				notifyRequest.Address = v
				notifyRequest.FromAddress = fromAddresses[i].String()
				notifyRequest.ToAddress = v

				err = notification.NotificationRequest(notifyRequest)
				if err != nil {
					global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
					// return
				}
			}

		}
	}
	return nil
}

func handleERC20(chainId int, publicKey *[]string, notifyRequest request.NotificationRequest, rpcDetail response.RPCTransactionDetail) {

	var (
		isSupportContract bool
		contractName      string
		decimals          int
	)

	if rpcDetail.Result.Input == "0x" {
		isSupportContract, contractName, _, decimals = sweepUtils.GetContractInfo(chainId, "0x0000000000000000000000000000000000000000")
	} else {
		isSupportContract, contractName, _, decimals = sweepUtils.GetContractInfo(chainId, rpcDetail.Result.To)
	}

	if !isSupportContract {
		return
	}

	fromAddress := rpcDetail.Result.From
	notifyRequest.FromAddress = fromAddress

	if decimals == 0 {
		return
	}

	if !(rpcDetail.Result.Input == "0x") {
		methodName, decodeFromAddress, decodeToAddress, transactionValue, err := erc20.DecodeERC20TransactionInputData(chainId, rpcDetail.Result.Hash, rpcDetail.Result.Input)
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}

		switch methodName {
		case erc20.TransferFrom:
			fromAddress = decodeFromAddress
		}

		notifyRequest.ToAddress = decodeToAddress
		notifyRequest.Token = contractName
		notifyRequest.Amount = utils.CalculateBalance(transactionValue, decimals)

		for _, v := range *publicKey {
			if utils.HexToAddress(fromAddress) == utils.HexToAddress(v) {
				notifyRequest.TransactType = "send"
				notifyRequest.Address = v

				err = notification.NotificationRequest(notifyRequest)
				if err != nil {
					global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
					// return
				}
			}

			if utils.HexToAddress(decodeToAddress) == utils.HexToAddress(v) {
				notifyRequest.TransactType = "receive"
				notifyRequest.Address = v

				err = notification.NotificationRequest(notifyRequest)
				if err != nil {
					global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
					// return
				}
			}
		}

	} else {
		toAddress := rpcDetail.Result.To
		notifyRequest.ToAddress = toAddress

		value, err := utils.HexStringToInt64(rpcDetail.Result.Value)
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}
		notifyRequest.Amount = utils.CalculateBalance(big.NewInt(value), decimals)
		notifyRequest.Token = contractName

		for _, v := range *publicKey {
			if utils.HexToAddress(fromAddress) == utils.HexToAddress(v) {
				notifyRequest.TransactType = "send"
				notifyRequest.Address = v

				err = notification.NotificationRequest(notifyRequest)
				if err != nil {
					global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
					// return
				}
			}

			if utils.HexToAddress(toAddress) == utils.HexToAddress(v) {
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
}

func SweepBlockchainPendingBlock(
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

	client.URL = constant.GetRPCUrlByNetwork(chainId)
	var rpcBlockDetail response.RPCBlockDetail
	var jsonRpcRequest request.JsonRpcRequest
	jsonRpcRequest.Id = 1
	jsonRpcRequest.Jsonrpc = "2.0"
	jsonRpcRequest.Method = "eth_getBlockByNumber"
	jsonRpcRequest.Params = []interface{}{"0x" + strconv.FormatInt(blockHeightInt, 16), true}

	err = client.HTTPPost(jsonRpcRequest, &rpcBlockDetail)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	height, err := utils.HexStringToInt64(rpcBlockDetail.Result.Number)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	if blockHeightInt == height {
		if len(rpcBlockDetail.Result.Transactions) > 0 {
			for _, transaction := range rpcBlockDetail.Result.Transactions {

				isMonitorTx := false

			outerCurrentTxLoop:
				for i := 0; i < len(*publicKey); i++ {
					if transaction.Input == "0x" {
						if utils.HexToAddress(transaction.From) == utils.HexToAddress((*publicKey)[i]) || utils.HexToAddress(transaction.To) == utils.HexToAddress((*publicKey)[i]) {
							isMonitorTx = true
						}
					} else {

						isSupportContract, contractName, _, _ := sweepUtils.GetContractInfo(chainId, transaction.To)
						if isSupportContract {
							if erc20.IsHandleTokenTransaction(chainId, transaction.Hash, contractName, transaction.From, transaction.To, (*publicKey)[i], transaction.Input) {
								isMonitorTx = true
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
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), fmt.Sprintf("Not the same height of block: %d - %d", blockHeightInt, height)))
	}
}
