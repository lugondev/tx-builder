package app

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(hostnameViperKey, hostnameDefault)
	_ = viper.BindEnv(hostnameViperKey, hostnameEnv)

	viper.SetDefault(HttpPortViperKey, httpPortDefault)
	_ = viper.BindEnv(HttpPortViperKey, httpPortEnv)

	viper.SetDefault(accessLogEnabledKey, accessLogEnabledDefault)
	_ = viper.BindEnv(accessLogEnabledKey, accessLogEnabledEnv)

}

const (
	hostnameFlag     = "rest-hostname"
	hostnameViperKey = "rest.hostname"
	hostnameDefault  = ""
	hostnameEnv      = "REST_HOSTNAME"
)

// Hostname register a flag for HTTP server address
func hostname(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hostname to expose REST services
Environment variable: %q`, hostnameEnv)
	f.String(hostnameFlag, hostnameDefault, desc)
	_ = viper.BindPFlag(hostnameViperKey, f.Lookup(hostnameFlag))
}

const (
	httpPortFlag     = "rest-port"
	HttpPortViperKey = "rest.port"
	httpPortDefault  = 8011
	httpPortEnv      = "REST_PORT"
)

// Port register a flag for HTTp server port
func port(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Port to expose REST services
Environment variable: %q`, httpPortEnv)
	f.Uint(httpPortFlag, httpPortDefault, desc)
	_ = viper.BindPFlag(HttpPortViperKey, f.Lookup(httpPortFlag))
}

const (
	accessLogEnabledFlag    = "accesslog-enabled"
	accessLogEnabledKey     = "accesslog.enabled"
	accessLogEnabledDefault = false
	accessLogEnabledEnv     = "ACCESSLOG_ENABLED"
)

func accessLogEnabled(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Enable http accesslog stout
Environment variable: %q`, accessLogEnabledEnv)
	f.Bool(accessLogEnabledFlag, accessLogEnabledDefault, desc)
	_ = viper.BindPFlag(accessLogEnabledKey, f.Lookup(accessLogEnabledFlag))
}

func Flags(f *pflag.FlagSet) {
	hostname(f)
	port(f)
	accessLogEnabled(f)
}
