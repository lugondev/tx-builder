package flags

import (
	"github.com/lugondev/tx-builder/pkg/toolkit/app"
	authjwt "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/jwt/jose"
	authkey "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/key"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/log"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/multitenancy"
	"github.com/lugondev/tx-builder/src/api"
	"github.com/lugondev/tx-builder/src/api/proxy"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func NewAPIFlags(f *pflag.FlagSet) {
	QKMFlags(f)
	PGFlags(f)

	log.Flags(f)
	multitenancy.Flags(f)
	authjwt.Flags(f)
	authkey.Flags(f)
	app.Flags(f)
	proxy.Flags(f)
}

func NewAPIConfig(vipr *viper.Viper) *api.Config {
	return &api.Config{
		Postgres:     NewPGConfig(vipr),
		Multitenancy: vipr.GetBool(multitenancy.EnabledViperKey),
		Proxy:        proxy.NewConfig(),
		QKM:          NewQKMConfig(vipr),
		App: &api.AppConfig{
			Port: vipr.GetString(app.HttpPortViperKey),
		},
	}
}
