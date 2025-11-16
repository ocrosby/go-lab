package http

//go:generate mockgen -destination=./mock_http_client.go -package=http github.com/ocrosby/golab/jsonplaceholder/pkg/http IHttpClient

import "net/http"

type IHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
