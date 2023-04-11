package controllers

import (
	"context"
	"net/http"

	qkm "github.com/lugondev/wallet-signer-manager/pkg/client"

	"github.com/gorilla/mux"
	usecases "github.com/lugondev/tx-builder/src/api/business/use-cases"
)

type Builder struct {
	accountsCtrl *AccountsController
}

func NewBuilder(ucs usecases.UseCases, keyManagerClient qkm.KeyManagerClient, qkmStoreID string) *Builder {
	return &Builder{
		accountsCtrl: NewAccountsController(ucs, keyManagerClient, qkmStoreID),
	}
}

func (b *Builder) Build(_ context.Context, _ string, _ func(response *http.Response) error) (http.Handler, error) {

	router := mux.NewRouter()
	b.accountsCtrl.Append(router)

	return router, nil
}
