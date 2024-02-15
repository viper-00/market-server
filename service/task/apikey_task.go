package task

import (
	"fmt"
	"market/global"
	"market/global/constant"
	"market/model/market/request"
	"market/model/market/response"
	"market/model/market/response/tatum"
	"market/utils"
	"strings"
	"time"
)

func RunApiKeyTestTask() {
	for {
		now := time.Now()

		nextHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
		durationUntilNextHour := nextHour.Sub(now)

		ticker := time.NewTicker(durationUntilNextHour)

		<-ticker.C

		RunApiKeyTestCore()
	}
}

func RunApiKeyTestCore() {
	global.MARKET_LOG.Info("---------- Run Node Testing Task ----------")

	ethNode, ethRate := testLikeEthNodeKey()
	btcNode, btcRate := testBtcNodeKey()
	ltcNode, ltcRate := testLtcNodeKey()
	tronKeys, tronRate := testTronNodeKey()

	testAllNode := []string{}
	testAllNode = append(testAllNode, "---------- Run Node Testing Task ----------")
	testAllNode = append(testAllNode, "\n\n")
	testAllNode = append(testAllNode, ethNode...)
	testAllNode = append(testAllNode, "\n")
	testAllNode = append(testAllNode, btcNode...)
	testAllNode = append(testAllNode, "\n")
	testAllNode = append(testAllNode, ltcNode...)
	testAllNode = append(testAllNode, "\n")
	testAllNode = append(testAllNode, tronKeys...)
	testAllNode = append(testAllNode, fmt.Sprintf("\n\n Total Success Rate: %.2f%%\n", (ethRate+btcRate+tronRate+ltcRate)/4*100))
	if len(testAllNode) > 0 {
		utils.InformToTelegram(strings.Join(testAllNode, ""))
	}
}

func testLikeEthNodeKey() ([]string, float64) {
	var failedUrl []string
	failedUrl = append(failedUrl, "Like Ethereum RPC Chain Testing:\n")

	var totalSuccessRate float64

	var markets []int
	if global.MARKET_CONFIG.Blockchain.SweepMainnet {
		markets = constant.LikeMainnetEthChain
	} else {
		markets = constant.LikeTestnetEthChain
	}

	for _, i := range markets {
		rpcUrl, successRate := testLikeEthByhChain(i)
		failedUrl = append(failedUrl, rpcUrl...)
		totalSuccessRate += successRate
	}

	return failedUrl, totalSuccessRate / float64(len(markets))
}

func testBtcNodeKey() ([]string, float64) {
	var failedUrl []string
	failedUrl = append(failedUrl, "Bitcoin Chain Testing:\n")

	var totalSuccessRate float64

	var markets []int
	if global.MARKET_CONFIG.Blockchain.SweepMainnet {
		markets = constant.LikeMainnetBtcChain
	} else {
		markets = constant.LikeTestnetBtcChain
	}

	for _, i := range markets {
		url, successRate := testLikeBtcByChain(i)
		failedUrl = append(failedUrl, url...)
		totalSuccessRate += successRate
	}

	return failedUrl, totalSuccessRate / float64(len(markets))
}

func testLtcNodeKey() ([]string, float64) {
	var failedUrl []string
	failedUrl = append(failedUrl, "Litecoin Chain Testing:\n")

	var totalSuccessRate float64

	var markets []int
	if global.MARKET_CONFIG.Blockchain.SweepMainnet {
		markets = constant.LikeMainnetLtcChain
	} else {
		markets = constant.LikeTestnetLtcChain
	}

	for _, i := range markets {
		url, successRate := testLikeLtcByChain(i)
		failedUrl = append(failedUrl, url...)
		totalSuccessRate += successRate
	}

	return failedUrl, totalSuccessRate / float64(len(markets))
}

func testTronNodeKey() ([]string, float64) {
	var failedUrl []string
	failedUrl = append(failedUrl, "Tron Chain Testing:\n")

	var totalSuccessRate float64

	var markets []int
	if global.MARKET_CONFIG.Blockchain.SweepMainnet {
		markets = constant.LikeMainnetTronChain
	} else {
		markets = constant.LikeTestnetTronChain
	}

	for _, i := range markets {
		url, successRate := testLikeTronByChain(i)
		failedUrl = append(failedUrl, url...)
		totalSuccessRate += successRate
	}

	return failedUrl, totalSuccessRate / float64(len(markets))
}

func testLikeTronByChain(chainId int) (status []string, successRate float64) {
	var err error
	allAPiKey := constant.GetAllHTTPKeyByNetwork(chainId)
	var successCount = 0
	if len(allAPiKey) > 0 {
		for _, v := range allAPiKey {
			client.URL = constant.TronGetBlockByNetwork(chainId)
			client.Headers = map[string]string{
				"TRON-PRO-API-KEY": v,
			}

			var blockRequest request.TronGetBlockRequest
			blockRequest.Detail = false
			var blockResponse response.TronGetBlockResponse
			err = client.HTTPPost(blockRequest, &blockResponse)
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				status = append(status, fmt.Sprintf("❌ | %s -> %s\n", constant.GetChainName(chainId), v))
				continue
			}

			if blockResponse.BlockHeader.RawData.Number > 0 {
				status = append(status, fmt.Sprintf("✅︎ | %s -> %s\n", constant.GetChainName(chainId), v))
				successCount += 1
				continue
			} else {
				status = append(status, fmt.Sprintf("❌ | %s -> %s\n", constant.GetChainName(chainId), v))
				continue
			}
		}
	}

	return status, float64(successCount) / float64(len(allAPiKey))
}

