package chain_test

import (
	"github.com/lugondev/tx-builder/blockchain/bitcoin/chain"
	"testing"
)

func TestChainFeeRate(t *testing.T) {
	feeRate, err := chain.SuggestFeeRate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(feeRate.Low)
	t.Log(feeRate.Average)
	t.Log(feeRate.High)
}
