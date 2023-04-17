package types

type SignatureResponse struct {
	Signature string              `json:"signature"`
	SignData  *SignMessageRequest `json:"sign_data"`
	Pubkey    string              `json:"pubkey"`
	Address   string              `json:"address"`
}
