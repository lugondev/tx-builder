package builder

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
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
	pubkey, err := btcec.ParsePubKey(common2.FromHex("02f564c5d9f932acbb0c81438f0e4389509f87383e22d4f203e0bb09c33135e86a"))
	if err != nil {
		t.Fatal(err)
	}
	addressType := common.Segwit
	netParams := &chaincfg.MainNetParams

	builder, err := NewTxBtcBuilder(pubkey.SerializeUncompressed(), addressType, netParams)
	if err != nil {
		t.Fatal(err)
	}
	addresses := chain.PubkeyToAddresses(pubkey, netParams)
	t.Log("address:", addresses[addressType])

	c := client.NewClient("https://blockstream.info", "", "", "")
	utxoService := utxo.BlockStreamService{Client: c}
	utxos, err := utxoService.SetAddress(builder.SourceAddressInfo.Address).
		Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("UTXOs: ", utxos.ToUTXOs().Len())
	wif, _ := btcutil.DecodeWIF("cVacJiScoPMAugWKRwMU2HVUPE4PhcJLgxVCexieWEWcTiYC8bSn")

	signedTx, err := builder.SetUtxos(*utxos.ToUTXOs()).
		SetPrivKey(wif.PrivKey).
		//SetSecretStore(pubkey.SerializeCompressed(), nil).
		SetFeeRate(1000).
		SetChangeSource(builder.SourceAddressInfo.Address).
		SetOutputs([]*Output{
			{
				Address: toAddress,
				Amount:  rand.Int63n(200) + 100,
			},
		}).
		Build()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("tx: %x", signedTx)
}
