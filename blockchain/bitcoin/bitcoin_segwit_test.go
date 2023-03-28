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
	"github.com/lugondev/tx-builder/blockchain/bitcoin/author"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"github.com/lugondev/tx-builder/pkg/common"
	"testing"
)

func TestSegwit(t *testing.T) {
	rawTx, err := CreateSegwitTx(
		"cP2gB7hrFoE4AccbB1qyfcgmzDicZ8bkr3XB9GhYzMUEQNkQRRwr",
		"mvBSG1p12WE14xnATXSa43wd8TppUzKwha",
		utxo.UnspentTxOutput{
			VOut:   1,
			TxHash: "b0f37aa5f4fdf30ad8c3ab17498c8d97ac3b754a924de2f8fe5a3e3203542f94",
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

	btcAddresses := PubkeyToAddresses(wif.PrivKey.PubKey(), &chaincfg.TestNet3Params)
	fromAddressInfo := common.GetBTCAddressInfo(btcAddresses[common.Segwit])

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
	redeemTxOut1 := wire.NewTxOut(90700-amount-212, fromAddressInfo.GetPayToAddrScript())

	redeemTx.AddTxOut(redeemTxOut)
	redeemTx.AddTxOut(redeemTxOut1)

	// now sign the transaction
	finalRawTx, err := SignSegwitTx(wif, fromAddressInfo, redeemTx)

	return finalRawTx, err
}

func SignSegwitTx(wif *btcutil.WIF, fromAddressInfo *common.BTCAddressInfo, redeemTx *wire.MsgTx) (string, error) {

	secretStore := author.NewMemorySecretStore(map[string]*btcec.PrivateKey{
		fromAddressInfo.Address: wif.PrivKey,
	}, &chaincfg.TestNet3Params)

	if err := author.AddAllInputScripts(redeemTx, [][]byte{fromAddressInfo.GetPayToAddrScript()}, []btcutil.Amount{90700}, secretStore); err != nil {
		return "", nil
	}
	var signedTx bytes.Buffer
	if err := redeemTx.Serialize(&signedTx); err != nil {
		return "", err
	}

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
}
