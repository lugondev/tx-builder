package builder

import (
	usecases "github.com/lugondev/tx-builder/src/api/business/use-cases"
	"github.com/lugondev/tx-builder/src/api/business/use-cases/accounts"
	"github.com/lugondev/tx-builder/src/api/store"
	qkmclient "github.com/lugondev/wallet-signer-manager/pkg/client"
)

type accountUseCases struct {
	create usecases.CreateAccountUseCase
	get    usecases.GetAccountUseCase
	search usecases.SearchAccountsUseCase
	update usecases.UpdateAccountUseCase
}

func newAccountUseCases(
	db store.DB,
	keyManagerClient qkmclient.KeyManagerClient,
) *accountUseCases {
	searchAccountsUC := accounts.NewSearchAccountsUseCase(db)

	return &accountUseCases{
		create: accounts.NewCreateAccountUseCase(db, searchAccountsUC, keyManagerClient),
		get:    accounts.NewGetAccountUseCase(db),
		search: searchAccountsUC,
		update: accounts.NewUpdateAccountUseCase(db),
	}
}

func (u *accountUseCases) Get() usecases.GetAccountUseCase {
	return u.get
}

func (u *accountUseCases) Search() usecases.SearchAccountsUseCase {
	return u.search
}

func (u *accountUseCases) Create() usecases.CreateAccountUseCase {
	return u.create
}

func (u *accountUseCases) Update() usecases.UpdateAccountUseCase {
	return u.update
}
