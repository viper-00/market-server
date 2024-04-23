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
	"sync"
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

	if rpcBlockInfo.Result.Number == "" {
		return
	}

	height, err := utils.HexStringToInt64(rpcBlockInfo.Result.Number)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	if height > 0 {
		setup.SetupLatestBlockHeight(context.Background(), chainId, int64(height))
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
		return
	}

	if *sweepBlockHeight >= *cacheBlockHeight {
		SetupLatestBlockHeight(client, chainId)
		setup.UpdateCacheBlockHeight(context.Background(), chainId)
		setup.UpdatePublicKey(context.Background(), chainId)
		return
	}

	var wg sync.WaitGroup
	mutex := sync.Mutex{}

	var (
		numWorkers = 20
	)

	if *sweepBlockHeight <= *cacheBlockHeight {
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				mutex.Lock()
				currentHeight := *sweepBlockHeight
				if currentHeight > *cacheBlockHeight {
					mutex.Unlock()
					return
				}
				*sweepBlockHeight++
				mutex.Unlock()

				if chainId == constant.ETH_MAINNET {
					err := SweepBlockchainTransactionCoreForEthereum(client, chainId, publicKey, sweepCount, currentHeight, constantSweepBlock, constantPendingBlock, constantPendingTransaction)
					if err != nil {
						// global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))

						_, err = global.MARKET_REDIS.RPush(context.Background(), constantPendingBlock, currentHeight).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
						}
					}
				} else {
					err := SweepBlockchainTransactionCore(client, chainId, publicKey, sweepCount, currentHeight, constantSweepBlock, constantPendingBlock, constantPendingTransaction)
					if err != nil {
						// global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))

						_, err = global.MARKET_REDIS.RPush(context.Background(), constantPendingBlock, currentHeight).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
						}
					}
				}
			}()
		}

		wg.Wait()

		_, err := global.MARKET_REDIS.Set(context.Background(), constantSweepBlock, *sweepBlockHeight+1, 0).Result()
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return
		}
	}

}

func SweepBlockchainTransactionCore(client MARKET_Client.Client,
	chainId int,
	publicKey *[]string,
	sweepCount *map[int64]int,
	sweepBlockHeight int64,
	constantSweepBlock, constantPendingBlock, constantPendingTransaction string) error {
	defer utils.HandlePanic()

	var err error

	client.URL = constant.GetRPCUrlByNetwork(chainId)
	var rpcBlockDetail response.RPCBlockDetail
	var jsonRpcRequest request.JsonRpcRequest
	jsonRpcRequest.Id = 1
	jsonRpcRequest.Jsonrpc = "2.0"
	jsonRpcRequest.Method = "eth_getBlockByNumber"
	jsonRpcRequest.Params = []interface{}{"0x" + strconv.FormatInt(sweepBlockHeight, 16), true}

	err = client.HTTPPost(jsonRpcRequest, &rpcBlockDetail)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return err
	}

	if rpcBlockDetail.Result.Number == "" {
		err = fmt.Errorf("can not get the number: %s", rpcBlockDetail.Result.Number)
		return err
	}

	height, err := utils.HexStringToInt64(rpcBlockDetail.Result.Number)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return err
	}

	if sweepBlockHeight == int64(height) {
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
						isSupportContract, contractName, _, _ := sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, transaction.To)
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
							return err
						}

						for _, redisTx := range redisTxs {
							if redisTx == transaction.Hash {
								continue outerCurrentTxLoop
							}
						}

						_, err = global.MARKET_REDIS.RPush(context.Background(), constantPendingTransaction, transaction.Hash).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return err
						}
						break
					}
				}
			}
		}

		return nil
	} else {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), fmt.Sprintf("Not the same height of block: %d - %d", sweepBlockHeight, height)))
		return errors.New("not the same height of block")
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
			time.Sleep(2 * time.Second)
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

	if chainId == constant.ETH_MAINNET {
		err = handleEthereumTx(client, chainId, publicKey, notifyRequest, rpcDetail)
	} else {
		if rpcDetail.Result.Input == "0x" {
			err = handleERC20(chainId, publicKey, notifyRequest, rpcDetail.Result.From, rpcDetail.Result.To, rpcDetail.Result.Hash, rpcDetail.Result.Input, rpcDetail.Result.Value)
		} else {
			err = handleERC20(chainId, publicKey, notifyRequest, rpcDetail.Result.From, rpcDetail.Result.To, rpcDetail.Result.Hash, rpcDetail.Result.Input, rpcDetail.Result.Value)
		}
	}

	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	_, err = global.MARKET_REDIS.LPop(context.Background(), constantPendingTransaction).Result()
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}
}

