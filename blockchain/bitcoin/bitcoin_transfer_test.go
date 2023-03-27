package bitcoin

import (
	"fmt"
	"github.com/lugondev/tx-builder/blockchain/bitcoin/utxo"
	"testing"
)

func TestTransfer(t *testing.T) {
	rawTx, err := CreateSegwitTx(
		"cP2gB7hrFoE4AccbB1qyfcgmzDicZ8bkr3XB9GhYzMUEQNkQRRwr",
		"tb1pr375lf8f88dzkxhhecpqarp9w5580eysuycu40czz8s2phd86gss9rwnaf",
		utxo.UnspentTxOutput{
			VOut:   1,
			TxHash: "34287f892662f88f68cadb4b29d51e3dcdd4241eee0f668fd254120316ba2e9c",
		},
		10000)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("raw signed transaction is: ", rawTx)
}
