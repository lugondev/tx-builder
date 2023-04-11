package utxo

import (
	"context"
	"github.com/lugondev/tx-builder/pkg/client"
	"testing"
)

func TestCallToBtcCom(t *testing.T) {
	client := client.NewClient("https://chain.api.btc.com", "", "", "")
	utxoService := BTCComService{Client: client}
	utxo, err := utxoService.SetAddress("3QS5z2ei7sPTUmonW88ZZAfjWXYzVtFsBF").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utxo.ToUTXOs())
	t.Log(string(utxo.ToUTXOs().ForceToUTXOsJSON()))
}
