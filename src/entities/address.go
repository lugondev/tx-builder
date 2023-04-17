package entities

import (
	"time"
)

type AddressType string // database max length 20

const (
	AddressTypeBTCSegwit  AddressType = "btc_segwit"
	AddressTypeBTCLegacy  AddressType = "btc_legacy"
	AddressTypeBTCTaproot AddressType = "btc_taproot"
	AddressTypeBTCNested  AddressType = "btc_nested"

	AddressTypeBTCSegwitTestnet  AddressType = "btc_segwit_testnet"
	AddressTypeBTCLegacyTestnet  AddressType = "btc_legacy_testnet"
	AddressTypeBTCTaprootTestnet AddressType = "btc_taproot_testnet"
	AddressTypeBTCNestedTestnet  AddressType = "btc_nested_testnet"

	AddressTypeEVM AddressType = "evm"
)

var AddressTypes = []AddressType{
	AddressTypeBTCSegwit,
	AddressTypeBTCLegacy,
	AddressTypeBTCTaproot,
	AddressTypeBTCNested,
	AddressTypeBTCSegwitTestnet,
	AddressTypeBTCLegacyTestnet,
	AddressTypeBTCTaprootTestnet,
	AddressTypeBTCNestedTestnet,
	AddressTypeEVM,
}

type Address struct {
	Address    string
	WalletType AddressType
	WalletID   int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
