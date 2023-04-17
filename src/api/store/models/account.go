package models

import (
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lugondev/tx-builder/src/entities"
)

type Wallet struct {
	tableName struct{} `pg:"wallets"` // nolint:unused,structcheck // reason

	ID                  int
	PublicKey           string
	CompressedPublicKey string
	TenantID            string
	OwnerID             string
	Attributes          map[string]string
	StoreID             string

	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time `pg:"default:now()"`
}

func NewWallet(wallet *entities.Wallet) *Wallet {
	return &Wallet{
		ID:                  wallet.ID,
		PublicKey:           wallet.PublicKey.String(),
		CompressedPublicKey: wallet.CompressedPublicKey.String(),
		TenantID:            wallet.TenantID,
		OwnerID:             wallet.OwnerID,
		StoreID:             wallet.StoreID,
		Attributes:          wallet.Attributes,
		CreatedAt:           wallet.CreatedAt,
		UpdatedAt:           wallet.UpdatedAt,
	}
}

func NewWallets(wallets []*Wallet) []*entities.Wallet {
	var res []*entities.Wallet
	for _, acc := range wallets {
		res = append(res, acc.ToEntity())
	}

	return res
}

func (w *Wallet) ToEntity() *entities.Wallet {
	return &entities.Wallet{
		ID:                  w.ID,
		PublicKey:           hexutil.MustDecode(w.PublicKey),
		CompressedPublicKey: hexutil.MustDecode(w.CompressedPublicKey),
		TenantID:            w.TenantID,
		OwnerID:             w.OwnerID,
		StoreID:             w.StoreID,
		Attributes:          w.Attributes,
		CreatedAt:           w.CreatedAt,
		UpdatedAt:           w.UpdatedAt,
	}
}