func handleERC20(chainId int, publicKey *[]string, notifyRequest request.NotificationRequest, from, to, hash, input, value string) error {

	var (
		isSupportContract bool
		contractName      string
		decimals          int
		err               error
	)

	defer func() {
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("complete %s -> %s", constant.GetChainName(chainId), err.Error()))
		}
	}()

	if input == "0x" {
		isSupportContract, contractName, _, decimals = sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, "0x0000000000000000000000000000000000000000")
	} else {
		isSupportContract, contractName, _, decimals = sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, to)
	}

	if !isSupportContract {
		err = errors.New("can not find the contract")
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return err
	}

	fromAddress := from
	notifyRequest.FromAddress = fromAddress

	if decimals == 0 {
		err = errors.New("decimals can not be 0")
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return err
	}

	if !(input == "0x") {
		methodName, decodeFromAddress, decodeToAddress, transactionValue, err := erc20.DecodeERC20TransactionInputData(chainId, hash, input)
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return err
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
		toAddress := to
		notifyRequest.ToAddress = toAddress

		value, err := utils.HexStringToInt64(value)
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
			return err
		}
		notifyRequest.Amount = utils.CalculateBalance(big.NewInt(0).SetUint64(value), decimals)
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

	return nil
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

	if rpcBlockDetail.Result.Number == "" {
		err = fmt.Errorf("can not get the number: %s", rpcBlockDetail.Result.Number)
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	height, err := utils.HexStringToInt64(rpcBlockDetail.Result.Number)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return
	}

	if blockHeightInt == int64(height) {
		if len(rpcBlockDetail.Result.Transactions) > 0 {
			for _, transaction := range rpcBlockDetail.Result.Transactions {

				isMonitorTx := false

				if chainId == constant.ETH_MAINNET {

					var infos response.RPCInnerTxInfo
					client.URL = constant.GetInnerTxRPCUrlByNetwork(chainId)
					payload := map[string]interface{}{
						"id":      1,
						"jsonrpc": "2.0",
						"method":  "debug_traceTransaction",
						"params": []interface{}{
							transaction.Hash,
							map[string]interface{}{
								"tracer": "callTracer",
							},
						},
					}

					err = client.HTTPPost(payload, &infos)
					if err != nil {
						global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
						return
					}

				outerETHCurrentTxLoop:
					for i := 0; i < len(*publicKey); i++ {

						if infos.Result.Input == "0x" {
							if utils.HexToAddress(infos.Result.From) == utils.HexToAddress((*publicKey)[i]) || utils.HexToAddress(infos.Result.To) == utils.HexToAddress((*publicKey)[i]) {
								isMonitorTx = true
							}
						} else {
							isSupportContract, contractName, _, _ := sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, infos.Result.To)
							if isSupportContract {
								if erc20.IsHandleTokenTransaction(chainId, transaction.Hash, contractName, infos.Result.From, infos.Result.To, (*publicKey)[i], infos.Result.Input) {
									isMonitorTx = true
								}
							}
						}

						for _, v := range infos.Result.Calls {
							if v.Type != "CALL" {
								continue
							}

							if v.Input == "0x" {
								if utils.HexToAddress(v.From) == utils.HexToAddress((*publicKey)[i]) || utils.HexToAddress(v.To) == utils.HexToAddress((*publicKey)[i]) {
									isMonitorTx = true
								}
							} else {
								isSupportContract, contractName, _, _ := sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, v.To)
								if isSupportContract {
									if erc20.IsHandleTokenTransaction(chainId, transaction.Hash, contractName, v.From, v.To, (*publicKey)[i], v.Input) {
										isMonitorTx = true
									}
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
									continue outerETHCurrentTxLoop
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

				} else {

				outerCurrentTxLoop:
					for i := 0; i < len(*publicKey); i++ {
						if transaction.Input == "0x" {
							if utils.HexToAddress(transaction.From) == utils.HexToAddress((*publicKey)[i]) || utils.HexToAddress(transaction.To) == utils.HexToAddress((*publicKey)[i]) {
								isMonitorTx = true
							}
						} else {

							isSupportContract, contractName, _, _ := sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, transaction.To)
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
		}

		_, err = global.MARKET_REDIS.LPop(context.Background(), constantPendingBlock).Result()
		if err != nil {
			global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		}
	} else {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), fmt.Sprintf("Not the same height of block: %d - %d", blockHeightInt, height)))
	}
}

func SweepBlockchainTransactionCoreForEthereum(client MARKET_Client.Client,
	chainId int,
	publicKey *[]string,
	sweepCount *map[int64]int,
	sweepBlockHeight int64,
	constantSweepBlock, constantPendingBlock, constantPendingTransaction string) error {
	defer utils.HandlePanic()

	var err error

	client.URL = constant.GetRPCUrlByNetwork(chainId)
	var rpcBlockDetail response.RPCBlockDetail
	var jsonRpcRequest request.JsonRpcRequest
	jsonRpcRequest.Id = 1
	jsonRpcRequest.Jsonrpc = "2.0"
	jsonRpcRequest.Method = "eth_getBlockByNumber"
	jsonRpcRequest.Params = []interface{}{"0x" + strconv.FormatInt(sweepBlockHeight, 16), true}

	err = client.HTTPPost(jsonRpcRequest, &rpcBlockDetail)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return err
	}

	if rpcBlockDetail.Result.Number == "" {
		err = fmt.Errorf("can not get the number: %s", rpcBlockDetail.Result.Number)
		return err
	}

	height, err := utils.HexStringToInt64(rpcBlockDetail.Result.Number)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
		return err
	}

	if sweepBlockHeight == int64(height) {

		if len(rpcBlockDetail.Result.Transactions) > 0 {
			for _, transaction := range rpcBlockDetail.Result.Transactions {

				isMonitorTx := false

				var infos response.RPCInnerTxInfo
				client.URL = constant.GetInnerTxRPCUrlByNetwork(chainId)
				payload := map[string]interface{}{
					"id":      1,
					"jsonrpc": "2.0",
					"method":  "debug_traceTransaction",
					"params": []interface{}{
						transaction.Hash,
						map[string]interface{}{
							"tracer": "callTracer",
						},
					},
				}

				err := client.HTTPPost(payload, &infos)
				if err != nil {
					global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
					return err
				}

			outerCurrentTxLoop:
				for i := 0; i < len(*publicKey); i++ {

					// if infos.Result.Type != "CALL" {
					// 	continue outerCurrentTxLoop
					// }

					if infos.Result.Input == "0x" {
						if utils.HexToAddress(infos.Result.From) == utils.HexToAddress((*publicKey)[i]) || utils.HexToAddress(infos.Result.To) == utils.HexToAddress((*publicKey)[i]) {
							isMonitorTx = true
						}
					} else {
						isSupportContract, contractName, _, _ := sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, infos.Result.To)
						if isSupportContract {
							if erc20.IsHandleTokenTransaction(chainId, transaction.Hash, contractName, infos.Result.From, infos.Result.To, (*publicKey)[i], infos.Result.Input) {
								isMonitorTx = true
							}
						}
					}

					for _, v := range infos.Result.Calls {
						if v.Type != "CALL" {
							continue
						}

						if v.Input == "0x" {
							if utils.HexToAddress(v.From) == utils.HexToAddress((*publicKey)[i]) || utils.HexToAddress(v.To) == utils.HexToAddress((*publicKey)[i]) {
								isMonitorTx = true
							}
						} else {
							isSupportContract, contractName, _, _ := sweepUtils.GetContractInfoByChainIdAndContractAddress(chainId, v.To)
							if isSupportContract {
								if erc20.IsHandleTokenTransaction(chainId, transaction.Hash, contractName, v.From, v.To, (*publicKey)[i], v.Input) {
									isMonitorTx = true
								}
							}
						}
					}

					if isMonitorTx {
						// Determine duplicate transactions
						redisTxs, err := global.MARKET_REDIS.LRange(context.Background(), constantPendingTransaction, 0, -1).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return err
						}

						for _, redisTx := range redisTxs {
							if redisTx == transaction.Hash {
								continue outerCurrentTxLoop
							}
						}

						_, err = global.MARKET_REDIS.RPush(context.Background(), constantPendingTransaction, transaction.Hash).Result()
						if err != nil {
							global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), err.Error()))
							return err
						}
						break
					}
				}
			}
		}

		return nil
	} else {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s", constant.GetChainName(chainId), fmt.Sprintf("Not the same height of block: %d - %d", sweepBlockHeight, height)))
		return errors.New("not the same height of block")
	}
}

