package bitcoin

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/txscript"
	"github.com/lugondev/tx-builder/pkg/common"
)

const compressed = true

func PubkeyToPubKeyHash(pubkey *btcec.PublicKey, params *chaincfg.Params) (btcutil.Address, error) {
	var raw []byte
	if compressed {
		raw = pubkey.SerializeCompressed()
	} else {
		raw = pubkey.SerializeUncompressed()
	}
	pkHash := btcutil.Hash160(raw)
	return btcutil.NewAddressPubKeyHash(pkHash, params)
}

func PubkeyToScriptHash(pubkey *btcec.PublicKey, params *chaincfg.Params) (btcutil.Address, error) {
	var raw []byte
	if compressed {
		raw = pubkey.SerializeCompressed()
	} else {
		raw = pubkey.SerializeUncompressed()
	}
	pkHash := btcutil.Hash160(raw)
	scriptSig, err := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pkHash).Script()
	if err != nil {
		return nil, err
	}
	return btcutil.NewAddressScriptHash(scriptSig, params)
}

func PubkeyToSegwit(pubkey *btcec.PublicKey, params *chaincfg.Params) (btcutil.Address, error) {
	var raw []byte
	if compressed {
		raw = pubkey.SerializeCompressed()
	} else {
		raw = pubkey.SerializeUncompressed()
	}

	pubKeyHash := btcutil.Hash160(raw)
	return btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, params)
}

func PubkeyToTaprootPubKey(pubkey *btcec.PublicKey, params *chaincfg.Params) (btcutil.Address, error) {
	tapKey := txscript.ComputeTaprootKeyNoScript(pubkey)
	return btcutil.NewAddressTaproot(schnorr.SerializePubKey(tapKey), params)
}

type KeyAddresses map[common.BTCAddressType]string

func PubkeyToAddresses(pubkey *btcec.PublicKey, params *chaincfg.Params) KeyAddresses {
	return KeyAddresses{
		common.Nested:  must(PubkeyToScriptHash(pubkey, params)).EncodeAddress(),
		common.Legacy:  must(PubkeyToPubKeyHash(pubkey, params)).EncodeAddress(),
		common.Segwit:  must(PubkeyToSegwit(pubkey, params)).EncodeAddress(),
		common.Taproot: must(PubkeyToTaprootPubKey(pubkey, params)).EncodeAddress(),
	}
}

func must(address btcutil.Address, err error) btcutil.Address {
	if err != nil {
		panic(err)
	}

	return address
}
