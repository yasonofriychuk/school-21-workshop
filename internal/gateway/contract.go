package gateway

import "net/http"

type HTTPGetter interface {
	Get(url string) (resp *http.Response, err error)
}
