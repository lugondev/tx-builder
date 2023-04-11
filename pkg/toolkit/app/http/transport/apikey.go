package transport

import (
	"net/http"

	authutils "github.com/lugondev/tx-builder/pkg/toolkit/app/auth/utils"
	"github.com/lugondev/tx-builder/pkg/toolkit/app/multitenancy"
)

type XAPIKeyHeadersTransport struct {
	apiKey string
	T      http.RoundTripper
}

// NewXAPIKeyHeadersTransport creates a new transport to attach API-KEY as part of request headers
func NewXAPIKeyHeadersTransport(apiKey string) Middleware {
	return func(nxt http.RoundTripper) http.RoundTripper {
		return &XAPIKeyHeadersTransport{
			T:      nxt,
			apiKey: apiKey,
		}
	}
}

func (t *XAPIKeyHeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	userInfo := multitenancy.UserInfoValue(req.Context())
	if t.apiKey == "" {
		return t.T.RoundTrip(req)
	}

	authutils.AddAPIKeyHeaderValue(req, t.apiKey)

	if userInfo != nil {
		if userInfo.TenantID != "" && userInfo.TenantID != multitenancy.WildcardTenant {
			authutils.AddTenantIDHeaderValue(req, userInfo.TenantID)
		}
		if userInfo.Username != "" {
			authutils.AddUsernameHeaderValue(req, userInfo.Username)
		}
	}

	return t.T.RoundTrip(req)
}
