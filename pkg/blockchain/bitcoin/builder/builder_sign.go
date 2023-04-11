package builder

import (
	"bytes"
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/wire"
	author2 "github.com/lugondev/tx-builder/pkg/blockchain/bitcoin/author"
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
//	if err := txauthor.AddAllInputScripts(redeemTx, [][]byte{t.SourceAddressInfo.GetPayToAddrScript()}, []btcutil.Amount{102100}, secretStore); err != nil {
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
	//	signature, err := txscript.SignatureScript(redeemTx, i, t.SourceAddressInfo.GetPayToAddrScript(), txscript.SigHashAll, t.privKey, true)
	//	if err != nil {
	//		return nil, err
	//	}
	//	redeemTx.TxIn[i].SignatureScript = signature
	//}

	secretStore := author2.NewMemorySecretStore(map[string]*btcec.PrivateKey{
		t.SourceAddressInfo.Address: t.privKey,
	}, map[string][]byte{
		t.SourceAddressInfo.Address: t.pubkey.SerializeCompressed(),
	}, t.chainCfg)

	if err := author2.AddAllInputScripts(tx, [][]byte{t.sourceScript, t.sourceScript, t.sourceScript, t.sourceScript}, t.amountsInput, secretStore); err != nil {
		return nil, err
	}

	var signedTx bytes.Buffer
	if err := tx.Serialize(&signedTx); err != nil {
		return nil, err
	}

	return signedTx.Bytes(), nil
}
