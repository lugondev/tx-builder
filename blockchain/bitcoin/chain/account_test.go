package chain

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/lugondev/tx-builder/pkg/common"
	"testing"
)

func TestGetAddressFromSeed(t *testing.T) {
	mnemonic := "furnace diesel fault piano wrap surface focus saddle chuckle absent range exact"
	addresses, err := GetAddressesFromSeed(mnemonic, &chaincfg.TestNet3Params, 3)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(addresses); i++ {
		pubkeyToAddress := addresses[i]
		fmt.Println("taproot address: ", pubkeyToAddress[common.Taproot])
		fmt.Println("legacy address: ", pubkeyToAddress[common.Legacy])
		fmt.Println("nested address: ", pubkeyToAddress[common.Nested])
		fmt.Println("segwit address: ", pubkeyToAddress[common.Segwit])
		fmt.Println("========================")
	}
}

func TestGetAddressFromGeneratedPrivateKey(t *testing.T) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	privKeyWif, err := btcutil.NewWIF(privateKey, &chaincfg.MainNetParams, false)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("private key: ", privKeyWif.String())

	privateKeyToAddresses(t, privateKey)
}

func TestGetWifFromPrivateKey(t *testing.T) {
	privateKey, _ := btcec.PrivKeyFromBytes(common2.FromHex("eea6db960d8537f33c922aa13ff3442f2cfa1e97a01023b2448b3af759c6833d"))
	wif, err := btcutil.NewWIF(privateKey, &chaincfg.TestNet3Params, true)
	if err != nil {
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("wif code: ", wif.String())

	privateKeyToAddresses(t, privateKey)
}

func TestGetAddressFromPrivateKey(t *testing.T) {
	wif, err := btcutil.DecodeWIF("cVacJiScoPMAugWKRwMU2HVUPE4PhcJLgxVCexieWEWcTiYC8bSn")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("private key: ", wif.String())

	privateKeyToAddresses(t, wif.PrivKey)
}

func privateKeyToAddresses(t *testing.T, privKey *btcec.PrivateKey) {
	addressesTestnet3, err := GetAddressFromPrivate(privKey, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("taproot testnet3: ", (*addressesTestnet3)[common.Taproot])
	fmt.Println("legacy testnet3: ", (*addressesTestnet3)[common.Legacy])
	fmt.Println("nested testnet3: ", (*addressesTestnet3)[common.Nested])
	fmt.Println("segwit testnet3: ", (*addressesTestnet3)[common.Segwit])

	addressesMainnet, err := GetAddressFromPrivate(privKey, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("taproot mainnet: ", (*addressesMainnet)[common.Taproot])
	fmt.Println("legacy mainnet: ", (*addressesMainnet)[common.Legacy])
	fmt.Println("nested mainnet: ", (*addressesMainnet)[common.Nested])
	fmt.Println("segwit mainnet: ", (*addressesMainnet)[common.Segwit])
}
