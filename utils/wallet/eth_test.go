package wallet

import (
	"crypto/ecdsa"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

const key = ""

const privateKeyHex = ""

func TestDepolyContract(t *testing.T) {
	conn, err := ethclient.Dial("https://sepolia.optimism.io")
	if err != nil {
		t.Logf("Failed to connect to the Ethereum client: %v", err)
	}
	defer conn.Close()

	file, err := os.Open("./Storage.json")
	if err != nil {
		t.Log(err)
	}
	defer file.Close()

	contractData, err := abi.JSON(file)
	if err != nil {
		t.Log(err)
	}

	// Convert private key hex to private key object
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		t.Log(err)
	}

	// Get the public key from the private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Log("Failed to convert public key to ECDSA")
	}

	// Get the Ethereum address from the public key
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Create an authorized transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155420))
	if err != nil {
		t.Log(err)
	}
	auth.From = fromAddress
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = big.NewInt(1000000000) // 1 Gwei

	// Deploy the contract
	contractAddress, _, _, err := bind.DeployContract(
		auth,
		contractData,
		[]byte("608060405234801561000f575f80fd5b506101438061001d5f395ff3fe608060405234801561000f575f80fd5b5060043610610034575f3560e01c80632e64cec1146100385780636057361d14610056575b5f80fd5b610040610072565b60405161004d919061009b565b60405180910390f35b610070600480360381019061006b91906100e2565b61007a565b005b5f8054905090565b805f8190555050565b5f819050919050565b61009581610083565b82525050565b5f6020820190506100ae5f83018461008c565b92915050565b5f80fd5b6100c181610083565b81146100cb575f80fd5b50565b5f813590506100dc816100b8565b92915050565b5f602082840312156100f7576100f66100b4565b5b5f610104848285016100ce565b9150509291505056fea2646970667358221220c6f7d4fbedfc2b04d6e9363c8b289a6ad5694781365bdbfb585ffe4ca734b45864736f6c63430008180033"),
		conn,
	)
	if err != nil {
		t.Log(err)
	}

	t.Logf("Contract deployed at address: %s\n", contractAddress.Hex())

	t.Fail()

}
