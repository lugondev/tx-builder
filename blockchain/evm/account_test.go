package evm_test

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lugondev/tx-builder/blockchain/evm"
	"github.com/tyler-smith/go-bip39"
	"testing"
)

func TestGetAddressesFromSeed(t *testing.T) {
	generateSeed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		t.Fatal(err)
	}
	mnemonic, err := bip39.NewMnemonic(generateSeed)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("mnemonic: ", mnemonic)
	addresses, err := evm.GetAddressesFromSeed("furnace diesel fault piano wrap surface focus saddle chuckle absent range exact", 3)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(addresses); i++ {
		pubkeyToAddresses := addresses[i]
		t.Log("address: ", pubkeyToAddresses.Address)
		t.Log("========================")
	}
}

func TestGetAddressFromPrivateKey(t *testing.T) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("private key: ", common.Bytes2Hex(privateKey.Serialize()))
	pubkeyToAddress, err := evm.GetAddressFromPrivate(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("address: ", pubkeyToAddress.Address)
}
