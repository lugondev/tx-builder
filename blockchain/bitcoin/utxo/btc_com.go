package utxo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lugondev/tx-builder/pkg/client"
	"github.com/lugondev/tx-builder/pkg/common"
	"net/http"
)

type BTCComService struct {
	Client  *client.Client
	address string
}

type BTCComResponse struct {
	Data struct {
		List []struct {
			TxHash        string `json:"tx_hash"`
			TxOutputN     int64  `json:"tx_output_n"`
			TxOutputN2    int64  `json:"tx_output_n2"`
			Value         int64  `json:"value"`
			Confirmations int64  `json:"confirmations"`
		} `json:"list"`
		Page       int `json:"page"`
		PageTotal  int `json:"page_total"`
		Pagesize   int `json:"pagesize"`
		TotalCount int `json:"total_count"`
	} `json:"data"`
	ErrCode int    `json:"err_code"`
	ErrNo   int    `json:"err_no"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Do send request
func (s *BTCComService) Do(ctx context.Context, opts ...client.RequestOption) (res *BTCComResponse, err error) {
	if s.address == "" {
		return nil, fmt.Errorf("address is empty or invalid")
	}
	r := &client.Request{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/v3/address/%s/unspent", s.address),
		SecType:  client.SecTypeNone,
	}

	data, err := s.Client.CallAPI(ctx, r, opts...)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (b *BTCComResponse) ToUTXOs() *UnspentTxsOutput {
	txs := make(UnspentTxsOutput, 0)
	for _, tx := range b.Data.List {
		txs = append(txs, &UnspentTxOutput{
			TxHash:        tx.TxHash,
			Value:         tx.Value,
			VOut:          tx.TxOutputN,
			Confirmations: &tx.Confirmations,
		})
	}

	return &txs
}

func (s *BTCComService) SetAddress(address string) *BTCComService {
	addressInfo := common.GetBTCAddressInfo(address)
	if addressInfo == nil || addressInfo.Chain != common.BTCMainnet {
		s.address = ""
	} else {
		s.address = address
	}

	return s
}
