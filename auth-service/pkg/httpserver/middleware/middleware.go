package httpmiddleware

import "net/http"

// Middleware is a function that wraps an http.Handler with additional logic.
// Middleware can be chained together using ChainMiddleware.
type Middleware func(http.Handler) http.Handler

// ChainMiddleware applies multiple middleware to a handler in order.
// The first middleware in the list wraps the handler first and runs first.
//
// Execution order for ChainMiddleware(h, A, B, C):
//   Request → A → B → C → handler
//   Response ← A ← B ← C ← handler
func ChainMiddleware(
	h http.Handler,
	m ...Middleware,
) http.Handler {
	if len(m) == 0 {
		return h
	}

	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}
