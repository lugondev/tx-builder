package utxo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lugondev/tx-builder/pkg/client"
	"github.com/lugondev/tx-builder/pkg/common"
	"net/http"
)

type BlockChainInfoService struct {
	c       *client.Client
	address string
}

type BlockchainInfoResponse struct {
	Notice         string `json:"notice"`
	UnspentOutputs []struct {
		TxHashBigEndian string `json:"tx_hash_big_endian"`
		TxHash          string `json:"tx_hash"`
		TxOutputN       int64  `json:"tx_output_n"`
		Script          string `json:"script"`
		Value           int64  `json:"value"`
		ValueHex        string `json:"value_hex"`
		Confirmations   int64  `json:"confirmations"`
		TxIndex         int64  `json:"tx_index"`
	} `json:"unspent_outputs"`
}

// Do send request
func (s *BlockChainInfoService) Do(ctx context.Context, opts ...client.RequestOption) (res *BlockchainInfoResponse, err error) {
	r := &client.Request{
		Method:   http.MethodGet,
		Endpoint: "/unspent",
		SecType:  client.SecTypeNone,
	}
	if s.address != "" {
		r.SetParam("active", s.address)
	} else {
		return nil, fmt.Errorf("address is empty or invalid")
	}

	data, err := s.c.CallAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (b *BlockchainInfoResponse) ToUTXOs() *UnspentTxsOutput {
	txs := make(UnspentTxsOutput, 0)
	for _, tx := range b.UnspentOutputs {
		txs = append(txs, &UnspentTxOutput{
			TxHash:        tx.TxHash,
			Value:         tx.Value,
			VOut:          tx.TxOutputN,
			Confirmations: &tx.Confirmations,
		})
	}

	return &txs
}

func (s *BlockChainInfoService) SetAddress(address string) *BlockChainInfoService {
	addressInfo := common.GetBTCAddressType(address)
	if addressInfo == nil || addressInfo.Chain != common.BTCMainnet {
		s.address = ""
	} else {
		s.address = address
	}

	return s
}
