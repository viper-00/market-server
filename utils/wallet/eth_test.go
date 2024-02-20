package wallet

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func TestWallet(t *testing.T) {
	privateKey, address, err := GenerateEthereumWallet()
	t.Log(privateKey, address)
	if err != nil {
		t.Log(err.Error())
	}

	t.Fail()
}

func TestCallContract(t *testing.T) {
	rpc := "https://optimism-sepolia.blockpi.network/v1/rpc/public"
	fromPri := ""
	fromPub := "0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a"
	marketContractAddress := "0xa04c49003a08485d927712c6678d828b644a013f"
	bindAddresses := []string{"0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a"}
	gasLimit := uint64(1000000)
	hash, err := SendEthereumCollectionContract(rpc, fromPri, fromPub, marketContractAddress, bindAddresses, gasLimit)
	if err != nil {
		t.Log(err.Error())
	}

	t.Log("hash", hash)

	t.Fail()

	// client, err := ethclient.Dial("https://optimism-sepolia.blockpi.network/v1/rpc/public")
	// if err != nil {
	// 	t.Log(err.Error())

	// }
	// defer client.Close()

	// file, err := os.Open("./market.json")
	// if err != nil {
	// 	t.Log(err.Error())

	// }
	// defer file.Close()

	// contractABI, err := abi.JSON(file)
	// if err != nil {
	// 	t.Log(err.Error())

	// }

	// // addresses := make([]common.Address, 0)
	// // addresses = append(addresses, )

	// // 设置合约函数参数
	// addresses := []common.Address{
	// 	common.HexToAddress("0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a"),
	// 	// 添加更多地址...
	// }

	// callData, err := contractABI.Pack(CreateNewContract, addresses)
	// if err != nil {
	// 	t.Log(err.Error())

	// }

	// privateKey, err := crypto.HexToECDSA("")
	// if err != nil {
	// 	t.Log(err)
	// }

	// fromAddress := common.HexToAddress("0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a")

	// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	t.Log(err)
	// }

	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	t.Log(err)
	// }

	// gasLimit := uint64(1000000)

	// contractAddress := common.HexToAddress("0xa04c49003a08485d927712c6678d828b644a013f")

	// tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, callData)

	// chainID, err := client.NetworkID(context.Background())
	// if err != nil {
	// 	t.Log(err)
	// }

	// signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	// if err != nil {
	// 	t.Log(err)
	// }

	// err = client.SendTransaction(context.Background(), signedTx)
	// if err != nil {
	// 	t.Log(err)
	// }

	// t.Logf("tx sent: %s", signedTx.Hash().Hex())

	// t.Fail()
}

func TestCallEthTransfer(t *testing.T) {
	rpc := "https://optimism-sepolia.blockpi.network/v1/rpc/public"
	fromPri := ""
	fromPub := "0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a"
	toAddress := "0x76F04327adA8CE7c7959BB0592329840cA6BD59C"
	value := big.NewInt(1000000)
	var data []byte
	var gasLimit uint64 = 96000

	hash, err := CallWalletTransactionCore(rpc, fromPri, fromPub, toAddress, value, data, gasLimit)
	if err != nil {
		t.Log(err.Error())
	}

	t.Log("hash", hash)

	t.Fail()
}

func TestCallTokenTransfer(t *testing.T) {
	rpc := "https://optimism-sepolia.blockpi.network/v1/rpc/public"
	fromPri := ""
	fromPub := "0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a"
	toAddress := "0x76F04327adA8CE7c7959BB0592329840cA6BD59C"
	value := big.NewInt(0)
	var data []byte
	var gasLimit uint64 = 96000

	tokenAddress := "0x257144bEb41Dd19c90aa71c7874D6a725829227b"

	file, err := os.Open("./erc20.json")
	if err != nil {
		t.Log(err.Error())

	}
	defer file.Close()

	contractABI, err := abi.JSON(file)
	if err != nil {
		t.Log(err.Error())
	}

	data, err = contractABI.Pack("transfer", common.HexToAddress(toAddress), big.NewInt(1000000000000000000))
	if err != nil {
		t.Log(err.Error())
	}

	hash, err := CallWalletTransactionCore(rpc, fromPri, fromPub, tokenAddress, value, data, gasLimit)
	if err != nil {
		t.Log(err.Error())
	}

	t.Log("hash", hash)

	t.Fail()
}

