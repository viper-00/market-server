package service

import (
	"errors"
	"fmt"
	"market/global"
	"market/global/constant"
	"market/model"
	"market/model/market/request"
	"market/model/market/response"
	"market/model/market/response/tatum"
	"market/utils"

	"gorm.io/gorm"
)

func (n *MService) SaveTx(chainId int, hash string) (err error) {
	if !constant.IsNetworkSupport(chainId) {
		return errors.New("do not support network")
	}

	hasWallet, err := n.HasTxByChainIdAndHash(chainId, hash)
	if err != nil {
		return
	}

	if hasWallet {
		return nil
	}

	var saveTx model.Transaction
	saveTx.ChainId = chainId
	saveTx.Hash = hash
	saveTx.Status = 1

	if err = global.MARKET_DB.Create(&saveTx).Error; err != nil {
		return
	}

	return nil
}

func (n *MService) HasTxByChainIdAndHash(chainId int, hash string) (hasWallet bool, err error) {
	var findTx model.Transaction

	err = global.MARKET_DB.Where("chain_id = ? AND hash = ?", chainId, hash).First(&findTx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	if findTx.ID > 0 {
		return true, nil
	}

	return false, nil
}

func (n *MService) SaveNotification(request request.NotificationRequest, ownId uint) (err error) {
	var noModel model.Notification
	noModel.Hash = request.Hash
	noModel.ChainId = request.Chain

	var uModel model.User
	err = global.MARKET_DB.Where("contract_address = ? AND status = 1", request.Address).First(&uModel).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}
	noModel.UserId = uModel.ID
	noModel.OwnId = ownId
	noModel.IsRead = int(constant.UNREAD)
	if request.TransactType == "send" {
		noModel.Title = fmt.Sprintf("Successfully sent %s", request.Token)
		noModel.Description = fmt.Sprintf("Successfully sent %s", request.Token)
		noModel.Content = fmt.Sprintf("%s %s %s received", request.Address, request.Amount, request.Token)
		noModel.NotificationType = string(constant.OUTGOING)
	} else if request.TransactType == "receive" {
		noModel.Title = fmt.Sprintf("Successfully received %s", request.Token)
		noModel.Description = fmt.Sprintf("Successfully received %s", request.Token)
		noModel.Content = fmt.Sprintf("%s %s %s sent", request.Address, request.Amount, request.Token)
		noModel.NotificationType = string(constant.INCOMING)
	} else {
		return errors.New("not support")
	}

	err = global.MARKET_DB.Save(&noModel).Error
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return
	}
	return nil
}

func (n *MService) SaveOwnTx(request request.NotificationRequest) (id uint, err error) {
	if !constant.IsNetworkSupport(request.Chain) {
		return 0, errors.New("do not support network")
	}

	hasWallet, err := n.HasOwnTxByNotificationObj(request)
	if err != nil {
		return
	}

	if hasWallet {
		return 0, nil
	}

	var ownTx model.OwnTransaction
	ownTx.ChainId = request.Chain
	ownTx.Hash = request.Hash
	ownTx.Address = request.Address
	ownTx.FromAddress = request.FromAddress
	ownTx.ToAddress = request.ToAddress
	ownTx.Token = request.Token
	ownTx.TransactType = request.TransactType
	ownTx.Amount = request.Amount
	ownTx.BlockTimestamp = request.BlockTimestamp
	ownTx.Status = 1

	if err = global.MARKET_DB.Create(&ownTx).Error; err != nil {
		return 0, err
	}

	return ownTx.ID, nil
}

