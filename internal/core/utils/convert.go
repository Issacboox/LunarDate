package utils

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/valyala/fasthttp"
)

func ConvertRequest(fasthttpRequest *fasthttp.Request) (*http.Request, error) {
	req := &http.Request{
		Method: string(fasthttpRequest.Header.Method()),
		URL: &url.URL{
			Scheme: string(fasthttpRequest.URI().Scheme()),
			Host:   string(fasthttpRequest.URI().Host()),
			Path:   string(fasthttpRequest.URI().Path()),
		},
		Header: make(http.Header),
		Body:   io.NopCloser(bytes.NewReader(fasthttpRequest.Body())),
	}

	fasthttpRequest.Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})

	return req, nil
}
