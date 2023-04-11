package api

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/lugondev/tx-builder/src/api/service/controllers"
	"net/http"
	"time"

	authjwt "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/jwt"
	authkey "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/key"
	"github.com/lugondev/tx-builder/src/infra/postgres/gopg"
	qkmhttp "github.com/lugondev/tx-builder/src/infra/signer-key-manager/http"
	nonclient "github.com/lugondev/tx-builder/src/infra/signer-key-manager/non-client"
	"github.com/lugondev/wallet-signer-manager/pkg/client"
)

type Daemon struct {
	*controllers.Builder
	config *Config
}

func New(ctx context.Context, cfg *Config) (*Daemon, error) {
	// Initialize infra dependencies
	qkmClient, err := QKMClient(cfg)
	if err != nil {
		return nil, err
	}

	postgresClient, err := gopg.New("orchestrate.api", cfg.Postgres)
	if err != nil {
		return nil, err
	}

	authjwt.Init(ctx)
	authkey.Init(ctx)

	api, err := NewAPI(
		cfg,
		postgresClient,
		authjwt.GlobalChecker(),
		authkey.GlobalChecker(),
		qkmClient,
		cfg.QKM.StoreName,
	)

	if err != nil {
		return nil, err
	}

	return &Daemon{api, cfg}, nil
}

func (d *Daemon) Run(ctx context.Context) error {
	build, err := d.Build(ctx, "orchestrate.api", func(response *http.Response) error {
		return nil
	})
	if err != nil {
		return err
	}
	r := mux.NewRouter()

	r.Handle("/v1", build)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8001",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv.ListenAndServe()
}

func (d *Daemon) GetConfig() *Config {
	return d.config
}

func QKMClient(cfg *Config) (client.KeyManagerClient, error) {
	if cfg.QKM.URL != "" {
		return qkmhttp.New(cfg.QKM)
	}

	return nonclient.NewNonClient(), nil
}
