package utxo

import "encoding/json"

type UnspentTxOutput struct {
	TxHash        string `json:"txHash"`
	Value         int64  `json:"value"`
	VOut          int64  `json:"vOut"`
	Confirmations *int64 `json:"confirmations"`
}

type UnspentTxsOutput []*UnspentTxOutput

func (u *UnspentTxsOutput) ToUTXOsJSON() ([]byte, error) {
	return json.Marshal(u)
}
func (u *UnspentTxsOutput) Len() int {
	return len(*u)
}

func (u *UnspentTxsOutput) ForceToUTXOsJSON() []byte {
	bytes, err := json.Marshal(u)
	if err != nil {
		return []byte{}
	}
	return bytes
}
