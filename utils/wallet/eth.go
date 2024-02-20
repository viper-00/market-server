package wallet

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	Transfer     = "transfer"
	TransferFrom = "transferFrom"
	Name         = "name"
	Symbol       = "symbol"
	Decimals     = "decimals"
	TotalSupply  = "totalSupply"
	BalanceOf    = "balanceOf"
	Approve      = "approve"

	CreateNewContract = "createNewContract"

	knownMethods = map[string]string{
		"0xa9059cbb": Transfer,
		"0x23b872dd": TransferFrom,

		"0x694e974c": CreateNewContract,
	}
)

func GenerateEthereumWallet() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	pKey := hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return pKey, address, nil
}

func SendEthereumCollectionContract(rpc, fromPri, fromPub, contractAddress string, bindAddresses []string, gasLimit uint64) (hash string, err error) {
	var value = big.NewInt(0)

	file, err := os.Open("./market.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	marketContractABI, err := abi.JSON(file)
	if err != nil {
		return "", err
	}

	var addresses = []common.Address{}
	for _, v := range bindAddresses {
		addresses = append(addresses, common.HexToAddress(v))
	}

	callData, err := marketContractABI.Pack(CreateNewContract, addresses)
	if err != nil {
		return "", err
	}

	hash, err = CallWalletTransactionCore(rpc, fromPri, fromPub, contractAddress, value, callData, gasLimit)
	if err != nil {
		return "", err
	}

	return
}

func CallEthTransfer(rpc, fromPri, fromPub, toAddress string, value *big.Int, gasLimit uint64) (hash string, err error) {
	var data []byte
	hash, err = CallWalletTransactionCore(rpc, fromPri, fromPub, toAddress, value, data, gasLimit)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func CallTokenTransfer(rpc, fromPri, fromPub, toAddress, tokenAddress string, tokenValue *big.Int, gasLimit uint64) (hash string, err error) {
	var value = big.NewInt(0)

	file, err := os.Open("./erc20.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	contractABI, err := abi.JSON(file)
	if err != nil {
		return "", err
	}

	data, err := contractABI.Pack(Transfer, common.HexToAddress(toAddress), tokenValue)
	if err != nil {
		return "", err
	}

	hash, err = CallWalletTransactionCore(rpc, fromPri, fromPub, tokenAddress, value, data, gasLimit)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func CallTokenTransferFrom(rpc, fromPri, fromPub, toAddress, tokenAddress string, tokenValue *big.Int, gasLimit uint64) (hash string, err error) {
	var value = big.NewInt(0)

	file, err := os.Open("./erc20.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	contractABI, err := abi.JSON(file)
	if err != nil {
		return "", err
	}

	data, err := contractABI.Pack(TransferFrom, common.HexToAddress(fromPub), common.HexToAddress(toAddress), tokenValue)
	if err != nil {
		return "", err
	}

	hash, err = CallWalletTransactionCore(rpc, fromPri, fromPub, tokenAddress, value, data, gasLimit)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func CallTokenApprove(rpc, fromPri, fromPub, approveAddress, tokenAddress string, approveValue *big.Int, gasLimit uint64) (hash string, err error) {
	var value = big.NewInt(0)

	file, err := os.Open("./erc20.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	contractABI, err := abi.JSON(file)
	if err != nil {
		return "", err
	}

	data, err := contractABI.Pack(Approve, common.HexToAddress(approveAddress), approveValue)
	if err != nil {
		return "", err
	}

	hash, err = CallWalletTransactionCore(rpc, fromPri, fromPub, tokenAddress, value, data, gasLimit)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func CallTokenName(rpc, tokenAddress string) (result interface{}, err error) {
	result, err = CallContractCore(rpc, tokenAddress, Name)
	if err != nil {
		return nil, err
	}

	return
}

func CallTokenSymbol(rpc, tokenAddress string) (result interface{}, err error) {
	result, err = CallContractCore(rpc, tokenAddress, Symbol)
	if err != nil {
		return nil, err
	}
	return
}

func CallTokenDecimals(rpc, tokenAddress string) (result interface{}, err error) {
	result, err = CallContractCore(rpc, tokenAddress, Decimals)
	if err != nil {
		return nil, err
	}
	return
}

func CallTokenTotalSupply(rpc, tokenAddress string) (result interface{}, err error) {
	result, err = CallContractCore(rpc, tokenAddress, TotalSupply)
	if err != nil {
		return nil, err
	}
	return
}

func CallTokenBalanceOf(rpc, fromPub, tokenAddress string) (result interface{}, err error) {
	result, err = CallContractCore(rpc, tokenAddress, BalanceOf, common.HexToAddress(fromPub))
	if err != nil {
		return nil, err
	}
	return
}
