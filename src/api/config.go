package api

import (
	"github.com/lugondev/tx-builder/src/api/proxy"
	"github.com/lugondev/tx-builder/src/infra/postgres/gopg"
	quorumkeymanager "github.com/lugondev/tx-builder/src/infra/signer-key-manager/http"
)

type AppConfig struct {
	Port string
}

type Config struct {
	Postgres     *gopg.Config
	Multitenancy bool
	Proxy        *proxy.Config
	QKM          *quorumkeymanager.Config
	App          *AppConfig
}
