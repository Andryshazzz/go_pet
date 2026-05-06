package core_http_server

import "net/http"

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

func NewRoute(
	method string,
	path string,
	handler http.HandlerFunc,
) *Route {
	return &Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
