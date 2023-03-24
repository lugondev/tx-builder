package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"testing"
)

const destinationAddress = "mkF4Rkh9bQoUujuk5zJnvcamXvTpUgSNss"

func GetPayToAddrScript(address string) []byte {
	rcvAddress, _ := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	rcvScript, _ := txscript.PayToAddrScript(rcvAddress)
	return rcvScript
}

type Utxo struct {
	Idx  uint32
	TxId string
}

func TestCC(t *testing.T) {
	rawTx, err := CreateTx(
		"cVojzAq1juw95rKh8khRVSLAdwJe2Z2CwP8t9tvWbQ9zKM8jNcXB",
		destinationAddress,
		Utxo{
			Idx:  0,
			TxId: "57739a1d6f8964443f38892ef94daaa95d53008fa7c70fe6e8e899a4dfe76538",
		},
		1000)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("raw signed transaction is: ", rawTx)
}

func CreateTx(privKey string, destination string, utxo Utxo, amount int64) (string, error) {

	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}

	// use TestNet3Params for interacting with bitcoin testnet
	// if we want to interact with main net should use MainNetParams
	addrPubKey, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}

	/*
	 * 1 or unit-amount in Bitcoin is equal to 1 satoshi and 1 Bitcoin = 100000000 satoshi
	 */

	// checking for sufficiency of account
	//if balance < amount {
	//	return "", fmt.Errorf("the balance of the account is not sufficient")
	//}

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

	// creating a new bitcoin transaction, different sections of the tx, including
	// input list (contain UTXOs) and outputlist (contain destination address and usually our address)
	// in next steps, sections will be field and pass to sign
	redeemTx := wire.NewMsgTx(wire.TxVersion)

	utxoHash, err := chainhash.NewHashFromStr(utxo.TxId)
	if err != nil {
		return "", err
	}

	// the second argument is vout or Tx-index, which is the index
	// of spending UTXO in the transaction that TxId referred to
	// in this case is 1, but can vary different numbers
	outPoint := wire.NewOutPoint(utxoHash, utxo.Idx)

	// making the input, and adding it to transaction
	txIn := wire.NewTxIn(outPoint, nil, nil)
	redeemTx.AddTxIn(txIn)

	// adding the destination address and the amount to
	// the transaction as output
	redeemTxOut := wire.NewTxOut(amount, destinationAddrByte)
	fmt.Println("GetPayToAddrScript: ", hexutil.Encode(GetPayToAddrScript(addrPubKey.EncodeAddress())))
	redeemTxOut1 := wire.NewTxOut(8600, GetPayToAddrScript(addrPubKey.EncodeAddress()))

	redeemTx.AddTxOut(redeemTxOut)
	redeemTx.AddTxOut(redeemTxOut1)

	// now sign the transaction
	finalRawTx, err := SignTx(wif, GetPayToAddrScript(addrPubKey.EncodeAddress()), redeemTx)

	return finalRawTx, err
}

func SignTx(wif *btcutil.WIF, sourceScript []byte, redeemTx *wire.MsgTx) (string, error) {
	fmt.Println("SerializePubKey: ", hexutil.Encode(sourceScript))
	// since there is only one input in our transaction
	// we use 0 as second argument, if the transaction
	// has more args, should pass related index
	signature, err := txscript.SignatureScript(redeemTx, 0, sourceScript, txscript.SigHashAll, wif.PrivKey, true)
	if err != nil {
		return "", nil
	}
	redeemTx.TxIn[0].SignatureScript = signature

	var signedTx bytes.Buffer
	if err := redeemTx.Serialize(&signedTx); err != nil {
		return "", err
	}

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
}
