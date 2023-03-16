package evm

import (
	"math/big"
)

// TxBuilder represents a transaction builder that builds transactions
type TxBuilder struct {
	*TxRequest
	ChainID *big.Int
}

// NewTxBuilder creates a new transaction builder.
func NewTxBuilder(chainID *big.Int) *TxBuilder {
	return &TxBuilder{
		&TxRequest{},
		chainID,
	}
}