func TestCallTokenTransferFrom(t *testing.T) {
	rpc := "https://optimism-sepolia.blockpi.network/v1/rpc/public"
	fromPri := ""
	fromPub := "0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a"
	toAddress := "0x76F04327adA8CE7c7959BB0592329840cA6BD59C"
	value := big.NewInt(0)
	var data []byte
	var gasLimit uint64 = 96000

	tokenAddress := "0x257144bEb41Dd19c90aa71c7874D6a725829227b"

	file, err := os.Open("./erc20.json")
	if err != nil {
		t.Log(err.Error())

	}
	defer file.Close()

	contractABI, err := abi.JSON(file)
	if err != nil {
		t.Log(err.Error())
	}

	data, err = contractABI.Pack("transferFrom", common.HexToAddress(fromPub), common.HexToAddress(toAddress), big.NewInt(1000000000000000000))
	if err != nil {
		t.Log(err.Error())
	}

	hash, err := CallWalletTransactionCore(rpc, fromPri, fromPub, tokenAddress, value, data, gasLimit)
	if err != nil {
		t.Log(err.Error())
	}

	t.Log("hash", hash)

	t.Fail()
}

func TestCallTokenApprove(t *testing.T) {
	rpc := "https://optimism-sepolia.blockpi.network/v1/rpc/public"
	fromPri := ""
	fromPub := "0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a"
	approveAddress := "0x76F04327adA8CE7c7959BB0592329840cA6BD59C"
	value := big.NewInt(0)
	var data []byte
	var gasLimit uint64 = 96000

	tokenAddress := "0x257144bEb41Dd19c90aa71c7874D6a725829227b"

	file, err := os.Open("./erc20.json")
	if err != nil {
		t.Log(err.Error())

	}
	defer file.Close()

	contractABI, err := abi.JSON(file)
	if err != nil {
		t.Log(err.Error())
	}

	data, err = contractABI.Pack("approve", common.HexToAddress(approveAddress), big.NewInt(1000000000000000000))
	if err != nil {
		t.Log(err.Error())
	}

	hash, err := CallWalletTransactionCore(rpc, fromPri, fromPub, tokenAddress, value, data, gasLimit)
	if err != nil {
		t.Log(err.Error())
	}

	t.Log("hash", hash)

	t.Fail()
}

func TestCallTokenName(t *testing.T) {
	rpc := "https://optimism-sepolia.blockpi.network/v1/rpc/public"
	tokenAddress := "0x257144bEb41Dd19c90aa71c7874D6a725829227b"

	result, err := CallContractCore(rpc, tokenAddress, "name")
	if err != nil {
		t.Log(err.Error())
	}

	t.Log(result)

	t.Fail()
}

func TestCallTokenSymbol(t *testing.T) {
	rpc := "https://optimism-sepolia.blockpi.network/v1/rpc/public"
	tokenAddress := "0x257144bEb41Dd19c90aa71c7874D6a725829227b"

	result, err := CallContractCore(rpc, tokenAddress, "symbol")
	if err != nil {
		t.Log(err.Error())
	}

	t.Log(result)

	t.Fail()
}

func TestCallTokenDecimals(t *testing.T) {
	rpc := "https://optimism-sepolia.blockpi.network/v1/rpc/public"
	tokenAddress := "0x257144bEb41Dd19c90aa71c7874D6a725829227b"

	result, err := CallContractCore(rpc, tokenAddress, "decimals")
	if err != nil {
		t.Log(err.Error())
	}

	t.Log(result)
}

func TestCallTokenTotalSupply(t *testing.T) {
	rpc := "https://optimism-sepolia.blockpi.network/v1/rpc/public"
	tokenAddress := "0x257144bEb41Dd19c90aa71c7874D6a725829227b"

	result, err := CallContractCore(rpc, tokenAddress, "totalSupply")
	if err != nil {
		t.Log(err.Error())
	}

	t.Log(result)

	t.Fail()
}

func TestCallTokenBalanceOf(t *testing.T) {
	rpc := "https://optimism-sepolia.blockpi.network/v1/rpc/public"
	fromPub := "0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a"
	tokenAddress := "0x257144bEb41Dd19c90aa71c7874D6a725829227b"

	result, err := CallContractCore(rpc, tokenAddress, "balanceOf", common.HexToAddress(fromPub))
	if err != nil {
		t.Log(err.Error())
	}

	t.Log(result)

	t.Fail()
}

func TestCallGetContractAddressFromHash(t *testing.T) {
	rpc := "https://optimism-sepolia.blockpi.network/v1/rpc/public"
	hash := "0xea48f36c01add64c6bda581787cf28c6a18a8d5975c2ab2e12add3742dbf3fe6"

	receipt, err := GetTransactionByHash(rpc, hash)
	if err != nil {
		t.Log(err.Error())
	}

	callContractAddress := "0xA04C49003a08485D927712c6678d828b644a013f"

	for _, v := range receipt.Logs {
		if common.HexToAddress(v.Topics[0].Hex()) == common.HexToAddress("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0") &&
			common.HexToAddress(v.Topics[1].Hex()) == common.HexToAddress("0x0000000000000000000000000000000000000000") &&
			common.HexToAddress(v.Topics[2].Hex()) == common.HexToAddress(callContractAddress) {
			t.Log(v.Address.String())
		}
	}

	t.Fail()
}
