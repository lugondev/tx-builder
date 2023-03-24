package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"testing"
)

func TestTaproot(t *testing.T) {
	rawTx, err := CreateTaprootTx(
		"eea6db960d8537f33c922aa13ff3442f2cfa1e97a01023b2448b3af759c6833d",
		destinationAddress,
		Utxo{
			Idx:  1,
			TxId: "111984ad61cc14f24f157de1ba8ccf9a38af68914f786983bf7bd96e38a60159",
		},
		1000)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("raw signed transaction is: ", rawTx)
}

func CreateTaprootTx(privKey string, destination string, utxo Utxo, amount int64) (string, error) {

	privateKey, pubKey := btcec.PrivKeyFromBytes(common.FromHex(privKey))

	// use TestNet3Params for interacting with bitcoin testnet
	// if we want to interact with main net should use MainNetParams
	addrPubKey, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(pubKey)), &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}
	fmt.Println("address: ", addrPubKey.EncodeAddress())
	// extracting destination address as []byte from function argument (destination string)
	destinationAddr, err := btcutil.DecodeAddress(destination, &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}
	fmt.Println("destination address: ", destinationAddr.EncodeAddress())

	destinationAddrByte, err := txscript.PayToAddrScript(destinationAddr)
	if err != nil {
		return "", err
	}

	redeemTx := wire.NewMsgTx(wire.TxVersion)
	utxoHash, err := chainhash.NewHashFromStr(utxo.TxId)
	if err != nil {
		return "", err
	}

	outPoint := wire.NewOutPoint(utxoHash, utxo.Idx)

	// making the input, and adding it to transaction
	txIn := wire.NewTxIn(outPoint, nil, nil)
	redeemTx.AddTxIn(txIn)

	// adding the destination address and the amount to
	// the transaction as output
	redeemTxOut := wire.NewTxOut(amount, destinationAddrByte)
	fmt.Println("GetPayToAddrScript: ", hexutil.Encode(GetPayToAddrScript(addrPubKey.EncodeAddress())))
	redeemTxOut1 := wire.NewTxOut(7600-amount-200, GetPayToAddrScript(addrPubKey.EncodeAddress()))

	redeemTx.AddTxOut(redeemTxOut)
	redeemTx.AddTxOut(redeemTxOut1)

	wif, err := btcutil.NewWIF(privateKey, &chaincfg.TestNet3Params, true)
	if err != nil {
		return "", err
	}
	// now sign the transaction
	finalRawTx, err := SignTaprootTx(wif, addrPubKey, redeemTx)

	return finalRawTx, err
}

func SignTaprootTx(wif *btcutil.WIF, addrPubKey *btcutil.AddressTaproot, redeemTx *wire.MsgTx) (string, error) {
	sourceScript := GetPayToAddrScript(addrPubKey.EncodeAddress())

	secretStore := NewMemorySecretStore(map[string]*btcutil.WIF{
		addrPubKey.EncodeAddress(): wif,
	}, &chaincfg.TestNet3Params)
	if err := txauthor.AddAllInputScripts(redeemTx, [][]byte{sourceScript}, []btcutil.Amount{7600}, secretStore); err != nil {
		return "", nil
	}

	var signedTx bytes.Buffer
	if err := redeemTx.Serialize(&signedTx); err != nil {
		return "", err
	}

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
}
