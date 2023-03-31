package builder

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/chain"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"github.com/lugondev/tx-builder/pkg/client"
	"github.com/lugondev/tx-builder/pkg/common"
	"testing"
)

const privKey = "cP2gB7hrFoE4AccbB1qyfcgmzDicZ8bkr3XB9GhYzMUEQNkQRRwr"

// const privKey = "cPeGCNhoftdg88EQgyrkdDVPe58d8MKfUHiyzz9eFEyKLxwxLXWb"
const toAddress = "tb1q5rvwj5fyh02ldstdk77ku0vc3g9utdq693tuet"

func TestGetBalance(t *testing.T) {
	wif, err := btcutil.DecodeWIF("cP2gB7hrFoE4AccbB1qyfcgmzDicZ8bkr3XB9GhYzMUEQNkQRRwr")

	builder, err := NewTxBtcBuilder(wif.SerializePubKey(), common.Legacy, &chaincfg.TestNet3Params)
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
	builder.SetUtxos(*utxos.ToUTXOs()).SetPrivKey(wif.PrivKey)

	fmt.Printf("balance %s: %d", builder.SourceAddressInfo.Address, builder.EstimateBalance)
}
func TestBuilder(t *testing.T) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		t.Fatal(err)
	}
	btcAddresses := chain.PubkeyToAddresses(wif.PrivKey.PubKey(), &chaincfg.TestNet3Params)
	fromAddressInfo := common.GetBTCAddressInfo(btcAddresses[common.Legacy])
	fmt.Println("address legacy: ", fromAddressInfo.Address)

	c := client.NewClient("https://blockstream.info", "", "", "")
	utxoService := utxo.BlockStreamService{Client: c}
	utxos, err := utxoService.SetAddress(fromAddressInfo.Address).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("UTXOs: ", utxos.ToUTXOs().Len())

	// Create a new transaction builder
	//builder, err := NewTxBtcBuilder(common.Legacy, &chaincfg.TestNet3Params)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//txBytes := CalculateTxBytes(fromAddressInfo.Address, float64(utxos.ToUTXOs().Len()), []string{toAddress, fromAddressInfo.Address})
	//
	//amount := int64(1231)
	//finalizedTx, err := builder.SetPrivKey(wif.PrivKey).
	//	SetFeeRate(1)
	//
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println("Finalized tx: ", hexutil.Encode(finalizedTx))
}
