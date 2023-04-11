package chain

import (
	"encoding/json"
	"fmt"
	"github.com/lugondev/tx-builder/pkg/httpUtil"
	"net/http"
)

type FeeRate struct {
	Low     int64
	Average int64
	High    int64
}

func SuggestFeeRate() (*FeeRate, error) {
	url := "https://mempool.space/api/v1/fees/recommended"

	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}
	respDict := make(map[string]interface{})
	err = json.Unmarshal(response.Body, &respDict)
	if err != nil {
		return nil, err
	}

	var low, avg, high float64
	var ok bool
	if low, ok = respDict["minimumFee"].(float64); !ok {
		low = 1
	}
	if avg, ok = respDict["halfHourFee"].(float64); !ok {
		avg = low
	}
	if high, ok = respDict["fastestFee"].(float64); !ok {
		high = avg
	}
	return &FeeRate{
		Low:     int64(low),
		Average: int64(avg),
		High:    int64(high),
	}, nil
}
