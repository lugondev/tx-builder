package builder

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/lugondev/tx-builder/blockchain/bitcoin"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"github.com/lugondev/tx-builder/pkg/common"
	"math"
)

type TxBtc struct {
	pubkey  *btcec.PublicKey
	privKey *btcec.PrivateKey

	FromAddressInfo *common.BTCAddressInfo
	fromAddressType common.BTCAddressType
	fromScript      []byte
	chainCfg        *chaincfg.Params

	utxos        []*utxo.UnspentTxOutput
	outputs      []Output
	amountsInput []btcutil.Amount

	TxBytes         int64
	FeeRate         int64
	EstimateBalance int64
}

type Output struct {
	Amount      int64
	Address     string
	Script      []byte
	AddressInfo *common.BTCAddressInfo
}

func NewTxBtcBuilder(addressType common.BTCAddressType, chainCfg *chaincfg.Params) (*TxBtc, error) {
	return &TxBtc{
		fromAddressType: addressType,
		chainCfg:        chainCfg,
	}, nil
}

func (t *TxBtc) SetPrivKey(privKey *btcec.PrivateKey) *TxBtc {
	t.privKey = privKey
	t.pubkey = privKey.PubKey()
	addresses := bitcoin.PubkeyToAddresses(t.pubkey, t.chainCfg)
	t.FromAddressInfo = common.GetBTCAddressInfo(addresses[t.fromAddressType])

	return t
}

func (t *TxBtc) SetPubkey(pubkey []byte) *TxBtc {
	pubKey, err := btcec.ParsePubKey(pubkey)
	if err != nil {
		fmt.Println("parse pubkey error", err)
		return nil
	}
	t.pubkey = pubKey

	addresses := bitcoin.PubkeyToAddresses(pubKey, t.chainCfg)
	t.FromAddressInfo = common.GetBTCAddressInfo(addresses[t.fromAddressType])
	t.fromScript = t.FromAddressInfo.GetPayToAddrScript()

	return t
}

func (t *TxBtc) GetPubKey() *btcec.PublicKey {
	return t.pubkey
}

func (t *TxBtc) SetOutputs(outputs []Output) *TxBtc {
	for i := range outputs {
		info := common.GetBTCAddressInfo(outputs[i].Address)
		if info == nil || info.GetChainConfig().Net != t.chainCfg.Net {
			fmt.Println("address type or chain config not match")
			return nil
		}
		outputs[i].AddressInfo = info
		outputs[i].Script = info.GetPayToAddrScript()
	}

	t.outputs = outputs
	return t
}

func (t *TxBtc) SetFeeRate(fee int64) *TxBtc {
	t.FeeRate = fee
	return t
}

func (t *TxBtc) SetTxBytes(txBytes float64) *TxBtc {
	t.TxBytes = int64(math.Ceil(txBytes))
	return t
}

func (t *TxBtc) CalcFee() int64 {
	return t.FeeRate * t.TxBytes
}

func (t *TxBtc) SetUtxos(utxos []*utxo.UnspentTxOutput) *TxBtc {
	t.utxos = utxos
	t.amountsInput = make([]btcutil.Amount, len(utxos))
	for i := range utxos {
		t.EstimateBalance += utxos[i].Value
		t.amountsInput[i] = btcutil.Amount(utxos[i].Value)
	}
	return t
}

func (t *TxBtc) LegacyTx() ([]byte, error) {
	redeemTx := wire.NewMsgTx(wire.TxVersion)

	if t.utxos == nil || len(t.utxos) == 0 {
		return nil, errors.New("utxos is empty")
	}

	for i := range t.utxos {
		utxoHash, err := chainhash.NewHashFromStr(t.utxos[i].TxHash)
		if err != nil {
			return nil, err
		}
		// the second argument is vout or Tx-index, which is the index
		// of spending UTXO in the transaction that TxId referred to
		// in this case is 1, but can vary different numbers
		outPoint := wire.NewOutPoint(utxoHash, uint32(t.utxos[i].VOut))

		// making the input, and adding it to transaction
		txIn := wire.NewTxIn(outPoint, nil, nil)
		redeemTx.AddTxIn(txIn)
	}

	if t.outputs == nil || len(t.outputs) == 0 {
		return nil, errors.New("outputs is empty")
	}

	for i := range t.outputs {
		// adding the destination address and the amount to
		// the transaction as output
		redeemTxOut := wire.NewTxOut(t.outputs[i].Amount, t.outputs[i].Script)
		redeemTx.AddTxOut(redeemTxOut)
	}

	// now sign the transaction
	finalRawTx, err := t.signLegacyTx(redeemTx)

	return finalRawTx, err
}
