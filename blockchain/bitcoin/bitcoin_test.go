package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/txscript"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"github.com/lugondev/tx-builder/pkg/common"
	"testing"
)

const destinationAddress = "mkF4Rkh9bQoUujuk5zJnvcamXvTpUgSNss"

func TestCC(t *testing.T) {
	rawTx, err := CreateTx(
		"cVojzAq1juw95rKh8khRVSLAdwJe2Z2CwP8t9tvWbQ9zKM8jNcXB",
		destinationAddress,
		utxo.UnspentTxOutput{
			VOut:   0,
			TxHash: "8bcb0a72620a7f55483c7cca7bf57d0c226474299c95825cb60da292bececa50",
		},
		1111)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("raw signed transaction is: ", rawTx)
}

func CreateTx(privKey string, destination string, utxo utxo.UnspentTxOutput, amount int64) (string, error) {

	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}

	// use TestNet3Params for interacting with bitcoin testnet
	// if we want to interact with main net should use MainNetParams
	//addrPubKey, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.TestNet3Params)
	//if err != nil {
	//	return "", err
	//}

	btcAddresses := PubkeyToAddresses(wif.PrivKey.PubKey(), &chaincfg.TestNet3Params)
	fromAddressInfo := common.GetBTCAddressInfo(btcAddresses[common.Segwit])

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

	utxoHash, err := chainhash.NewHashFromStr(utxo.TxHash)
	if err != nil {
		return "", err
	}

	// the second argument is vout or Tx-index, which is the index
	// of spending UTXO in the transaction that TxId referred to
	// in this case is 1, but can vary different numbers
	outPoint := wire.NewOutPoint(utxoHash, uint32(utxo.VOut))

	// making the input, and adding it to transaction
	txIn := wire.NewTxIn(outPoint, nil, nil)
	redeemTx.AddTxIn(txIn)

	// adding the destination address and the amount to
	// the transaction as output
	redeemTxOut := wire.NewTxOut(amount, destinationAddrByte)
	//fmt.Println("GetPayToAddrScript: ", hexutil.Encode(GetPayToAddrScript(addrPubKey.EncodeAddress())))
	redeemTxOut1 := wire.NewTxOut(8600, fromAddressInfo.GetPayToAddrScript())

	redeemTx.AddTxOut(redeemTxOut)
	redeemTx.AddTxOut(redeemTxOut1)

	// now sign the transaction
	finalRawTx, err := SignTx(wif, fromAddressInfo.GetPayToAddrScript(), redeemTx)

	return finalRawTx, err
}

func SignTx(wif *btcutil.WIF, sourceScript []byte, redeemTx *wire.MsgTx) (string, error) {
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
	fmt.Println("hexSignedTx: ", len(signedTx.Bytes()))
	return hexSignedTx, nil
}
