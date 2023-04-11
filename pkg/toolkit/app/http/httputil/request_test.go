// +build unit

package httputil

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMethod(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	assert.Equal(t, http.MethodGet, GetMethod(req))
}

func TestGetProtocol(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	assert.Equal(t, "http", GetProtocol(req))

	req, _ = http.NewRequest(http.MethodGet, "http://test.com", nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	assert.Equal(t, "websocket", GetProtocol(req))

	req, _ = http.NewRequest(http.MethodGet, "http://test.com", nil)
	req.Header.Set("Accept", "text/event-stream")
	assert.Equal(t, "sse", GetProtocol(req))
}

func TestTLS(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
	req.TLS = &tls.ConnectionState{
		Version:     tls.VersionTLS11,
		CipherSuite: tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	}
	assert.Equal(t, "1.1", GetTLSVersion(req))
	assert.Equal(t, "TLS_RSA_WITH_3DES_EDE_CBC_SHA", GetTLSCipher(req))
}

func TestToFilters(t *testing.T) {
	v := url.Values{}
	v.Add("key1", "val1")
	v.Add("key2", "val2")
	m := ToFilters(v)

	assert.Equal(t, m["key1"], "val1")
	assert.Equal(t, m["key2"], "val2")
}
