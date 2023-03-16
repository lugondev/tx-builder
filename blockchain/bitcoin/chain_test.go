package bitcoin_test

import (
	"github.com/lugondev/tx-builder/blockchain/bitcoin"
	"testing"
)

func TestChainFeeRate(t *testing.T) {
	feeRate, err := bitcoin.SuggestFeeRate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(feeRate.Low)
	t.Log(feeRate.Average)
	t.Log(feeRate.High)
}
