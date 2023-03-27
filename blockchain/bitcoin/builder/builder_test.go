package builder

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/lugondev/tx-builder/blockchain/bitcoin"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"github.com/lugondev/tx-builder/pkg/client"
	"github.com/lugondev/tx-builder/pkg/common"
	"github.com/status-im/keycard-go/hexutils"
	"testing"
)

// const privKey = "cP2gB7hrFoE4AccbB1qyfcgmzDicZ8bkr3XB9GhYzMUEQNkQRRwr"
const privKey = "cPeGCNhoftdg88EQgyrkdDVPe58d8MKfUHiyzz9eFEyKLxwxLXWb"
const toAddress = "mvBSG1p12WE14xnATXSa43wd8TppUzKwha"

func TestGetBalance(t *testing.T) {

	client := client.NewClient("https://blockstream.info", "", "", "")
	utxoService := utxo.BlockStreamService{Client: client}
	utxos, err := utxoService.SetAddress(toAddress).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Create a new transaction builder
	builder, err := NewTxBtcBuilder(common.Legacy, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	builder.SetUtxos(*utxos.ToUTXOs())

	fmt.Println("balance:", builder.EstimateBalance)
}
func TestBuilder(t *testing.T) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		t.Fatal(err)
	}
	btcAddresses := bitcoin.PubkeyToAddresses(wif.PrivKey.PubKey(), &chaincfg.TestNet3Params)
	fromAddressInfo := common.GetBTCAddressInfo(btcAddresses[common.Legacy])
	fmt.Println("address legacy: ", fromAddressInfo.Address)

	client := client.NewClient("https://blockstream.info", "", "", "")
	utxoService := utxo.BlockStreamService{Client: client}
	utxos, err := utxoService.SetAddress(fromAddressInfo.Address).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("UTXOs: ", utxos.ToUTXOs().Len(), string(utxos.ToUTXOs().ForceToUTXOsJSON()))

	// Create a new transaction builder
	builder, err := NewTxBtcBuilder(common.Legacy, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	builder.SetUtxos(*utxos.ToUTXOs())

	txBytes := CalculateTxBytes(fromAddressInfo.Address, float64(utxos.ToUTXOs().Len()), []string{toAddress, fromAddressInfo.Address})
	fmt.Println("txBytes: ", txBytes)

	amount := int64(1231)
	finalizedTx, err := builder.SetPrivKey(wif.PrivKey).
		SetFeeRate(1).
		SetTxBytes(txBytes).
		SetOutputs([]Output{
			{
				Amount:  amount,
				Address: toAddress,
			},
			{
				Amount:  builder.EstimateBalance - builder.CalcFee() - amount,
				Address: fromAddressInfo.Address,
			},
		}).
		LegacyTx()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Finalized tx: ", hexutils.BytesToHex(finalizedTx))
}
