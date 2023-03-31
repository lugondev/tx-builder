package builder

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"github.com/lugondev/tx-builder/pkg/client"
	"github.com/lugondev/tx-builder/pkg/common"
	"github.com/status-im/keycard-go/hexutils"
	"testing"
)

func TestBuilderSignExternal(t *testing.T) {
	//wif, err := btcutil.DecodeWIF("cP2gB7hrFoE4AccbB1qyfcgmzDicZ8bkr3XB9GhYzMUEQNkQRRwr")
	wif, err := btcutil.DecodeWIF("cVacJiScoPMAugWKRwMU2HVUPE4PhcJLgxVCexieWEWcTiYC8bSn")
	if err != nil {
		t.Fatal(err)
	}

	builder, err := NewTxBtcBuilder(wif.SerializePubKey(), common.Nested, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}

	c := client.NewClient("https://blockstream.info", "", "", "")
	utxoService := utxo.BlockStreamService{Client: c}
	utxos, err := utxoService.SetAddress(builder.SourceAddressInfo.Address).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("UTXOs: ", utxos.ToUTXOs().Len())

	signedTx, err := builder.SetUtxos(*utxos.ToUTXOs()).
		SetPrivKey(wif.PrivKey).
		SweepTo("tb1pr375lf8f88dzkxhhecpqarp9w5580eysuycu40czz8s2phd86gss9rwnaf").
		SetFeeRate(1000).
		//SetChangeSource(builder.SourceAddressInfo.Address).
		//SetOutputs([]*Output{
		//	{
		//		Address: "n47m9mEe4vXHvgWfXFZZit38NSranTwC6f",
		//		Amount:  rand.Int63n(200) + 1500,
		//	},
		//}).
		Build()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("tx: ", hexutils.BytesToHex(signedTx))
}
