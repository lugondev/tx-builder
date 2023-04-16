package proxy

import (
	"fmt"
	"time"

	"github.com/lugondev/tx-builder/pkg/toolkit/app/multitenancy"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(CacheTTLViperKey, cacheTTLEnv)
	viper.SetDefault(CacheTTLViperKey, cacheDefault)
	_ = viper.BindEnv(MaxIdleConnsPerHostViperKey, maxIdleConnsPerHostEnv)
	viper.SetDefault(MaxIdleConnsPerHostViperKey, maxIdleConnsPerHostDefault)
}

var (
	cacheTTLFlag     = "proxy-cache-ttl"
	CacheTTLViperKey = "proxy.cache.ttl"
	cacheDefault     = 0 * time.Second
	cacheTTLEnv      = "PROXY_CACHE_TTL"

	maxIdleConnsPerHostFlag     = "proxy-max-idle-connections-per-host"
	MaxIdleConnsPerHostViperKey = "proxy.max-idle-connections-per-host"
	maxIdleConnsPerHostDefault  = 50
	maxIdleConnsPerHostEnv      = "PROXY_MAXIDLECONNSPERHOST"
)

type Config struct {
	ProxyCacheTTL *time.Duration
	Multitenancy  bool
}

func Flags(f *pflag.FlagSet) {
	cacheDesc := fmt.Sprintf(`Proxy Cache TTL duration (Disabled by default). Environment variable: %q`, cacheTTLEnv)
	f.Duration(cacheTTLFlag, cacheDefault, cacheDesc)
	_ = viper.BindPFlag(CacheTTLViperKey, f.Lookup(cacheTTLFlag))

	maxIdleConnsPerHostDesc := fmt.Sprintf(`Maximum number of open HTTP connections to a chain proxied. Environment variable: %q`, maxIdleConnsPerHostEnv)
	f.Int(maxIdleConnsPerHostFlag, maxIdleConnsPerHostDefault, maxIdleConnsPerHostDesc)
	_ = viper.BindPFlag(MaxIdleConnsPerHostViperKey, f.Lookup(maxIdleConnsPerHostFlag))
}

func NewConfig() *Config {
	cfg := &Config{
		Multitenancy: viper.GetBool(multitenancy.EnabledViperKey),
	}

	cacheStr := viper.GetDuration(CacheTTLViperKey)
	if cacheStr != time.Duration(0) {
		cfg.ProxyCacheTTL = &cacheStr
	}

	return cfg
}
