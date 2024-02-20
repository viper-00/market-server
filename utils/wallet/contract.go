package wallet

import (
	"context"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CallWalletTransactionCore(rpc, fromPrivateKey, fromPublicKey, toPublicKey string, ethValue *big.Int, data []byte, gasLimit uint64) (hash string, err error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(fromPrivateKey)
	if err != nil {
		return
	}

	fromAddress := common.HexToAddress(fromPublicKey)
	toAddress := common.HexToAddress(toPublicKey)

	// nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return
	}

	// eth_gasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return
	}

	// eth_maxPriorityFeePerGas
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		return
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		return
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:    chainId,
		Nonce:      nonce,
		GasTipCap:  gasTipCap,
		GasFeeCap:  gasPrice,
		Gas:        gasLimit,
		To:         &toAddress,
		Value:      ethValue,
		Data:       data,
		AccessList: nil,
	})

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), privateKey)
	if err != nil {
		return
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return
	}

	return signedTx.Hash().Hex(), nil
}

func CallContractCore(rpc, contractAddress, contractFunc string, args ...interface{}) (interface{}, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	file, err := os.Open("./erc20.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	contractABI, err := abi.JSON(file)
	if err != nil {
		return nil, err
	}

	callData, err := contractABI.Pack(contractFunc, args...)
	if err != nil {
		return nil, err
	}

	ca := common.HexToAddress(contractAddress)

	msg := ethereum.CallMsg{
		To:   &ca,
		Data: callData,
	}

	callResult, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	unPackResult, err := contractABI.Unpack(contractFunc, callResult)
	if err != nil {
		return nil, err
	}

	return unPackResult, nil
}

func GetTransactionByHash(rpc, hash string) (tx *types.Receipt, err error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return
	}
	defer client.Close()

	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(hash))
	if err != nil {
		return
	}

	return receipt, nil
}
