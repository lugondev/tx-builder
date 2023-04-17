package formatters

import "github.com/lugondev/tx-builder/src/api/service/types"

func FormatSignatureResponse(signData *types.SignMessageRequest, signature, pubkey string) *types.SignatureResponse {
	return &types.SignatureResponse{
		Signature: signature,
		SignData:  signData,
		Pubkey:    pubkey,
		Address:   "",
	}
}
