package builder

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/lugondev/tx-builder/pkg/blockchain/bitcoin/chain"
	"github.com/lugondev/tx-builder/pkg/blockchain/bitcoin/utxo"
	"github.com/lugondev/tx-builder/pkg/client"
	"github.com/lugondev/tx-builder/pkg/common"
	"math/rand"
	"testing"
)

func TestBuilderBTCWallet(t *testing.T) {
	//pubkey, err := btcec.ParsePubKey(common2.FromHex("03201d82ca8f8ebe6542459911a5275e8650504050af538d469fda12568c49ed7b"))
	pubkey, err := btcec.ParsePubKey(common2.FromHex("037d418ddabdf94074dd9e11a7f297329f2bf3352d5df480ef3a8ceac87926db93"))
	//pubkey, err := btcec.ParsePubKey(common2.FromHex("03282c5432b0b716e99caf54fafb5a83b41ee5213940b3daf178ad7e3ecd41ae8b"))
	if err != nil {
		t.Fatal(err)
	}
	addressType := common.Taproot

	builder, err := NewTxBtcBuilder(pubkey.SerializeUncompressed(), addressType, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	addresses := chain.PubkeyToAddresses(pubkey, &chaincfg.TestNet3Params)
	t.Log("address legacy:", addresses[addressType])

	c := client.NewClient("https://blockstream.info", "", "", "")
	utxoService := utxo.BlockStreamService{Client: c}
	utxos, err := utxoService.SetAddress(builder.SourceAddressInfo.Address).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("UTXOs: ", utxos.ToUTXOs().Len())

	signedTx, err := builder.SetUtxos(*utxos.ToUTXOs()).
		SetFeeRate(1000).
		SetChangeSource(builder.SourceAddressInfo.Address).
		SetSecretStore(pubkey.SerializeCompressed(), nil).
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

	fmt.Printf("tx: %x", signedTx)
}
