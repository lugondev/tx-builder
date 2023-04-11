package store

import (
	"context"

	"github.com/lugondev/tx-builder/src/entities"
)

//go:generate mockgen -source=store.go -destination=mocks/mock.go -package=mocks

type DB interface {
	Account() AccountAgent
	RunInTransaction(ctx context.Context, persistFunc func(db DB) error) error
}

type AccountAgent interface {
	Insert(ctx context.Context, account *entities.Account) (*entities.Account, error)
	Update(ctx context.Context, account *entities.Account) (*entities.Account, error)
	FindOneByAddress(ctx context.Context, address string, tenants []string, ownerID string) (*entities.Account, error)
	Search(ctx context.Context, filters *entities.AccountFilters, tenants []string, ownerID string) ([]*entities.Account, error)
}
