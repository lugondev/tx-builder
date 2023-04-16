package accounts

import (
	"context"
	"github.com/lugondev/tx-builder/pkg/errors"

	"github.com/lugondev/tx-builder/pkg/toolkit/app/log"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/multitenancy"
	usecases "github.com/lugondev/tx-builder/src/api/business/use-cases"
	"github.com/lugondev/tx-builder/src/api/store"
	"github.com/lugondev/tx-builder/src/entities"
)

const searchAccountsComponent = "use-cases.search-accounts"

type searchAccountsUseCase struct {
	db     store.DB
	logger *log.Logger
}

func NewSearchAccountsUseCase(db store.DB) usecases.SearchAccountsUseCase {
	return &searchAccountsUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(searchAccountsComponent),
	}
}

func (uc *searchAccountsUseCase) Execute(ctx context.Context, filters *entities.AccountFilters, userInfo *multitenancy.UserInfo) ([]*entities.Wallet, error) {
	accs, err := uc.db.Account().Search(ctx, filters, userInfo.AllowedTenants, userInfo.Username)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchAccountsComponent)
	}

	uc.logger.WithContext(ctx).Debug("accounts found successfully")
	return accs, nil
}
