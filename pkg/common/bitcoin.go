package common

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"strings"
)

type BTCChainType int

const (
	BTCMainnet BTCChainType = iota
	BTCTestnet
)

type BTCAddressInfo struct {
	Prefix  string         `json:"prefix"`
	Version string         `json:"version"`
	Type    BTCAddressType `json:"type"`
	Chain   BTCChainType   `json:"chain"`
	Address string         `json:"address"`
}

type BTCAddressType string

const (
	Nested  BTCAddressType = "nested"
	Legacy                 = "legacy"
	Segwit                 = "segwit"
	Taproot                = "taproot"
)

var BTCAddressTypes = []BTCAddressInfo{
	{Prefix: "1", Version: "p2pkh", Chain: BTCMainnet, Type: Legacy},
	{Prefix: "3", Version: "p2sh", Chain: BTCMainnet, Type: Nested},
	{Prefix: "bc1q", Version: "p2wpkh", Chain: BTCMainnet, Type: Segwit},
	{Prefix: "bc1p", Version: "p2tr", Chain: BTCMainnet, Type: Taproot},

	{Prefix: "tb1q", Version: "p2wpkh", Chain: BTCTestnet, Type: Segwit},
	{Prefix: "tb1p", Version: "p2tr", Chain: BTCTestnet, Type: Taproot},
	{Prefix: "m", Version: "p2pkh", Chain: BTCTestnet, Type: Legacy},
	{Prefix: "n", Version: "p2pkh", Chain: BTCTestnet, Type: Legacy},
	{Prefix: "2", Version: "p2sh", Chain: BTCTestnet, Type: Nested},
}

func GetBTCAddressInfo(address string) *BTCAddressInfo {
	for _, info := range BTCAddressTypes {
		if address[:len(info.Prefix)] == info.Prefix {
			info.Address = address
			return &info
		}
	}
	return nil
}

func (b *BTCAddressInfo) GetBTCRouterBlockStream() (router string) {
	router = ""
	if b.Chain == BTCTestnet {
		router = "/testnet"
	}

	return router
}

func (b *BTCAddressInfo) GetChainConfig() *chaincfg.Params {
	if b.Chain == BTCMainnet {
		return &chaincfg.MainNetParams
	} else {
		return &chaincfg.TestNet3Params
	}
}

func (b *BTCAddressInfo) GetVersion() string {
	return strings.ToUpper(b.Version)
}

func (b *BTCAddressInfo) GetBTCRouterCryptoAPIs() (router string) {
	router = "mainnet"
	if b.Chain == BTCTestnet {
		router = "testnet"
	}

	return router
}

func (b *BTCAddressInfo) GetPayToAddrScript() []byte {
	rcvAddress, _ := btcutil.DecodeAddress(b.Address, b.GetChainConfig())
	rcvScript, _ := txscript.PayToAddrScript(rcvAddress)
	return rcvScript
}
