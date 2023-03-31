package builder

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/author"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/chain"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"github.com/lugondev/tx-builder/pkg/common"
	"strings"
)

func NewTxBtcBuilder(pubkey []byte, addressType common.BTCAddressType, chainCfg *chaincfg.Params) (*TxBtc, error) {
	txBtc := &TxBtc{
		sourceAddressType: addressType,
		chainCfg:          chainCfg,
	}
	txBtc.SetPubkey(pubkey)

	return txBtc, nil
}

func (t *TxBtc) SetPrivKey(privKey *btcec.PrivateKey) *TxBtc {
	t.privKey = privKey
	if t.pubkey == nil {
		t.SetPubkey(privKey.PubKey().SerializeUncompressed())
	} else if t.pubkey.IsEqual(privKey.PubKey()) == false {
		return nil
	}
	pubkey := hexutil.Encode(t.pubkey.SerializeCompressed())

	t.secretStore = author.NewMemorySecretStore(map[string]*btcec.PrivateKey{
		t.SourceAddressInfo.Address: privKey,
	}, map[string]*btcec.PrivateKey{
		strings.ToLower(pubkey): privKey,
	}, t.SourceAddressInfo.GetChainConfig())

	return t
}

func (t *TxBtc) SetPubkey(pubkey []byte) *TxBtc {
	if pubkey == nil || len(pubkey) == 0 {
		return nil
	}
	pubKey, err := btcec.ParsePubKey(pubkey)
	if err != nil {
		fmt.Println("parse pubkey error", err)
		return nil
	}
	t.pubkey = pubKey

	addresses := chain.PubkeyToAddresses(pubKey, t.chainCfg)
	t.SourceAddressInfo = common.GetBTCAddressInfo(addresses[t.sourceAddressType])
	t.sourceScript = t.SourceAddressInfo.GetPayToAddrScript()

	return t
}

func (t *TxBtc) GetPubKey() *btcec.PublicKey {
	return t.pubkey
}

func (t *TxBtc) SetOutputs(outputs []*Output) *TxBtc {
	t.outputs = make([]*wire.TxOut, len(outputs))
	for i := range outputs {
		err := outputs[i].HandleAddressInfo(t.chainCfg)
		if err != nil {
			return nil
		}

		t.outputs[i] = wire.NewTxOut(outputs[i].Amount, outputs[i].GetScript())
	}

	return t
}

func (t *TxBtc) SetFeeRate(fee int64) *TxBtc {
	if fee < 1000 {
		return nil
	}
	t.FeeRate = fee
	return t
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

func (t *TxBtc) getFetchInputs() author.InputSource {
	return func(target btcutil.Amount) (total btcutil.Amount, inputs []*wire.TxIn,
		inputValues []btcutil.Amount, scripts [][]byte, err error) {

		for _, utx := range t.utxos {
			total += btcutil.Amount(utx.Value)

			utxoHash, err := chainhash.NewHashFromStr(utx.TxHash)
			if err != nil {
				continue
			}
			outPoint := wire.NewOutPoint(utxoHash, uint32(utx.VOut))
			inputs = append(inputs, wire.NewTxIn(outPoint, nil, nil))
			inputValues = append(inputValues, btcutil.Amount(utx.Value))
			scripts = append(scripts, t.sourceScript)
		}

		return total, inputs, inputValues, scripts, nil
	}
}

func (t *TxBtc) SetChangeSource(address string) *TxBtc {
	addressInfo := common.GetBTCAddressInfo(address)
	if addressInfo == nil || addressInfo.GetChainConfig().Net != t.chainCfg.Net {
		return nil
	}

	t.changeSource = &author.ChangeSource{
		NewScript: func() ([]byte, error) {
			return addressInfo.GetPayToAddrScript(), nil
		},
		ScriptSize: len(addressInfo.GetPayToAddrScript()),
	}
	return t
}

func (t *TxBtc) SweepTo(address string) *TxBtc {
	return t.SetChangeSource(address)
}

func (t *TxBtc) Build() ([]byte, error) {
	if t.utxos == nil || len(t.utxos) == 0 {
		return nil, errors.New("utxos is empty")
	}
	if t.FeeRate == 0 || t.FeeRate < 1000 {
		return nil, errors.New("fee rate is too low")
	}

	var outputs []*wire.TxOut
	if t.outputs == nil || len(t.outputs) == 0 {
		outputs = []*wire.TxOut{}
	} else {
		outputs = t.outputs
	}

	transaction, err := author.NewUnsignedTransaction(outputs, btcutil.Amount(t.FeeRate), t.getFetchInputs(), t.changeSource)
	if err != nil {
		return nil, err
	}
	if err := transaction.AddAllInputScripts(t.secretStore); err != nil {
		return nil, err
	}

	var signedTx bytes.Buffer
	if err := transaction.Tx.Serialize(&signedTx); err != nil {
		return nil, err
	}

	return signedTx.Bytes(), err
}

func (t *TxBtc) SignWithECDSA(privKey *btcec.PrivateKey, msgHash []byte) (rsv string, err error) {
	sig := ecdsa.Sign(privKey, msgHash)
	return hexutil.Encode(sig.Serialize()), nil
}
