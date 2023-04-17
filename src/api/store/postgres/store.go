package postgres

import (
	"context"

	"github.com/lugondev/tx-builder/src/api/store"
	"github.com/lugondev/tx-builder/src/infra/postgres"
)

type PGStore struct {
	account store.AccountAgent
	address store.AddressAgent
	client  postgres.Client
}

var _ store.DB = &PGStore{}

func New(client postgres.Client) *PGStore {
	return &PGStore{
		account: NewPGAccount(client),
		address: NewPGAddress(client),
		client:  client,
	}
}

func (s *PGStore) Account() store.AccountAgent {
	return s.account
}

func (s *PGStore) Address() store.AddressAgent {
	return s.address
}

func (s *PGStore) RunInTransaction(ctx context.Context, persist func(a store.DB) error) error {
	return s.client.RunInTransaction(ctx, func(dbTx postgres.Client) error {
		return persist(New(dbTx))
	})
}
