package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

type SecType int

const (
	SecTypeNone SecType = iota
	SecTypeAPIKey
	SecTypeSigned // if the 'timestamp' parameter is required
)

type params map[string]interface{}

// Request define an API Request
type Request struct {
	Method     string
	Endpoint   string
	query      url.Values
	form       url.Values
	recvWindow int64
	SecType    SecType
	header     http.Header
	body       io.Reader
	fullURL    string
}

// AddParam add param with key/value to query string
func (r *Request) AddParam(key string, value interface{}) *Request {
	if r.query == nil {
		r.query = url.Values{}
	}
	r.query.Add(key, fmt.Sprintf("%v", value))
	return r
}

// SetParam set param with key/value to query string
func (r *Request) SetParam(key string, value interface{}) *Request {
	if r.query == nil {
		r.query = url.Values{}
	}

	if reflect.TypeOf(value).Kind() == reflect.Slice {
		v, err := json.Marshal(value)
		if err == nil {
			value = string(v)
		}
	}

	r.query.Set(key, fmt.Sprintf("%v", value))
	return r
}

// SetParams set params with key/values to query string
func (r *Request) SetParams(m params) *Request {
	for k, v := range m {
		r.SetParam(k, v)
	}
	return r
}

// SetFormParam set param with key/value to Request form body
func (r *Request) SetFormParam(key string, value interface{}) *Request {
	if r.form == nil {
		r.form = url.Values{}
	}
	r.form.Set(key, fmt.Sprintf("%v", value))
	return r
}

// SetFormParams set params with key/values to Request form body
func (r *Request) SetFormParams(m params) *Request {
	for k, v := range m {
		r.SetFormParam(k, v)
	}
	return r
}

func (r *Request) validate() (err error) {
	if r.query == nil {
		r.query = url.Values{}
	}
	if r.form == nil {
		r.form = url.Values{}
	}
	return nil
}

// RequestOption define option type for Request
type RequestOption func(*Request)

// WithRecvWindow set recvWindow param for the Request
func WithRecvWindow(recvWindow int64) RequestOption {
	return func(r *Request) {
		r.recvWindow = recvWindow
	}
}

// WithHeader set or add a header value to the Request
func WithHeader(key, value string, replace bool) RequestOption {
	return func(r *Request) {
		if r.header == nil {
			r.header = http.Header{}
		}
		if replace {
			r.header.Set(key, value)
		} else {
			r.header.Add(key, value)
		}
	}
}

// WithHeaders set or replace the headers of the Request
func WithHeaders(header http.Header) RequestOption {
	return func(r *Request) {
		r.header = header.Clone()
	}
}