func handleEthereumTx(client MARKET_Client.Client, chainId int, publicKey *[]string, notifyRequest request.NotificationRequest, rpcDetail response.RPCTransactionDetail) error {

	var infos response.RPCInnerTxInfo
	client.URL = constant.GetInnerTxRPCUrlByNetwork(chainId)
	payload := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "debug_traceTransaction",
		"params": []interface{}{
			rpcDetail.Result.Hash,
			map[string]interface{}{
				"tracer": "callTracer",
			},
		},
	}

	err := client.HTTPPost(payload, &infos)
	if err != nil {
		global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s, hash: %s", constant.GetChainName(chainId), err.Error(), rpcDetail.Result.Hash))
		return err
	}

	_ = handleERC20(chainId, publicKey, notifyRequest, infos.Result.From, infos.Result.To, rpcDetail.Result.Hash, infos.Result.Input, infos.Result.Value)
	// if err != nil {
	// 	global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s, hash: %s", constant.GetChainName(chainId), err.Error(), rpcDetail.Result.Hash))
	// }

	for _, v := range infos.Result.Calls {
		if v.Type != "CALL" {
			continue
		}

		_ = handleERC20(chainId, publicKey, notifyRequest, v.From, v.To, rpcDetail.Result.Hash, v.Input, v.Value)
		// if innerErr != nil {
		// 	global.MARKET_LOG.Error(fmt.Sprintf("%s -> %s, hash: %s", constant.GetChainName(chainId), innerErr.Error(), rpcDetail.Result.Hash))
		// }
	}

	return nil
}
