package models

import (
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lugondev/tx-builder/src/entities"
)

type Account struct {
	tableName struct{} `pg:"accounts"` // nolint:unused,structcheck // reason

	ID                  int
	Alias               string
	Address             string
	PublicKey           string
	CompressedPublicKey string
	TenantID            string
	OwnerID             string
	Attributes          map[string]string
	// TODO add internal labels to store accountID
	StoreID string

	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time `pg:"default:now()"`
}

func NewAccount(account *entities.Account) *Account {
	return &Account{
		Alias:               account.Alias,
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

func NewAccounts(accounts []*Account) []*entities.Account {
	res := []*entities.Account{}
	for _, acc := range accounts {
		res = append(res, acc.ToEntity())
	}

	return res
}

func (acc *Account) ToEntity() *entities.Account {
	return &entities.Account{
		Alias:               acc.Alias,
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
