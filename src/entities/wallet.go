package entities

import (
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Wallet struct {
	ID                  int
	PublicKey           hexutil.Bytes
	CompressedPublicKey hexutil.Bytes
	TenantID            string
	OwnerID             string
	StoreID             string
	Attributes          map[string]string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
