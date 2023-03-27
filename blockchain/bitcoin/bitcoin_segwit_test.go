package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/author"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"testing"
)

func TestSegwit(t *testing.T) {
	rawTx, err := CreateSegwitTx(
		"cP2gB7hrFoE4AccbB1qyfcgmzDicZ8bkr3XB9GhYzMUEQNkQRRwr",
		"mvBSG1p12WE14xnATXSa43wd8TppUzKwha",
		utxo.UnspentTxOutput{
			VOut:   1,
			TxHash: "4a2723a32169baab6d0bf9030d56dd156a7396883cdc37067ea2caa238a17a52",
		},
		1000)

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
	redeemTxOut1 := wire.NewTxOut(90700, GetPayToAddrScript(addrPubKey.EncodeAddress()))

	redeemTx.AddTxOut(redeemTxOut)
	redeemTx.AddTxOut(redeemTxOut1)

	// now sign the transaction
	finalRawTx, err := SignSegwitTx(wif, GetPayToAddrScript(addrPubKey.EncodeAddress()), redeemTx)

	return finalRawTx, err
}

func SignSegwitTx(wif *btcutil.WIF, sourceScript []byte, redeemTx *wire.MsgTx) (string, error) {
	fmt.Println("SerializePubKey: ", hexutil.Encode(sourceScript))

	addrPubKey, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}

	secretStore := author.NewMemorySecretStore(map[string]*btcec.PrivateKey{
		addrPubKey.EncodeAddress(): wif.PrivKey,
	}, &chaincfg.TestNet3Params)

	if err := author.AddAllInputScripts(redeemTx, [][]byte{sourceScript}, []btcutil.Amount{91900}, secretStore); err != nil {
		return "", nil
	}
	var signedTx bytes.Buffer
	if err := redeemTx.Serialize(&signedTx); err != nil {
		return "", err
	}

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
}
