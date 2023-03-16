package evm

import (
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type TxRequest struct {
	Nonce    *big.Int        `json:"nonce"`
	GasPrice *big.Int        `json:"gasPrice"`
	GasLimit uint64          `json:"gasLimit"`
	Value    *big.Int        `json:"value"`
	To       *common.Address `json:"to"`
	From     common.Address  `json:"from"`
	Data     []byte          `json:"data"`
}

func (t *TxRequest) PrepareTransaction(client *Client) (*types.Transaction, error) {
	if t.From == common.BytesToAddress([]byte{0}) {
		return nil, errors.New("from address is required")
	}
	if t.To == nil {
		return nil, errors.New("to address is required")
	}

	var err error
	if t.GasLimit == 0 {
		if t.GasLimit, err = client.EthClient.EstimateGas(client.Ctx, ethereum.CallMsg{
			From:      t.From,
			To:        t.To,
			GasPrice:  t.GasPrice,
			GasFeeCap: t.GasPrice,
			GasTipCap: t.GasPrice,
			Value:     t.Value,
			Data:      t.Data,
		}); err != nil {
			return nil, err
		}
	}
	if t.Nonce == nil {
		if t.Nonce, err = client.AccountNonce(t.From); err != nil {
			return nil, err
		}
	}
	if t.GasPrice == nil {
		if t.GasPrice, err = client.EthClient.SuggestGasPrice(client.Ctx); err != nil {
			return nil, err
		}
	}

	return types.NewTx(&types.LegacyTx{
		Nonce:    t.Nonce.Uint64(),
		GasPrice: t.GasPrice,
		Gas:      t.GasLimit,
		To:       t.To,
		Value:    t.Value,
		Data:     t.Data,
	}), nil
}
