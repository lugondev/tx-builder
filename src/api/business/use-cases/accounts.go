package usecases

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/multitenancy"
	"github.com/lugondev/tx-builder/src/entities"
)

//go:generate mockgen -source=accounts.go -destination=mocks/accounts.go -package=mocks

type AccountUseCases interface {
	Get() GetAccountUseCase
	Create() CreateAccountUseCase
	Update() UpdateAccountUseCase
	Search() SearchAccountsUseCase
}

type GetAccountUseCase interface {
	Execute(ctx context.Context, pubkey string, userInfo *multitenancy.UserInfo) (*entities.Account, error)
}

type CreateAccountUseCase interface {
	Execute(ctx context.Context, identity *entities.Account, privateKey hexutil.Bytes, chainName string, userInfo *multitenancy.UserInfo) (*entities.Account, error)
}

type SearchAccountsUseCase interface {
	Execute(ctx context.Context, filters *entities.AccountFilters, userInfo *multitenancy.UserInfo) ([]*entities.Account, error)
}

type UpdateAccountUseCase interface {
	Execute(ctx context.Context, identity *entities.Account, userInfo *multitenancy.UserInfo) (*entities.Account, error)
}
