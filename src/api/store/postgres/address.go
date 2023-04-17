package postgres

import (
	"context"
	"time"

	"github.com/lugondev/tx-builder/src/infra/postgres"

	"github.com/lugondev/tx-builder/pkg/errors"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/log"
	"github.com/lugondev/tx-builder/src/api/store"
	"github.com/lugondev/tx-builder/src/api/store/models"
	"github.com/lugondev/tx-builder/src/entities"
)

type PGAccount struct {
	client postgres.Client
	logger *log.Logger
}

var _ store.AccountAgent = &PGAccount{}

func NewPGAccount(client postgres.Client) *PGAccount {
	return &PGAccount{
		client: client,
		logger: log.NewLogger().SetComponent("data-agents.account"),
	}
}

func (agent *PGAccount) Insert(ctx context.Context, account *entities.Wallet) (*entities.Wallet, error) {
	model := models.NewWallet(account)
	model.CreatedAt = time.Now().UTC()
	model.UpdatedAt = time.Now().UTC()

	err := agent.client.ModelContext(ctx, model).Insert()
	if err != nil {
		errMsg := "failed to insert account"
		agent.logger.WithContext(ctx).WithError(err).Error(errMsg)
		return nil, errors.FromError(err).SetMessage(errMsg)
	}

	return model.ToEntity(), nil
}

func (agent *PGAccount) Update(ctx context.Context, account *entities.Wallet) (*entities.Wallet, error) {
	model := models.NewWallet(account)
	model.UpdatedAt = time.Now().UTC()

	q := agent.client.ModelContext(ctx, model).
		Where("compressed_public_key = ?", account.CompressedPublicKey.String()).
		Where("tenant_id = ?", account.TenantID)

	if account.OwnerID != "" {
		q = q.Where("owner_id = ?", account.OwnerID)
	}

	err := q.UpdateNotZero()
	if err != nil {
		errMsg := "failed to update account"
		agent.logger.WithContext(ctx).WithError(err).Error(errMsg)
		return nil, errors.FromError(err).SetMessage(errMsg)
	}

	return model.ToEntity(), nil
}

func (agent *PGAccount) Search(ctx context.Context, filters *entities.AccountFilters, tenants []string, ownerID string) ([]*entities.Wallet, error) {
	var accounts []*models.Wallet

	q := agent.client.ModelContext(ctx, &accounts)
	if filters.TenantID != "" {
		q = q.Where("tenant_id = ?", filters.TenantID)
	}

	err := q.WhereAllowedTenants("", tenants).WhereAllowedOwner("", ownerID).Order("id ASC").Select()
	if err != nil && !errors.IsNotFoundError(err) {
		errMsg := "failed to search accounts"
		agent.logger.WithContext(ctx).WithError(err).Error(errMsg)
		return nil, errors.FromError(err).SetMessage(errMsg)
	}

	return models.NewWallets(accounts), nil
}

func (agent *PGAccount) FindOneByPubkey(ctx context.Context, pubkey string, tenants []string, ownerID string) (*entities.Wallet, error) {
	account := &models.Wallet{}

	err := agent.client.
		ModelContext(ctx, account).
		Where("public_key = ?", pubkey).
		WhereAllowedTenants("", tenants).
		WhereAllowedOwner("", ownerID).
		SelectOne()
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, errors.FromError(err).SetMessage("account not found")
		}

		errMsg := "failed to find one account by pubkey"
		agent.logger.WithContext(ctx).WithError(err).Error(errMsg)
		return nil, errors.FromError(err).SetMessage(errMsg)
	}

	return account.ToEntity(), nil
}
