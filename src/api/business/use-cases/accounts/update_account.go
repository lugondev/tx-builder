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

const updateAccountComponent = "use-cases.update-account"

type updateAccountUseCase struct {
	db     store.DB
	logger *log.Logger
}

func NewUpdateAccountUseCase(db store.DB) usecases.UpdateAccountUseCase {
	return &updateAccountUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(updateAccountComponent),
	}
}

func (uc *updateAccountUseCase) Execute(ctx context.Context, acc *entities.Wallet, userInfo *multitenancy.UserInfo) (*entities.Wallet, error) {
	ctx = log.WithFields(ctx, log.Field("pubkey", acc.PublicKey))
	logger := uc.logger.WithContext(ctx)

	curAcc, err := uc.db.Account().FindOneByPubkey(ctx, acc.PublicKey.String(), userInfo.AllowedTenants, userInfo.Username)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateAccountComponent)
	}

	if acc.Attributes != nil {
		curAcc.Attributes = acc.Attributes
	}
	if acc.StoreID != "" {
		curAcc.StoreID = acc.StoreID
	}

	updatedAcc, err := uc.db.Account().Update(ctx, curAcc)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateAccountComponent)
	}

	logger.Info("account updated successfully")
	return updatedAcc, nil
}
