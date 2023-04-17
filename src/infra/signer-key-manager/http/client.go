package http

import (
	"crypto/tls"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/http"
	qkm "github.com/lugondev/wallet-signer-manager/pkg/client"
)

func New(cfg *Config) (*qkm.HTTPClient, error) {
	httpCfg := http.NewDefaultConfig()
	// Support user's JWT forwarding
	httpCfg.AuthHeaderForward = true

	httpCfg.InsecureSkipVerify = cfg.TLSSkipVerify
	if cfg.APIKey != "" {
		httpCfg.AuthHeaderForward = false
		httpCfg.Authorization = "Basic " + cfg.APIKey
	}

	if cfg.TLSCert != "" && cfg.TLSKey != "" {
		cert, err := tls.LoadX509KeyPair(cfg.TLSCert, cfg.TLSKey)
		if err != nil {
			return nil, err
		}
		httpCfg.ClientCert = &cert
	}

	return qkm.NewHTTPClient(http.NewClient(httpCfg), &qkm.Config{
		URL: cfg.URL,
	}), nil
}
