package builder

import (
	"bytes"
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/lugondev/tx-builder/blockchain/bitcoin"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/author"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"github.com/lugondev/tx-builder/pkg/client"
	"github.com/lugondev/tx-builder/pkg/common"
	"github.com/status-im/keycard-go/hexutils"
	"testing"
)

func TestBuilderSignExternal(t *testing.T) {
	wif, err := btcutil.DecodeWIF("cVacJiScoPMAugWKRwMU2HVUPE4PhcJLgxVCexieWEWcTiYC8bSn")
	if err != nil {
		t.Fatal(err)
	}
	btcAddresses := bitcoin.PubkeyToAddresses(wif.PrivKey.PubKey(), &chaincfg.TestNet3Params)
	fromAddressInfo := common.GetBTCAddressInfo(btcAddresses[common.Taproot])
	fmt.Println("address: ", fromAddressInfo.Address)

	c := client.NewClient("https://blockstream.info", "", "", "")
	utxoService := utxo.BlockStreamService{Client: c}
	utxos, err := utxoService.SetAddress(fromAddressInfo.Address).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("UTXOs: ", utxos.ToUTXOs().Len())

	toAddressInfo := common.GetBTCAddressInfo(toAddress)
	changeSource := author.ChangeSource{
		NewScript: func() ([]byte, error) {
			return fromAddressInfo.GetPayToAddrScript(), nil
		},
		ScriptSize: len(fromAddressInfo.GetPayToAddrScript()),
	}
	secretStore := author.NewMemorySecretStore(map[string]*btcec.PrivateKey{
		fromAddressInfo.Address: wif.PrivKey,
	}, fromAddressInfo.GetChainConfig())

	fetchInputs := func(target btcutil.Amount) (total btcutil.Amount, inputs []*wire.TxIn,
		inputValues []btcutil.Amount, scripts [][]byte, err error) {

		p2pkhScript := fromAddressInfo.GetPayToAddrScript()

		for _, utx := range utxos.ToUTXOsArray() {
			total += btcutil.Amount(utx.Value)

			utxoHash, err := chainhash.NewHashFromStr(utx.TxHash)
			if err != nil {
				continue
			}
			outPoint := wire.NewOutPoint(utxoHash, uint32(utx.VOut))
			inputs = append(inputs, wire.NewTxIn(outPoint, nil, nil))
			inputValues = append(inputValues, btcutil.Amount(utx.Value))
			scripts = append(scripts, p2pkhScript)
		}

		return total, inputs, inputValues, scripts, nil
	}

	transaction, err := author.NewUnsignedTransaction([]*wire.TxOut{
		wire.NewTxOut(1000, toAddressInfo.GetPayToAddrScript()),
	}, btcutil.Amount(1000), fetchInputs, &changeSource)
	if err != nil {
		t.Fatal(err)
	}

	if err := transaction.AddAllInputScripts(secretStore); err != nil {
		t.Fatal(err)
	}

	var signedTx bytes.Buffer
	if err := transaction.Tx.Serialize(&signedTx); err != nil {
		t.Fatal(err)
	}

	fmt.Println("tx: ", hexutils.BytesToHex(signedTx.Bytes()))
}
