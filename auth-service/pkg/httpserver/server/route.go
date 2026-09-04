package httpserver

import "net/http"

// Route defines an HTTP endpoint with its method, path, and handler.
type Route struct {
	// Method is the HTTP method (GET, POST, PUT, DELETE, etc.)
	Method string

	// Path is the URL path relative to the API version prefix.
	Path string

	// Handler is the function that processes requests matching this route.
	Handler http.HandlerFunc
}

// NewRoute creates a new Route with the given method, path, and handler.
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
