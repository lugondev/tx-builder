package utils

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
)

func ParseHexToPubkey(pubkey string) (*btcec.PublicKey, error) {
	if !IsHexString(pubkey) {
		return nil, fmt.Errorf("expected hex string")
	}
	bytes := StringToHexBytes(pubkey)
	publicKey, err := btcec.ParsePubKey(bytes)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func ParseHexToPublicKey(pubkey string) (string, error) {
	publicKey, err := ParseHexToPubkey(pubkey)
	if err != nil {
		return "", err
	}
	return HexBytesToString(publicKey.SerializeUncompressed()), nil
}

func ParseHexToCompressedPublicKey(pubkey string) (string, error) {
	publicKey, err := ParseHexToPubkey(pubkey)
	if err != nil {
		return "", err
	}
	return HexBytesToString(publicKey.SerializeCompressed()), nil
}
