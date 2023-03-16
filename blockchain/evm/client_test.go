package evm_test

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lugondev/tx-builder/blockchain/evm"
	"math"
	"math/big"
	"testing"
)

const (
	sampleAddress = "0xD453deF0B97911be60d0899C62c445e6c4096582"
	toAddress     = "0x01504761F5Ec308Fc0BAf3e705f31F2466535d94"
	tokenAddress  = "0x2DF9398abC26759fB88aaD3FCF04b4b9F74c01cD"

	privateKeyHex = "3d153b43d2b05ed7cbdd4262ec4600bb8def570421d97d73dd59d00b4584be0c"
)

func getClient(t *testing.T) (client *evm.Client) {
	clientMumbai, err := evm.NewClientMumbai()
	if err != nil {
		t.Fatal(err)
	}
	return clientMumbai
}

func TestGetBalance(t *testing.T) {
	client := getClient(t)
	address := common.HexToAddress(sampleAddress)
	accountBalance, err := client.AccountBalance(address)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(accountBalance.String())

	fbalance := new(big.Float)
	fbalance.SetString(accountBalance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	t.Log(ethValue.String())
}

func TestTransferNativeArbitrumGoerli(t *testing.T) {
	privateKey, pubkey := btcec.PrivKeyFromBytes(common.FromHex(privateKeyHex))
	amount := new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil)
	addressFromPubkey := evm.PubkeyToAddress(pubkey).Address
	fmt.Println("addressFromPubkey", addressFromPubkey.Hex())
	fmt.Println("amount", amount.String())

	client, err := evm.NewClient("https://endpoints.omniatech.io/v1/arbitrum/goerli/public", big.NewInt(421613))
	if err != nil {
		t.Fatal(err)
	}
	accountNonce, err := client.AccountNonce(addressFromPubkey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("accountNonce", accountNonce.String())
	gasPrice, err := client.EthClient.SuggestGasPrice(client.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("gasPrice", gasPrice.String())
	to := common.HexToAddress(toAddress)
	tx, err := client.Transfer(&evm.TxRequest{
		From:     addressFromPubkey,
		Nonce:    accountNonce,
		To:       &to,
		GasPrice: gasPrice,
		Value:    amount,
	}, func(txHash []byte) ([]byte, error) {
		return crypto.Sign(txHash, privateKey.ToECDSA())
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.SubmitTx(tx); err != nil {
		t.Fatal(err)
	}
	t.Log(tx.Hash().Hex())
}

func TestTransferNative(t *testing.T) {
	privateKey, pubkey := btcec.PrivKeyFromBytes(common.FromHex(privateKeyHex))
	amount := new(big.Int).Exp(big.NewInt(10), big.NewInt(16), nil)
	addressFromPubkey := evm.PubkeyToAddress(pubkey).Address
	fmt.Println("addressFromPubkey", addressFromPubkey.Hex())
	fmt.Println("amount", amount.String())

	client := getClient(t)
	accountNonce, err := client.AccountNonce(addressFromPubkey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("accountNonce", accountNonce.String())
	gasPrice, err := client.EthClient.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	to := common.HexToAddress(toAddress)
	tx, err := client.Transfer(&evm.TxRequest{
		From:     addressFromPubkey,
		Nonce:    accountNonce,
		To:       &to,
		GasPrice: gasPrice,
		Value:    amount,
	}, func(txHash []byte) ([]byte, error) {
		return crypto.Sign(txHash, privateKey.ToECDSA())
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := client.SubmitTx(tx); err != nil {
		t.Fatal(err)
	}
	t.Log(tx.Hash().Hex())
}
func TestCallContract(t *testing.T) {
	privateKey, pubkey := btcec.PrivKeyFromBytes(common.FromHex(privateKeyHex))
	addressFromPubkey := evm.PubkeyToAddress(pubkey).Address
	fmt.Println("addressFromPubkey", addressFromPubkey.Hex())

	client := getClient(t)
	accountNonce, err := client.AccountNonce(addressFromPubkey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("accountNonce", accountNonce.String())
	gasPrice, err := client.EthClient.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	to := common.HexToAddress(tokenAddress)
	tx, err := client.TransactContract(&evm.TxRequest{
		From:     addressFromPubkey,
		Nonce:    accountNonce,
		To:       &to,
		GasPrice: gasPrice,
		Data:     common.FromHex("0xa9059cbb00000000000000000000000001504761F5Ec308Fc0BAf3e705f31F2466535d9400000000000000000000000000000000000000000000000619ce65bcce760000"),
	}, func(txHash []byte) ([]byte, error) {
		return crypto.Sign(txHash, privateKey.ToECDSA())
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := client.SubmitTx(tx); err != nil {
		t.Fatal(err)
	}
	t.Log(tx.Hash().Hex())
}
