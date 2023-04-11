package flags

import (
	"fmt"
	"github.com/lugondev/tx-builder/src/infra/signer-key-manager/http"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(URLViperKey, urlDefault)
	_ = viper.BindEnv(URLViperKey, urlEnv)

	viper.SetDefault(QKMMetricsURLViperKey, qkmMetricsURLDefault)
	_ = viper.BindEnv(QKMMetricsURLViperKey, qkmMetricsURLEnv)

	viper.SetDefault(storeNameViperKey, storeNameDefault)
	_ = viper.BindEnv(storeNameViperKey, storeNameEnv)

	viper.SetDefault(keyManagerTLSSkipVerifyViperKey, keyManagerTLSSKipVerifyDefault)
	_ = viper.BindEnv(keyManagerTLSSkipVerifyViperKey, keyManagerTLSSkipVerifyEnv)

	viper.SetDefault(keyManagerAPIKeyViperKey, keyManagerAPIKeyDefault)
	_ = viper.BindEnv(keyManagerAPIKeyViperKey, keyManagerAPIKeyEnv)

	viper.SetDefault(keyManagerTLSCertViperKey, keyManagerTLSCertDefault)
	_ = viper.BindEnv(keyManagerTLSCertViperKey, keyManagerTLSCertEnv)

	viper.SetDefault(keyManagerTLSKeyViperKey, keyManagerTLSKeyDefault)
	_ = viper.BindEnv(keyManagerTLSKeyViperKey, keyManagerTLSKeyEnv)
}

func QKMFlags(f *pflag.FlagSet) {
	url(f)
	storeName(f)
	tlsSkipVerify(f)
	authAPIKey(f)
	authTLSCert(f)
	authTLSKey(f)
}

const (
	urlFlag     = "key-manager-url"
	URLViperKey = "key.manager.url"
	urlDefault  = ""
	urlEnv      = "KEY_MANAGER_URL"
)

const (
	QKMMetricsURLViperKey = "key.manager.metrics-url"
	qkmMetricsURLDefault  = ""
	qkmMetricsURLEnv      = "KEY_MANAGER_METRICS_URL"
)

func url(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager HTTP domain.
Environment variable: %q`, urlEnv)
	f.String(urlFlag, urlDefault, desc)
	_ = viper.BindPFlag(URLViperKey, f.Lookup(urlFlag))
}

const (
	storeNameFlag     = "key-manager-store-name"
	storeNameViperKey = "key.manager.store.name"
	storeNameDefault  = ""
	storeNameEnv      = "KEY_MANAGER_STORE_NAME"
)

func storeName(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager ethereum account store name.
Environment variable: %q`, storeNameEnv)
	f.String(storeNameFlag, storeNameDefault, desc)
	_ = viper.BindPFlag(storeNameViperKey, f.Lookup(storeNameFlag))
}

const (
	keyManagerTLSSkipVerifyFlag     = "key-manager-tls-skip-verify"
	keyManagerTLSSkipVerifyViperKey = "key.manager.tls.skip.verify"
	keyManagerTLSSKipVerifyDefault  = false
	keyManagerTLSSkipVerifyEnv      = "KEY_MANAGER_TLS_SKIP_VERIFY"
)

func tlsSkipVerify(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager, disables SSL certificate verification.
Environment variable: %q`, keyManagerTLSSkipVerifyEnv)
	f.Bool(keyManagerTLSSkipVerifyFlag, keyManagerTLSSKipVerifyDefault, desc)
	_ = viper.BindPFlag(keyManagerTLSSkipVerifyViperKey, f.Lookup(keyManagerTLSSkipVerifyFlag))
}

const (
	keyManagerAPIKeyFlag     = "key-manager-api-key"
	keyManagerAPIKeyViperKey = "key.manager.api.key"
	keyManagerAPIKeyDefault  = ""
	keyManagerAPIKeyEnv      = "KEY_MANAGER_API_KEY"
)

func authAPIKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager API-KEY authentication.
Environment variable: %q`, keyManagerAPIKeyEnv)
	f.String(keyManagerAPIKeyFlag, keyManagerAPIKeyDefault, desc)
	_ = viper.BindPFlag(keyManagerAPIKeyViperKey, f.Lookup(keyManagerAPIKeyFlag))
}

const (
	keyManagerTLSCertFlag     = "key-manager-client-tls-cert"
	keyManagerTLSCertViperKey = "key.manager.client.tls.cert"
	keyManagerTLSCertDefault  = ""
	keyManagerTLSCertEnv      = "KEY_MANAGER_CLIENT_TLS_CERT"
)

const (
	keyManagerTLSKeyFlag     = "key-manager-client-tls-key"
	keyManagerTLSKeyViperKey = "key.manager.client.tls.key"
	keyManagerTLSKeyDefault  = ""
	keyManagerTLSKeyEnv      = "KEY_MANAGER_CLIENT_TLS_KEY"
)

func authTLSCert(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager mutual TLS authentication (crt file).
Environment variable: %q`, keyManagerTLSCertEnv)
	f.String(keyManagerTLSCertFlag, keyManagerTLSCertDefault, desc)
	_ = viper.BindPFlag(keyManagerTLSCertViperKey, f.Lookup(keyManagerTLSCertFlag))
}

func authTLSKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager mutual TLS authentication (key file).
Environment variable: %q`, keyManagerTLSKeyEnv)
	f.String(keyManagerTLSKeyFlag, keyManagerTLSKeyDefault, desc)
	_ = viper.BindPFlag(keyManagerTLSKeyViperKey, f.Lookup(keyManagerTLSKeyFlag))
}

func NewQKMConfig(vipr *viper.Viper) *http.Config {
	return &http.Config{
		URL:           vipr.GetString(URLViperKey),
		StoreName:     vipr.GetString(storeNameViperKey),
		APIKey:        vipr.GetString(keyManagerAPIKeyViperKey),
		TLSSkipVerify: vipr.GetBool(keyManagerTLSSkipVerifyViperKey),
		TLSCert:       vipr.GetString(keyManagerTLSCertViperKey),
		TLSKey:        vipr.GetString(keyManagerTLSKeyViperKey),
	}
}
