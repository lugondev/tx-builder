package utxo

import (
	"context"
	"github.com/lugondev/tx-builder/pkg/client"
	"testing"
)

func TestCallToMemPoolSpace(t *testing.T) {
	client := client.NewClient("https://mempool.space", "", "", "")
	utxoService := MemPoolSpaceService{c: client}
	utxo, err := utxoService.SetAddress("tb1q5rvwj5fyh02ldstdk77ku0vc3g9utdq693tuet").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utxo.ToUTXOs())
	t.Log(string(utxo.ToUTXOs().ForceToUTXOsJSON()))
}
