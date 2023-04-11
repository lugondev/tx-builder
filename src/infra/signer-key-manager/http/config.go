package http

type Config struct {
	URL           string
	StoreName     string
	APIKey        string
	TLSSkipVerify bool
	TLSCert       string
	TLSKey        string
}
