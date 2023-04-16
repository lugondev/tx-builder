package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	authjwt "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/jwt"
	authkey "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/key"
	"github.com/lugondev/tx-builder/src/api/service/controllers"
	"github.com/lugondev/tx-builder/src/infra/postgres/gopg"
	qkmhttp "github.com/lugondev/tx-builder/src/infra/signer-key-manager/http"
	nonclient "github.com/lugondev/tx-builder/src/infra/signer-key-manager/non-client"
	"github.com/lugondev/wallet-signer-manager/pkg/client"
	"net/http"
	"time"
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
	r := mux.NewRouter()
	r = d.BuildRouter(r, "/v1")

	r.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	appPort := d.GetConfig().App.Port
	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%s", appPort),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("server run port: %s\n", appPort)

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
