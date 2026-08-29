package httpmiddleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Andryshazzz/go_pet/pkg/logger"
	httpresponse "github.com/Andryshazzz/go_pet/pkg/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

// CORS returns a middleware that handles Cross-Origin Resource Sharing.
// It allows requests from specified origins and handles preflight OPTIONS requests.
//
// Allowed origins:
//   - http://localhost:5050
//
// Allowed methods: GET, POST, PATCH, DELETE, OPTIONS
// Allowed headers: Content-Type, Authorization
func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := map[string]struct{}{
				"http://localhost:5050": {},
			}

			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequestID returns a middleware that ensures every request has a unique ID.
// If the client sends an X-Request-ID header, it is used as-is.
// Otherwise, a new UUID v4 is generated.
//
// The request ID is set on both the request and response headers,
// allowing end-to-end request tracing.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// Logger returns a middleware that injects a request-scoped logger
// into the request context. The logger includes the request ID and URL.
func Logger(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "log", l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Panic returns a middleware that recovers from panics in downstream handlers.
// When a panic occurs, it logs the error and returns a 500 Internal Server Error
// response instead of crashing the server.
//
// This middleware MUST be placed before any handlers that might panic.
func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			responseHandler := httpresponse.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// Trace returns a middleware that logs incoming requests and their outcomes.
// It records the HTTP method, URL, status code, and request latency.
//
// Log format:
//   >>> Incoming HTTP request  {"time": "..."}
//   <<< done HTTP request      {"status_code": 200, "latency": "1.5ms"}
func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			rw := httpresponse.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> Incoming HTTP request",
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", rw.GetStatusCodeOrPanic()),
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}