func (n *MService) HasOwnTxByNotificationObj(request request.NotificationRequest) (hasWallet bool, err error) {
	var findOwnTx model.OwnTransaction

	err = global.MARKET_DB.Where("chain_id = ? AND hash = ? AND address = ? AND from_address = ? AND to_address = ? AND token = ? AND transact_type = ? AND amount = ? AND block_timestamp = ?",
		request.Chain, request.Hash, request.Address, request.FromAddress, request.ToAddress, request.Token, request.TransactType, request.Amount, request.BlockTimestamp).First(&findOwnTx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	if findOwnTx.ID > 0 {
		return true, nil
	}

	return false, nil
}

func (n *MService) GetOwnTxById(id string) (findOwnTx model.OwnTransaction, err error) {
	err = global.MARKET_DB.Where("id = ?", id).First(&findOwnTx).Error
	return
}

func (n *MService) GetTransactionByChainAndHash(chainId int, hash string) (interface{}, error) {
	var findTx model.Transaction

	err := global.MARKET_DB.Where("chain_id = ? AND hash = ?", chainId, hash).First(&findTx).Error
	if err == nil && findTx.ID > 0 {
		return findTx, nil
	}

	return n.getTxByChainAndHash(chainId, hash)
}

func (n *MService) GetTransactionsByChainAndAddress(req request.TransactionsByChainAndAddress) ([]model.OwnTransaction, int64, error) {
	var txs []model.OwnTransaction

	var total int64
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.MARKET_DB.Model(&model.OwnTransaction{})

	if err := db.Where("chain_id = ? AND (from_address = ? OR to_address = ?)", req.ChainId, req.Address, req.Address).Count(&total).Order("created_at desc").Offset(offset).Limit(limit).Find(&txs).Error; err != nil {
		return nil, total, err
	}

	return txs, total, nil
}

func (n *MService) getTxByChainAndHash(chainId int, hash string) (interface{}, error) {

	if constant.IsNetworkLikeEth(chainId) {
		return n.handleLikeEthChain(chainId, hash)
	}

	if constant.IsNetworkLikeTron(chainId) {
		return n.handleLinkTronChain(chainId, hash)
	}

	if constant.IsNetworkLikeBtc(chainId) {
		return n.handleLinkBtcChain(chainId, hash)
	}

	if constant.IsNetworkLikeLtc(chainId) {
		return n.handleLinkLtcChain(chainId, hash)
	}

	return nil, errors.New("network not support")
}

func (n *MService) handleLikeEthChain(chainId int, hash string) (interface{}, error) {
	var err error
	var tx model.Transaction

	client.URL = constant.GetRPCUrlByNetwork(chainId)
	var rpcDetail response.RPCTransactionDetail
	var jsonRpcRequest request.JsonRpcRequest
	jsonRpcRequest.Id = 1
	jsonRpcRequest.Jsonrpc = "2.0"
	jsonRpcRequest.Method = "eth_getTransactionByHash"
	jsonRpcRequest.Params = []interface{}{hash}

	err = client.HTTPPost(jsonRpcRequest, &rpcDetail)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	var rpcBlockInfo response.RPCBlockInfo
	jsonRpcRequest.Id = 1
	jsonRpcRequest.Jsonrpc = "2.0"
	jsonRpcRequest.Method = "eth_getBlockByNumber"
	jsonRpcRequest.Params = []interface{}{rpcDetail.Result.BlockNumber, false}

	err = client.HTTPPost(jsonRpcRequest, &rpcBlockInfo)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	tx.Hash = rpcDetail.Result.Hash
	tx.ChainId = chainId
	tx.BlockNumber, _ = utils.HexStringToInt64(rpcDetail.Result.BlockNumber)
	tx.BlockHash = rpcDetail.Result.BlockHash
	tx.From = rpcDetail.Result.From
	tx.To = rpcDetail.Result.To
	tx.Gas, _ = utils.HexStringToInt64(rpcDetail.Result.Gas)
	tx.GasPrice, _ = utils.HexStringToInt64(rpcDetail.Result.GasPrice)
	tx.Input = rpcDetail.Result.Input
	tx.MaxFeePerGas, _ = utils.HexStringToInt64(rpcDetail.Result.MaxFeePerGas)
	tx.MaxPriorityFeePerGas, _ = utils.HexStringToInt64(rpcDetail.Result.MaxPriorityFeePerGas)
	tx.Nonce, _ = utils.HexStringToInt64(rpcDetail.Result.Nonce)
	tx.TransactionIndex, _ = utils.HexStringToInt64(rpcDetail.Result.TransactionIndex)
	tx.Type, _ = utils.HexStringToInt64(rpcDetail.Result.Type)
	tx.Value, _ = utils.HexStringToInt64(rpcDetail.Result.Value)

	blockTimeStamp, _ := utils.HexStringToInt64(rpcBlockInfo.Result.Timestamp)
	tx.BlockTimestamp = int(blockTimeStamp) * 1000

	if err = global.MARKET_DB.Create(&tx).Error; err != nil {
		return nil, err
	}

	return tx, nil
}

func (n *MService) handleLinkTronChain(chainId int, hash string) (interface{}, error) {
	var err error

	client.URL = constant.TronGetTxByIdByNetwork(chainId)
	client.Headers = map[string]string{
		"TRON-PRO-API-KEY": constant.GetRandomHTTPKeyByNetwork(chainId),
	}

	var txRequest request.TronGetBlockTxByIdRequest
	txRequest.Value = hash
	var txResponse response.TronGetTxResponse
	err = client.HTTPPost(txRequest, &txResponse)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	return txResponse, nil
}

func (n *MService) handleLinkBtcChain(chainId int, hash string) (interface{}, error) {
	var err error

	client.URL = constant.TatumGetBitcoinTxByHash + hash
	client.Headers = map[string]string{
		"x-api-key": constant.GetTatumRandomKeyByNetwork(chainId),
	}
	var bitcoinTxResponse tatum.TatumBitcoinTx
	err = client.HTTPGet(&bitcoinTxResponse)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	return bitcoinTxResponse, nil
}

func (n *MService) handleLinkLtcChain(chainId int, hash string) (interface{}, error) {
	var err error

	client.URL = constant.TatumGetLitecoinTxByHash + hash
	client.Headers = map[string]string{
		"x-api-key": constant.GetTatumRandomKeyByNetwork(chainId),
	}

	var litecoinTxResponse tatum.TatumLitecoinTx
	err = client.HTTPGet(&litecoinTxResponse)
	if err != nil {
		global.MARKET_LOG.Error(err.Error())
		return nil, err
	}

	return litecoinTxResponse, nil
}
