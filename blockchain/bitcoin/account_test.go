package bitcoin_test

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/lugondev/tx-builder/blockchain/bitcoin"
	"testing"
)

func TestGetAddressFromSeed(t *testing.T) {
	//generateSeed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//mnemonic, err := bip39.NewMnemonic(generateSeed)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log("mnemonic: ", mnemonic)
	mnemonic := "furnace diesel fault piano wrap surface focus saddle chuckle absent range exact"
	addresses, err := bitcoin.GetAddressesFromSeed(mnemonic, &chaincfg.MainNetParams, 3)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(addresses); i++ {
		pubkeyToAddress := addresses[i]
		fmt.Println("taproot address: ", pubkeyToAddress.Taproot)
		fmt.Println("legacy address: ", pubkeyToAddress.Legacy)
		fmt.Println("nested address: ", pubkeyToAddress.Nested)
		fmt.Println("segwit address: ", pubkeyToAddress.Segwit)
		fmt.Println("========================")
	}
}

func TestGetAddressFromPrivateKey(t *testing.T) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	privKeyWif, err := btcutil.NewWIF(privateKey, &chaincfg.MainNetParams, false)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("private key: ", privKeyWif.String())

	addresses, err := bitcoin.GetAddressFromPrivate(privateKey, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("taproot address: ", addresses.Taproot)
	fmt.Println("legacy address: ", addresses.Legacy)
	fmt.Println("nested address: ", addresses.Nested)
	fmt.Println("segwit address: ", addresses.Segwit)
}
