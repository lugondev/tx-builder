package common

type BTCChainType int

const (
	BTCMainnet BTCChainType = iota
	BTCTestnet
)

type BTCAddressInfo struct {
	Prefix  string       `json:"prefix"`
	Type    string       `json:"type"`
	Chain   BTCChainType `json:"chain"`
	Address string       `json:"address"`
}

var BTCAddressTypes = []BTCAddressInfo{
	{Prefix: "1", Type: "p2pkh", Chain: BTCMainnet},
	{Prefix: "3", Type: "p2sh", Chain: BTCMainnet},
	{Prefix: "bc1", Type: "p2wpkh", Chain: BTCMainnet},
	{Prefix: "tb1", Type: "p2wpkh", Chain: BTCTestnet},
	{Prefix: "m", Type: "p2pkh", Chain: BTCTestnet},
	{Prefix: "n", Type: "p2pkh", Chain: BTCTestnet},
	{Prefix: "2", Type: "p2sh", Chain: BTCTestnet},
}

func GetBTCAddressType(address string) *BTCAddressInfo {
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

func (b *BTCAddressInfo) GetBTCRouterCryptoAPIs() (router string) {
	router = "mainnet"
	if b.Chain == BTCTestnet {
		router = "testnet"
	}

	return router
}
