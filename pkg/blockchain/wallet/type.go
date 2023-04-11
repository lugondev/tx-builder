package wallet

type Blockchain string

const (
	Bitcoin  Blockchain = "bitcoin"
	Ethereum            = "ethereum"
	BSC                 = "binance-smart-chain"
	Polygon             = "polygon"
	Arbitrum            = "arbitrum"
	Optimism            = "optimism"
	Cosmos              = "cosmos"
	Polkadot            = "polkadot"
	Dogecoin            = "dogecoin"
)

var SupportChains = []Blockchain{Bitcoin, Ethereum, BSC, Polygon, Arbitrum, Optimism}

var EVMChains = []Blockchain{Ethereum, BSC, Polygon, Arbitrum, Optimism}

var FullChains = []Blockchain{Bitcoin, Ethereum, BSC, Polygon, Arbitrum, Optimism, Cosmos, Polkadot, Dogecoin}
