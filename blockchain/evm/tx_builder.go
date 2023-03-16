package evm

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

// TxBuilder represents a transaction builder that builds transactions
type TxBuilder struct {
	tx  *TxRequest
	ctx context.Context
}

// NewTxBuilder creates a new transaction builder.
func NewTxBuilder(ctx context.Context) *TxBuilder {
	return &TxBuilder{
		&TxRequest{},
		ctx,
	}
}

// SetFrom sets the from address of the transaction.
func (b *TxBuilder) SetFrom(from common.Address) *TxBuilder {
	if from == common.BytesToAddress([]byte{}) {
		panic("from address is not set")
	}
	b.tx.From = from
	return b
}

// SetTo sets the to address of the transaction.
func (b *TxBuilder) SetTo(to common.Address) *TxBuilder {
	if to == common.BytesToAddress([]byte{}) {
		panic("from address is not set")
	}
	b.tx.To = &to
	return b
}

// SetNonce sets the nonce of the transaction.
func (b *TxBuilder) SetNonce(nonce uint64) *TxBuilder {
	b.tx.Nonce = new(big.Int).SetUint64(nonce)
	return b
}

// SetGasPrice sets the gas price of the transaction.
func (b *TxBuilder) SetGasPrice(gasPrice *big.Int) *TxBuilder {
	// check gas price greater than 0
	if gasPrice == nil || gasPrice.Cmp(big.NewInt(0)) <= 0 {
		panic("gas price is invalid")
	}
	b.tx.GasPrice = gasPrice
	return b
}

// SetGasLimit sets the gas limit of the transaction.
func (b *TxBuilder) SetGasLimit(gasLimit uint64) *TxBuilder {
	// check gas limit greater than 0
	if gasLimit == 0 {
		panic("gas limit is invalid")
	}
	b.tx.GasLimit = gasLimit
	return b
}

// SetValue sets the value of the transaction.
func (b *TxBuilder) SetValue(value *big.Int) *TxBuilder {
	b.tx.Value = value
	return b
}

// SetData sets the data of the transaction.
func (b *TxBuilder) SetData(data []byte) *TxBuilder {
	b.tx.Data = data
	return b
}

// PrepareTransferToken builds the transaction to transfer token.
func (b *TxBuilder) PrepareTransferToken(token common.Address, amount *big.Int) *TxBuilder {
	// check token address
	if token == common.BytesToAddress([]byte{}) {
		panic("token address is not set")
	}
	// check amount greater than 0
	if amount == nil || amount.Cmp(big.NewInt(0)) <= 0 {
		panic("amount is invalid")
	}

	methodID, _ := GetMethodID("transfer(address,uint256)")
	paddedAddress := common.LeftPadBytes(b.tx.To.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	b.SetTo(token)
	b.SetData(append(methodID, append(paddedAddress, paddedAmount...)...))

	return b
}

// PrepareApproveToken builds the transaction to approve token.
func (b *TxBuilder) PrepareApproveToken(token common.Address, amount *big.Int) *TxBuilder {
	// check token address
	if token == common.BytesToAddress([]byte{}) {
		panic("token address is not set")
	}
	// check amount greater than 0
	if amount == nil || amount.Cmp(big.NewInt(0)) <= 0 {
		panic("amount is invalid")
	}

	methodID, _ := GetMethodID("approve(address,uint256)")
	paddedAddress := common.LeftPadBytes(b.tx.To.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	b.SetTo(token)
	b.SetData(append(methodID, append(paddedAddress, paddedAmount...)...))

	return b
}

// PrepareRevokeToken builds the transaction to revoke token.
func (b *TxBuilder) PrepareRevokeToken(token common.Address) *TxBuilder {
	return b.PrepareApproveToken(token, big.NewInt(0))
}

// Build builder the transaction.
func (b *TxBuilder) Build(client *Client) (*types.Transaction, error) {
	return b.tx.PrepareTransaction(client)
}

// GetTxRequest returns the transaction request.
func (b *TxBuilder) GetTxRequest() *TxRequest {
	return b.tx
}
