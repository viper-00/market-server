package wallet

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GenerateEthereumWallet() (error, string, string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return err, "", ""
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	pKey := hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA"), "", ""
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return nil, pKey, address
}

var (
	Transfer     = "transfer"
	TransferFrom = "transferFrom"

	CreateNewContract = "createNewContract"

	knownMethods = map[string]string{
		"0xa9059cbb": Transfer,
		"0x23b872dd": TransferFrom,

		"0x694e974c": CreateNewContract,
	}
)

func GenerateEthereumCollectionContract(bindAddress string) (error, string) {
	client, err := ethclient.Dial("https://sepolia.optimism.io")
	if err != nil {
		return err, ""
	}
	defer client.Close()

	file, err := os.Open("./market.json")
	if err != nil {
		return err, ""
	}
	defer file.Close()

	contractAddress := common.HexToAddress("0xa04c49003a08485d927712c6678d828b644a013f")

	contractABI, err := abi.JSON(file)
	if err != nil {
		return err, ""
	}

	callData, err := contractABI.Pack(CreateNewContract)
	if err != nil {
		return err, ""
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return err, ""
	}

	// var returnValue *big.Int
	inputsMap := make(map[string]interface{})

	err = contractABI.UnpackIntoMap(inputsMap, CreateNewContract, result)
	if err != nil {
		return err, ""
	}

	fmt.Printf("Result of CreateNewContract: %s\n", inputsMap)

	return nil, ""

	// 	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// 	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	gasPrice, err := client.SuggestGasPrice(context.Background())
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(5))
	// 	auth.Nonce = big.NewInt(int64(nonce))
	// 	auth.Value = big.NewInt(0)     // in wei
	// 	auth.GasLimit = uint64(300000) // in units
	// 	auth.GasPrice = gasPrice

	// 	input := "1.0"

	// 	address, tx, instance, err := bind.DeployContract(auth, client, []byte(""), nil)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
	// 	fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0

	// _ = instance
}
