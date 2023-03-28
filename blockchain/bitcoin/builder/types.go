package builder

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/author"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"github.com/lugondev/tx-builder/pkg/common"
)

type TxBtc struct {
	pubkey      *btcec.PublicKey
	privKey     *btcec.PrivateKey
	secretStore author.MemorySecretStore

	SourceAddressInfo *common.BTCAddressInfo
	sourceAddressType common.BTCAddressType
	sourceScript      []byte
	chainCfg          *chaincfg.Params
	changeSource      *author.ChangeSource

	utxos        []*utxo.UnspentTxOutput
	outputs      []*wire.TxOut
	amountsInput []btcutil.Amount

	TxBytes         int64
	FeeRate         int64
	EstimateBalance int64
}

type Output struct {
	Amount      int64
	Address     string
	script      []byte
	addressInfo *common.BTCAddressInfo
}
