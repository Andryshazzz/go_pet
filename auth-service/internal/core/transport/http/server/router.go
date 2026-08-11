package httpserver

import (
	"fmt"
	"net/http"
)

// ApiVersion.
type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
)

// APIVersionRouter is a versioned HTTP router that groups routes
// under a specific API version prefix. It embeds http.ServeMux
// and adds method-based routing support.
type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
}

// NewAPIVersionRouter creates a new router for the given API version.
func NewAPIVersionRouter(apiVersion ApiVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

// RegisterRoutes registers one or more routes with this router.
func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.Handler)
	}
}
