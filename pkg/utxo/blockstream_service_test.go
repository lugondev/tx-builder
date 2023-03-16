package utxo

import (
	"context"
	"github.com/lugondev/tx-builder/pkg/client"
	"testing"
)

func TestCallToBlockStream(t *testing.T) {
	client := client.NewClient("https://blockstream.info", "", "", "")
	utxoService := BlockStreamService{c: client}
	utxo, err := utxoService.SetAddress("tb1pr375lf8f88dzkxhhecpqarp9w5580eysuycu40czz8s2phd86gss9rwnaf").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utxo.ToUTXOs())
	t.Log(string(utxo.ToUTXOs().ForceToUTXOsJSON()))
}
