package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Andryshazzz/go_pet/docs"
	httpmiddleware "github.com/Andryshazzz/go_pet/pkg/httpserver/middleware"
	"github.com/Andryshazzz/go_pet/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// Config holds the configuration for the HTTP server.
// Values are passed from the application configuration.
type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
}

// HTTPServer wraps the standard http.Server with middleware support,
// API versioning, Swagger documentation, and graceful shutdown.
type HTTPServer struct {
	mux        *http.ServeMux
	config     Config
	log        *logger.Logger
	middleware []httpmiddleware.Middleware
}

// NewHTTPServer creates a new HTTPServer with the given configuration,
// logger, and optional middleware chain.
//
// Middleware is applied in order to all registered routes.
func NewHTTPServer(
	config Config,
	log *logger.Logger,
	middleware ...httpmiddleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

// RegisterAPIRouters registers versioned API routers under the /api/{version} prefix.
// Each router handles its own set of routes.
func (h *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		h.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

// RegisterSwagger sets up Swagger UI and JSON spec endpoints.
//
// Endpoints:
//   - GET /swagger/       — Swagger UI interface
//   - GET /swagger/doc.json — OpenAPI JSON specification
func (h *HTTPServer) RegisterSwagger() {
	if docs.SwaggerInfo.ReadDoc() == "" {
		h.log.Warn("Swagger documentation is empty, skipping registration")

		return
	}

	h.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		),
	)

	h.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}

// Run starts the HTTP server and blocks until the context is cancelled
// or a fatal error occurs. On context cancellation, it performs a graceful
// shutdown with the configured timeout.
func (h *HTTPServer) Run(ctx context.Context) error {
	mux := httpmiddleware.ChainMiddleware(h.mux, h.middleware...)

	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		h.log.Debug("Start HTTP server", zap.String("Addr", h.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("Listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("Shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			h.config.ShutdownTimeout,
		)

		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			if err := server.Close(); err != nil {
				h.log.Error("Force close HTTP server", zap.Error(err))
			}

			return fmt.Errorf("Shutdown HTTP server: %w", err)
		}

		h.log.Warn("HTTP server stopped")
	}

	return nil
}
