package evm

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/samber/lo"
	"github.com/tyler-smith/go-bip39"
)

func GetAddressesFromSeed(mnemonic string, number int) ([]KeyAddress, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	keys := make([]int64, number)

	addresses := lo.Map[int64, KeyAddress](keys, func(_ int64, i int) KeyAddress {
		path, _ := accounts.ParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", i))
		key := masterKey
		for _, n := range path {
			key, _ = key.Derive(n)
		}
		pubKey, _ := key.ECPubKey()
		return PubkeyToAddress(pubKey)
	})

	return addresses, nil
}

func GetAddressFromPrivate(privateKey *btcec.PrivateKey) (*KeyAddress, error) {
	pubkeyBytes := privateKey.PubKey().SerializeCompressed()
	pubkey, err := btcec.ParsePubKey(pubkeyBytes)

	if err != nil {
		return nil, err
	}
	addresses := PubkeyToAddress(pubkey)

	return &addresses, nil
}
