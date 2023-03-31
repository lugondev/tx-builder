package builder

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"github.com/lugondev/tx-builder/pkg/client"
	"github.com/lugondev/tx-builder/pkg/common"
	"math/rand"
	"testing"
)

func TestBuilderBTCWallet(t *testing.T) {
	//wif, err := btcutil.DecodeWIF("cP2gB7hrFoE4AccbB1qyfcgmzDicZ8bkr3XB9GhYzMUEQNkQRRwr") // mvBSG1p12WE14xnATXSa43wd8TppUzKwha
	wif, err := btcutil.DecodeWIF("cVacJiScoPMAugWKRwMU2HVUPE4PhcJLgxVCexieWEWcTiYC8bSn")
	if err != nil {
		t.Fatal(err)
	}
	builder, err := NewTxBtcBuilder(wif.SerializePubKey(), common.Taproot, &chaincfg.TestNet3Params)
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
		SetFeeRate(1000).
		SetChangeSource(builder.SourceAddressInfo.Address).
		SetOutputs([]*Output{
			{
				Address: toAddress,
				Amount:  rand.Int63n(200) + 300,
			},
		}).
		Build()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("tx: ", hexutil.Encode(signedTx))
}
