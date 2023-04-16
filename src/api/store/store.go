package store

import (
	"context"

	"github.com/lugondev/tx-builder/src/entities"
)

type DB interface {
	Account() AccountAgent
	RunInTransaction(ctx context.Context, persistFunc func(db DB) error) error
}

type AccountAgent interface {
	Insert(ctx context.Context, account *entities.Wallet) (*entities.Wallet, error)
	Update(ctx context.Context, account *entities.Wallet) (*entities.Wallet, error)
	FindOneByPubkey(ctx context.Context, pubkey string, tenants []string, ownerID string) (*entities.Wallet, error)
	Search(ctx context.Context, filters *entities.AccountFilters, tenants []string, ownerID string) ([]*entities.Wallet, error)
}
