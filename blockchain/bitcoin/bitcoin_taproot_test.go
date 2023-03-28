package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/author"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"github.com/lugondev/tx-builder/pkg/common"
	"testing"
)

func TestTaproot(t *testing.T) {
	rawTx, err := CreateTaprootTx(
		"cVacJiScoPMAugWKRwMU2HVUPE4PhcJLgxVCexieWEWcTiYC8bSn",
		destinationAddress,
		utxo.UnspentTxOutput{
			VOut:   1,
			TxHash: "23ce41a1bdd37836b4beedb6ebd51e485834fcddade96a61f9a56b56c088e5e4",
		},
		1000)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("raw signed transaction is: ", rawTx)
}

func CreateTaprootTx(privKey string, destination string, utxo utxo.UnspentTxOutput, amount int64) (string, error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}
	//privateKey, pubKey := btcec.PrivKeyFromBytes(common2.FromHex(privKey))

	btcAddresses := PubkeyToAddresses(wif.PrivKey.PubKey(), &chaincfg.TestNet3Params)
	fromAddressInfo := common.GetBTCAddressInfo(btcAddresses[common.Taproot])
	// use TestNet3Params for interacting with bitcoin testnet
	// if we want to interact with main net should use MainNetParams
	//addrPubKey, err := btcutil.NewAddressTaproot(
	//	schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(pubKey)), &chaincfg.TestNet3Params)
	//if err != nil {
	//	return "", err
	//}
	fmt.Println("address: ", fromAddressInfo.Address)
	// extracting destination address as []byte from function argument (destination string)
	toAddressInfo := common.GetBTCAddressInfo(destination)

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
	redeemTxOut := wire.NewTxOut(amount, toAddressInfo.GetPayToAddrScript())
	redeemTxOut1 := wire.NewTxOut(1657-amount-200, fromAddressInfo.GetPayToAddrScript())

	redeemTx.AddTxOut(redeemTxOut)
	redeemTx.AddTxOut(redeemTxOut1)

	// now sign the transaction
	finalRawTx, err := SignTaprootTx(wif.PrivKey, fromAddressInfo, redeemTx)

	return finalRawTx, err
}

func SignTaprootTx(privateKey *btcec.PrivateKey, fromAddressInfo *common.BTCAddressInfo, redeemTx *wire.MsgTx) (string, error) {
	secretStore := author.NewMemorySecretStore(map[string]*btcec.PrivateKey{
		fromAddressInfo.Address: privateKey,
	}, fromAddressInfo.GetChainConfig())

	if err := author.AddAllInputScripts(redeemTx, [][]byte{fromAddressInfo.GetPayToAddrScript()}, []btcutil.Amount{1657}, secretStore); err != nil {
		return "", nil
	}

	var signedTx bytes.Buffer
	if err := redeemTx.Serialize(&signedTx); err != nil {
		return "", err
	}

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
}
