package builder

import (
	"bytes"
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/wire"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/author"
)

//func (t *TxBtc)SignSegwitTx( redeemTx *wire.MsgTx) ([]byte, error) {
//
//	addrPubKey, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.TestNet3Params)
//	if err != nil {
//		return nil, err
//	}
//	secretStore := bitcoin.NewMemorySecretStore(map[string]*btcutil.WIF{
//		addrPubKey.EncodeAddress(): wif,
//	}, &chaincfg.TestNet3Params)
//
//	if err := txauthor.AddAllInputScripts(redeemTx, [][]byte{t.FromAddressInfo.GetPayToAddrScript()}, []btcutil.Amount{102100}, secretStore); err != nil {
//		return nil, err
//	}
//	var signedTx bytes.Buffer
//	if err := redeemTx.Serialize(&signedTx); err != nil {
//		return nil, err
//	}
//
//
//	return signedTx.Bytes(), nil
//}

func (t *TxBtc) signLegacyTx(tx *wire.MsgTx) ([]byte, error) {
	if t.utxos == nil || len(t.utxos) == 0 {
		return nil, errors.New("utxos is empty")
	}
	//for i := range t.utxos {
	//	signature, err := txscript.SignatureScript(redeemTx, i, t.FromAddressInfo.GetPayToAddrScript(), txscript.SigHashAll, t.privKey, true)
	//	if err != nil {
	//		return nil, err
	//	}
	//	redeemTx.TxIn[i].SignatureScript = signature
	//}

	secretStore := author.NewMemorySecretStore(map[string]*btcec.PrivateKey{
		t.FromAddressInfo.Address: t.privKey,
	}, t.chainCfg)

	if err := author.AddAllInputScripts(tx, [][]byte{t.fromScript, t.fromScript, t.fromScript, t.fromScript}, t.amountsInput, secretStore); err != nil {
		return nil, err
	}

	var signedTx bytes.Buffer
	if err := tx.Serialize(&signedTx); err != nil {
		return nil, err
	}

	return signedTx.Bytes(), nil
}
