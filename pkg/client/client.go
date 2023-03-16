package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"github.com/lugondev/tx-builder/pkg/common"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// Config
const (
	UseAgent  = "RawTx-Bitcoin/golang"
	LogPrefix = "RawTx-golang: "
)

const (
	TimestampKey  = "timestamp"
	SignatureKey  = "signature"
	RecvWindowKey = "recvWindow"
)

// Redefining the standard package
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func currentTimestamp() int64 {
	return FormatTimestamp(time.Now())
}

// FormatTimestamp formats a time into Unix timestamp in milliseconds, as requested by Binance.
func FormatTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// NewClient initialize an API client instance with API key and secret key.
// You should always call this function before using this SDK.
// Services will be created by the form client.NewXXXService().
func NewClient(baseURL, apiKey, secretKey, headerKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		APIHeader:  headerKey,
		BaseURL:    baseURL,
		UserAgent:  UseAgent,
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, LogPrefix, log.LstdFlags),
		IsDebug:    true,
	}
}

// NewProxyClient passing a proxy url
func NewProxyClient(baseURL, apiKey, secretKey, headerKey, proxyUrl string) *Client {
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		log.Fatal(err)
	}
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		APIHeader: headerKey,
		BaseURL:   baseURL,
		UserAgent: UseAgent,
		HTTPClient: &http.Client{
			Transport: tr,
		},
		Logger: log.New(os.Stderr, LogPrefix, log.LstdFlags),
	}
}

type doFunc func(req *http.Request) (*http.Response, error)

// Client define API client
type Client struct {
	APIKey     string
	SecretKey  string
	APIHeader  string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	IsDebug    bool
	Logger     *log.Logger
	TimeOffset int64
	do         doFunc
}

func (c *Client) Debug(format string, v ...interface{}) {
	if c.IsDebug {
		c.Logger.Printf(format, v...)
	}
}

func (c *Client) ParseRequest(r *Request, opts ...RequestOption) (err error) {
	// set Request options from user
	for _, opt := range opts {
		opt(r)
	}
	err = r.validate()
	if err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.Endpoint)
	if r.recvWindow > 0 {
		r.SetParam(RecvWindowKey, r.recvWindow)
	}
	if r.SecType == SecTypeSigned {
		r.SetParam(TimestampKey, currentTimestamp()-c.TimeOffset)
	}
	queryString := r.query.Encode()
	body := &bytes.Buffer{}
	bodyString := r.form.Encode()
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}
	if r.SecType == SecTypeAPIKey || r.SecType == SecTypeSigned {
		header.Set(c.APIHeader, c.APIKey)
	}

	if r.SecType == SecTypeSigned {
		raw := fmt.Sprintf("%s%s", queryString, bodyString)
		mac := hmac.New(sha256.New, []byte(c.SecretKey))
		_, err = mac.Write([]byte(raw))
		if err != nil {
			return err
		}
		v := url.Values{}
		v.Set(SignatureKey, fmt.Sprintf("%x", mac.Sum(nil)))
		if queryString == "" {
			queryString = v.Encode()
		} else {
			queryString = fmt.Sprintf("%s&%s", queryString, v.Encode())
		}
	}
	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}
	c.Debug("full url: %s, body: %s", fullURL, bodyString)

	r.fullURL = fullURL
	r.header = header
	r.body = body
	return nil
}

func (c *Client) CallAPI(ctx context.Context, r *Request, opts ...RequestOption) (data []byte, err error) {
	err = c.ParseRequest(r, opts...)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest(r.Method, r.fullURL, r.body)
	if err != nil {
		return []byte{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	c.Debug("Request: %#v", req)
	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return []byte{}, err
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	c.Debug("response: %#v", res)
	c.Debug("response body: %s", string(data))
	c.Debug("response status code: %d", res.StatusCode)

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(common.APIError)
		e := json.Unmarshal(data, apiErr)
		if e != nil {
			c.Debug("failed to unmarshal json: %s", e)
		}
		return nil, apiErr
	}

	return data, nil
}

// SetApiEndpoint set api Endpoint
func (c *Client) SetApiEndpoint(url string) *Client {
	c.BaseURL = url
	return c
}
