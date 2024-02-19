package wallet

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestWallet(t *testing.T) {
	err, privateKey, address := GenerateEthereumWallet()
	t.Log(privateKey, address)
	if err != nil {
		t.Log(err.Error())
	}

	t.Fail()
}

// func TestDepolyContract(t *testing.T) {
// 	conn, err := ethclient.Dial("https://sepolia.optimism.io")
// 	if err != nil {
// 		t.Logf("Failed to connect to the Ethereum client: %v", err)
// 	}
// 	defer conn.Close()

// 	file, err := os.Open("./Storage.json")
// 	if err != nil {
// 		t.Log(err)
// 	}
// 	defer file.Close()

// 	contractData, err := abi.JSON(file)
// 	if err != nil {
// 		t.Log(err)
// 	}

// 	// Convert private key hex to private key object
// 	privateKey, err := crypto.HexToECDSA(privateKeyHex)
// 	if err != nil {
// 		t.Log(err)
// 	}

// 	// Get the public key from the private key
// 	publicKey := privateKey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		t.Log("Failed to convert public key to ECDSA")
// 	}

// 	// Get the Ethereum address from the public key
// 	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

// 	// Create an authorized transactor
// 	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155420))
// 	if err != nil {
// 		t.Log(err)
// 	}
// 	auth.From = fromAddress
// 	auth.GasLimit = uint64(3000000)
// 	auth.GasPrice = big.NewInt(1000000000) // 1 Gwei

// 	// Deploy the contract
// 	contractAddress, _, _, err := bind.DeployContract(
// 		auth,
// 		contractData,
// 		[]byte("608060405234801561000f575f80fd5b506101438061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610034575f3560e01c80632e64cec1146100385780636057361d14610056575b5f80fd5b610040610072565b60405161004d919061009b565b60405180910390f35b610070600480360381019061006b91906100e2565b61007a565b005b5f8054905090565b805f8190555050565b5f819050919050565b61009581610083565b82525050565b5f6020820190506100ae5f83018461008c565b92915050565b5f80fd5b6100c181610083565b81146100cb575f80fd5b50565b5f813590506100dc816100b8565b92915050565b5f602082840312156100f7576100f66100b4565b5b5f610104848285016100ce565b9150509291505056fea2646970667358221220c6f7d4fbedfc2b04d6e9363c8b289a6ad5694781365bdbfb585ffe4ca734b45864736f6c63430008180033"),
// 		conn,
// 	)
// 	if err != nil {
// 		t.Log(err)
// 	}

// 	t.Logf("Contract deployed at address: %s\n", contractAddress.Hex())

// 	t.Fail()
// }

func TestCallContract(t *testing.T) {
	client, err := ethclient.Dial("https://optimism-sepolia.blockpi.network/v1/rpc/public")
	if err != nil {
		t.Log(err.Error())

	}
	defer client.Close()

	file, err := os.Open("./market.json")
	if err != nil {
		t.Log(err.Error())

	}
	defer file.Close()

	contractABI, err := abi.JSON(file)
	if err != nil {
		t.Log(err.Error())

	}

	// addresses := make([]common.Address, 0)
	// addresses = append(addresses, )

	// 设置合约函数参数
	addresses := []common.Address{
		common.HexToAddress("0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a"),
		// 添加更多地址...
	}

	callData, err := contractABI.Pack(CreateNewContract, addresses)
	if err != nil {
		t.Log(err.Error())

	}

	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		t.Log(err)
	}

	fromAddress := common.HexToAddress("0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a")

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		t.Log(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Log(err)
	}

	gasLimit := uint64(1000000)

	contractAddress := common.HexToAddress("0xa04c49003a08485d927712c6678d828b644a013f")

	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, callData)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		t.Log(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		t.Log(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Log(err)
	}

	// msg := ethereum.CallMsg{
	// 	From: common.HexToAddress("0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a"),
	// 	To:   &contractAddress,
	// 	Data: callData,
	// 	Gas:  1000000,
	// }

	// // client.Client().CallContext(context.Background())

	// result, err := client.CallContract(context.Background(), msg, nil)
	// if err != nil {
	// 	t.Log(err.Error())
	// }

	// var returnValue *big.Int
	// inputsMap := make(map[string]interface{})

	// err = contractABI.UnpackIntoMap(inputsMap, CreateNewContract, result)
	// if err != nil {
	// 	t.Log(err.Error())
	// }

	// for i, v := range inputsMap {
	// 	t.Logf("sss: %s %s", i, v)
	// }

	t.Logf("tx sent: %s", signedTx.Hash().Hex())

	// t.Logf("Result of CreateNewContract: %s\n", inputsMap)

	t.Fail()
}

func TestCallErc20Contract(t *testing.T) {
	client, err := ethclient.Dial("https://optimism-sepolia.blockpi.network/v1/rpc/public")
	if err != nil {
		t.Log(err.Error())

	}
	defer client.Close()

	file, err := os.Open("./erc20.json")
	if err != nil {
		t.Log(err.Error())

	}
	defer file.Close()

	contractABI, err := abi.JSON(file)
	if err != nil {
		t.Log(err.Error())
	}

	callData, err := contractABI.Pack("decimals")
	if err != nil {
		t.Log(err.Error())
	}

	contractAddress := common.HexToAddress("0x257144beb41dd19c90aa71c7874d6a725829227b")

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
		Gas:  300000,
	}

	// 在发送交易之前，你需要解锁发起交易的账户
	// client.CallContext(context.Background(), "personal_unlockAccount", yourAccountAddress, "yourAccountPassword", nil)

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		t.Log(err.Error())
	}

	inputsMap := make(map[string]interface{})

	err = contractABI.UnpackIntoMap(inputsMap, "decimals", result)
	if err != nil {
		t.Log(err.Error())
	}

	for i, v := range inputsMap {
		t.Logf("sss: %s %s", i, v)
	}

	t.Fail()
}

func TestCallEthTransfer(t *testing.T) {
	client, err := ethclient.Dial("https://optimism-sepolia.blockpi.network/v1/rpc/public")
	if err != nil {
		t.Log(err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		t.Log(err)
	}

	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	t.Log("error casting public key to ECDSA")
	// }

	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fromAddress := common.HexToAddress("0x4e16f68b13f15b40b0313f35E01bF2e6F636eB9a")
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		t.Log(err)
	}

	value := big.NewInt(100000000000000) // in wei (1 eth)
	gasLimit := uint64(21000)            // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Log(err)
	}

	toAddress := common.HexToAddress("0x7072579d5551Af8f77C960364923f305dEB1A521")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		t.Log(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		t.Log(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Log(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())

	t.Fail()
}
