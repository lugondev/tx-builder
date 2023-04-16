package entities

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lugondev/tx-builder/pkg/common"
)

type BTCWallet struct {
	Address             string                `json:"address"`
	AddressType         common.BTCAddressType `json:"addressType"`
	PublicKey           hexutil.Bytes         `json:"publicKey"`
	CompressedPublicKey hexutil.Bytes         `json:"compressedPublicKey"`
	Namespace           string                `json:"namespace,omitempty" example:"tenant_id"`
}
