package models

import (
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lugondev/tx-builder/src/entities"
)

type Account struct {
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

func NewAccount(account *entities.Wallet) *Account {
	return &Account{
		PublicKey:           account.PublicKey.String(),
		CompressedPublicKey: account.CompressedPublicKey.String(),
		TenantID:            account.TenantID,
		OwnerID:             account.OwnerID,
		StoreID:             account.StoreID,
		Attributes:          account.Attributes,
		CreatedAt:           account.CreatedAt,
		UpdatedAt:           account.UpdatedAt,
	}
}

func NewAccounts(accounts []*Account) []*entities.Wallet {
	var res []*entities.Wallet
	for _, acc := range accounts {
		res = append(res, acc.ToEntity())
	}

	return res
}

func (acc *Account) ToEntity() *entities.Wallet {
	return &entities.Wallet{
		PublicKey:           hexutil.MustDecode(acc.PublicKey),
		CompressedPublicKey: hexutil.MustDecode(acc.CompressedPublicKey),
		TenantID:            acc.TenantID,
		OwnerID:             acc.OwnerID,
		StoreID:             acc.StoreID,
		Attributes:          acc.Attributes,
		CreatedAt:           acc.CreatedAt,
		UpdatedAt:           acc.UpdatedAt,
	}
}
