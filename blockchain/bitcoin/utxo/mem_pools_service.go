package utxo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lugondev/tx-builder/pkg/client"
	"github.com/lugondev/tx-builder/pkg/common"
	"net/http"
)

type MemPoolSpaceService struct {
	Client      *client.Client
	addressInfo *common.BTCAddressInfo
}

type MemPoolItemResponse struct {
	TxId   string `json:"txid"`
	VOut   int64  `json:"vout"`
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight int    `json:"block_height"`
		BlockHash   string `json:"block_hash"`
		BlockTime   int    `json:"block_time"`
	} `json:"status"`
	Value int `json:"value"`
}

type MemPoolResponse []*MemPoolItemResponse

// Do send request
func (s *MemPoolSpaceService) Do(ctx context.Context, opts ...client.RequestOption) (res *MemPoolResponse, err error) {
	if s.addressInfo == nil || s.addressInfo.Address == "" {
		return nil, fmt.Errorf("address is empty or invalid")
	}
	r := &client.Request{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("%s/api/address/%s/utxo", s.addressInfo.GetBTCRouterBlockStream(), s.addressInfo.Address),
		SecType:  client.SecTypeNone,
	}

	data, err := s.Client.CallAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}

	return res, err
}

func (b *MemPoolResponse) ToUTXOs() *UnspentTxsOutput {
	txs := make(UnspentTxsOutput, 0)
	for _, tx := range *b {
		txs = append(txs, &UnspentTxOutput{
			TxHash: tx.TxId,
			Value:  int64(tx.Value),
			VOut:   tx.VOut,
		})
	}

	return &txs
}

func (s *MemPoolSpaceService) SetAddress(address string) *MemPoolSpaceService {
	addressInfo := common.GetBTCAddressInfo(address)
	s.addressInfo = addressInfo
	return s
}
