package api

import (
	postgresstore "github.com/lugondev/tx-builder/src/api/store/postgres"
	"github.com/lugondev/tx-builder/src/infra/postgres"
	qkmclient "github.com/lugondev/wallet-signer-manager/pkg/client"

	"github.com/lugondev/tx-builder/pkg/toolkit/app/auth"
	"github.com/lugondev/tx-builder/src/api/business/builder"
	"github.com/lugondev/tx-builder/src/api/service/controllers"
)

func NewAPI(
	cfg *Config,
	db postgres.Client,
	jwt, key auth.Checker,
	keyManagerClient qkmclient.KeyManagerClient,
	qkmStoreID string,
) (*controllers.Builder, error) {

	ucs := builder.NewUseCases(postgresstore.New(db), keyManagerClient)
	return controllers.NewBuilder(cfg.Multitenancy, ucs, keyManagerClient, qkmStoreID, jwt, key), nil
}