func testLikeBtcByChain(chainId int) (status []string, successRate float64) {
	var err error
	var totalNumber int
	allAPiKey := constant.GetAllTatumAPiKey(chainId)
	var successCount = 0
	if len(allAPiKey) > 0 {
		totalNumber += len(allAPiKey)
		for _, v := range allAPiKey {
			client.URL = constant.TatumGetBitcoinInfo
			client.Headers = map[string]string{
				"x-api-key": v,
			}

			var bitcoinInfoResponse tatum.TatumGetBitcoinInfo
			err = client.HTTPGet(&bitcoinInfoResponse)
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				status = append(status, fmt.Sprintf("❌ | %s -> %s\n", constant.GetChainName(chainId), v))
				continue
			}

			if bitcoinInfoResponse.Blocks > 0 {
				status = append(status, fmt.Sprintf("✅︎ | %s -> %s\n", constant.GetChainName(chainId), v))
				successCount += 1
				continue
			} else {
				status = append(status, fmt.Sprintf("❌ | %s -> %s\n", constant.GetChainName(chainId), v))
				continue
			}
		}
	}

	// mempool
	client.URL = constant.MempoolGetBlockHeightByNetwork(chainId)
	var bitcoinHeight int64
	err = client.HTTPGetUnique(&bitcoinHeight)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		status = append(status, fmt.Sprintf("❌ | %s -> %s\n", constant.GetChainName(chainId), client.URL))
	}

	if bitcoinHeight > 0 {
		status = append(status, fmt.Sprintf("✅︎ | %s -> %s\n", constant.GetChainName(chainId), client.URL))
		successCount += 1
	} else {
		status = append(status, fmt.Sprintf("❌ | %s -> %s\n", constant.GetChainName(chainId), client.URL))
	}

	totalNumber += 1

	return status, float64(successCount) / float64(totalNumber)
}

func testLikeLtcByChain(chainId int) (status []string, successRate float64) {
	var err error
	var totalNumber int
	allAPiKey := constant.GetAllTatumAPiKey(chainId)
	var successCount = 0
	if len(allAPiKey) > 0 {
		totalNumber += len(allAPiKey)
		for _, v := range allAPiKey {
			client.URL = constant.TatumGetLitecoinInfo
			client.Headers = map[string]string{
				"x-api-key": v,
			}

			var litecoinInfoResponse tatum.TatumGetLitecoinInfo
			err = client.HTTPGet(&litecoinInfoResponse)
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				status = append(status, fmt.Sprintf("❌ | %s -> %s\n", constant.GetChainName(chainId), v))
				continue
			}

			if litecoinInfoResponse.Blocks > 0 {
				status = append(status, fmt.Sprintf("✅︎ | %s -> %s\n", constant.GetChainName(chainId), v))
				successCount += 1
				continue
			} else {
				status = append(status, fmt.Sprintf("❌ | %s -> %s\n", constant.GetChainName(chainId), v))
				continue
			}
		}
	}

	if global.MARKET_CONFIG.Blockchain.SweepMainnet {
		// mempool
		client.URL = constant.MempoolGetBlockHeightByNetwork(chainId)
		var bitcoinHeight int64
		err = client.HTTPGetUnique(&bitcoinHeight)
		if err != nil {
			global.MARKET_LOG.Error(err.Error())
			status = append(status, fmt.Sprintf("❌ | %s -> %s\n", constant.GetChainName(chainId), client.URL))
		}

		if bitcoinHeight > 0 {
			status = append(status, fmt.Sprintf("✅︎ | %s -> %s\n", constant.GetChainName(chainId), client.URL))
			successCount += 1
		} else {
			status = append(status, fmt.Sprintf("❌ | %s -> %s\n", constant.GetChainName(chainId), client.URL))
		}

		totalNumber += 1
	}

	return status, float64(successCount) / float64(totalNumber)
}

func testLikeEthByhChain(chainId int) (status []string, successRate float64) {
	var err error
	allRpc := constant.GetAllRPCUrlByNetwork(chainId)
	var successCount = 0
	if len(allRpc) > 0 {
		for _, v := range allRpc {
			client.URL = v
			var rpcBlockInfo response.RPCBlockInfo
			var jsonRpcRequest request.JsonRpcRequest
			jsonRpcRequest.Id = 1
			jsonRpcRequest.Jsonrpc = "2.0"
			jsonRpcRequest.Method = "eth_getBlockByNumber"
			jsonRpcRequest.Params = []interface{}{"latest", false}
			err = client.HTTPPost(jsonRpcRequest, &rpcBlockInfo)
			if err != nil {
				global.MARKET_LOG.Error(err.Error())
				status = append(status, fmt.Sprintf("❌ | %s -> %s\n", constant.GetChainName(chainId), v))
				continue
			}

			height, err := utils.HexStringToInt64(rpcBlockInfo.Result.Number)
			if err != nil || !(height > 0) {
				global.MARKET_LOG.Error(err.Error())
				status = append(status, fmt.Sprintf("❌ | %s -> %s\n", constant.GetChainName(chainId), v))
				continue
			}

			status = append(status, fmt.Sprintf("✅︎ | %s -> %s\n", constant.GetChainName(chainId), v))
			successCount += 1
		}
	}

	return status, float64(successCount) / float64(len(allRpc))
}
