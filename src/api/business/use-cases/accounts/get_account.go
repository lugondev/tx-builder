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

const getAccountComponent = "use-cases.get-account"

type getAccountUseCase struct {
	db     store.DB
	logger *log.Logger
}

func NewGetAccountUseCase(db store.DB) usecases.GetAccountUseCase {
	return &getAccountUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(getAccountComponent),
	}
}

func (uc *getAccountUseCase) Execute(ctx context.Context, pubkey string, userInfo *multitenancy.UserInfo) (*entities.Wallet, error) {
	ctx = log.WithFields(ctx, log.Field("pubkey", pubkey))
	logger := uc.logger.WithContext(ctx)

	acc, err := uc.db.Account().FindOneByPubkey(ctx, pubkey, userInfo.AllowedTenants, userInfo.Username)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getAccountComponent)
	}

	logger.Debug("account found successfully")
	return acc, nil
}
