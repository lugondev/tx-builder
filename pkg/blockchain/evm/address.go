package evm

import (
	"github.com/btcsuite/btcd/btcec/v2"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

type KeyAddress struct {
	Pubkey           []byte
	AddressHex       string
	Address          common2.Address
	AddressLowerCase string
}

func PubkeyToAddress(pubKey *btcec.PublicKey) KeyAddress {
	pubkey, err := crypto.DecompressPubkey(pubKey.SerializeCompressed())
	if err != nil {
		panic(err)
	}
	address := crypto.PubkeyToAddress(*pubkey).Hex()

	return KeyAddress{
		Pubkey:           pubKey.SerializeCompressed(),
		AddressHex:       address,
		Address:          common2.HexToAddress(address),
		AddressLowerCase: strings.ToLower(address),
	}
}
