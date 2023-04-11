package builder

import (
	usecases "github.com/lugondev/tx-builder/src/api/business/use-cases"
	"github.com/lugondev/tx-builder/src/api/store"
	qkmclient "github.com/lugondev/wallet-signer-manager/pkg/client"
)

type useCases struct {
	*accountUseCases
}

func NewUseCases(
	db store.DB,
	keyManagerClient qkmclient.KeyManagerClient,
) usecases.UseCases {

	accountUseCases := newAccountUseCases(db, keyManagerClient)

	return &useCases{
		accountUseCases: accountUseCases,
	}
}
