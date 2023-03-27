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
	"github.com/btcsuite/btcwallet/wallet/txauthor"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"testing"
)

func TestSegwit(t *testing.T) {
	rawTx, err := CreateSegwitTx(
		"cP2gB7hrFoE4AccbB1qyfcgmzDicZ8bkr3XB9GhYzMUEQNkQRRwr",
		"tb1pr375lf8f88dzkxhhecpqarp9w5580eysuycu40czz8s2phd86gss9rwnaf",
		utxo.UnspentTxOutput{
			VOut:   1,
			TxHash: "34287f892662f88f68cadb4b29d51e3dcdd4241eee0f668fd254120316ba2e9c",
		},
		10000)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("raw signed transaction is: ", rawTx)
}

func CreateSegwitTx(privKey string, destination string, utxo utxo.UnspentTxOutput, amount int64) (string, error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}

	// use TestNet3Params for interacting with bitcoin testnet
	// if we want to interact with main net should use MainNetParams
	addrPubKey, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}
	fmt.Println("src address: ", addrPubKey.EncodeAddress())

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
	utxoHash, err := chainhash.NewHashFromStr(utxo.TxHash)
	if err != nil {
		return "", err
	}

	outPoint := wire.NewOutPoint(utxoHash, uint32(utxo.VOut))

	// making the input, and adding it to transaction
	txIn := wire.NewTxIn(outPoint, nil, nil)
	redeemTx.AddTxIn(txIn)

	// adding the destination address and the amount to
	// the transaction as output
	redeemTxOut := wire.NewTxOut(amount, destinationAddrByte)
	fmt.Println("GetPayToAddrScript: ", hexutil.Encode(GetPayToAddrScript(addrPubKey.EncodeAddress())))
	redeemTxOut1 := wire.NewTxOut(91900, GetPayToAddrScript(addrPubKey.EncodeAddress()))

	redeemTx.AddTxOut(redeemTxOut)
	redeemTx.AddTxOut(redeemTxOut1)

	// now sign the transaction
	finalRawTx, err := SignSegwitTx(wif, GetPayToAddrScript(addrPubKey.EncodeAddress()), redeemTx)

	return finalRawTx, err
}

func SignSegwitTx(wif *btcutil.WIF, sourceScript []byte, redeemTx *wire.MsgTx) (string, error) {
	fmt.Println("SerializePubKey: ", hexutil.Encode(sourceScript))
	// since there is only one input in our transaction
	// we use 0 as second argument, if the transaction
	// has more args, should pass related index
	//prevOutputFetcher, err := txauthor.TXPrevOutFetcher(redeemTx, [][]byte{sourceScript}, []btcutil.Amount{btcutil.Amount(103300)})
	//if err != nil {
	//	return "", nil
	//}
	//signature, err := txscript.WitnessSignature(redeemTx, txscript.NewTxSigHashes(redeemTx, prevOutputFetcher), 0, 103300, sourceScript, txscript.SigHashAll, wif.PrivKey, true)
	//if err != nil {
	//	return "", nil
	//}
	//redeemTx.TxIn[0].Witness = signature

	addrPubKey, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}
	secretStore := NewMemorySecretStore(map[string]*btcutil.WIF{
		addrPubKey.EncodeAddress(): wif,
	}, &chaincfg.TestNet3Params)

	if err := txauthor.AddAllInputScripts(redeemTx, [][]byte{sourceScript}, []btcutil.Amount{102100}, secretStore); err != nil {
		return "", nil
	}
	var signedTx bytes.Buffer
	if err := redeemTx.Serialize(&signedTx); err != nil {
		return "", err
	}

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
}
