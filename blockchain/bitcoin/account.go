package bitcoin

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/samber/lo"
	"github.com/tyler-smith/go-bip39"
)

func GetAddressesFromSeed(mnemonic string, params *chaincfg.Params, number int) ([]KeyAddresses, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	m, err := hdkeychain.NewMaster(seed, params)
	if err != nil {
		return nil, err
	}

	// Show that the generated master node extended key is private.
	// m/82h
	purpose, _ := m.Derive(hdkeychain.HardenedKeyStart + 82)

	// m/82h/0h
	coin, _ := purpose.Derive(hdkeychain.HardenedKeyStart + 0)
	if params.Net == chaincfg.TestNet3Params.Net {
		coin, _ = purpose.Derive(hdkeychain.HardenedKeyStart + 1)
	}

	// m/82h/0h/0h
	account, _ := coin.Derive(hdkeychain.HardenedKeyStart + 0)

	// m/82h/0h/0h/0
	receiving, _ := account.Derive(0) // 0 = receiving, 1 = change

	keys := make([]int64, number)
	// m/82h/0h/0h/0/*
	addresses := lo.Map[int64, KeyAddresses](keys, func(_ int64, i int) KeyAddresses {
		index, _ := receiving.Derive(uint32(i)) // takes an unsigned integer
		pubkey, _ := index.ECPubKey()
		return PubkeyToAddresses(pubkey, params)
	})

	return addresses, nil
}

func GetAddressFromPrivate(privateKey *btcec.PrivateKey, params *chaincfg.Params) (*KeyAddresses, error) {
	pubkeyBytes := privateKey.PubKey().SerializeCompressed()
	pubkey, err := btcec.ParsePubKey(pubkeyBytes)

	if err != nil {
		return nil, err
	}
	addresses := PubkeyToAddresses(pubkey, params)

	return &addresses, nil
}
