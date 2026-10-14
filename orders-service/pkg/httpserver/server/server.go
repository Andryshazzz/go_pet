package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Andryshazzz/go_pet/orders-service/docs"
	"github.com/Andryshazzz/go_pet/orders-service/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// Config holds the configuration for the HTTP server.
type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
}

// HTTPServer wraps the Gin engine with graceful shutdown support.
type HTTPServer struct {
	router *gin.Engine
	config Config
	log    *logger.Logger
}

// NewHTTPServer creates a new HTTPServer.
func NewHTTPServer(
	config Config,
	log *logger.Logger,
) *HTTPServer {
	router := gin.New()

	return &HTTPServer{
		router: router,
		config: config,
		log:    log,
	}
}

// Router returns the underlying Gin engine.
func (h *HTTPServer) Router() *gin.Engine {
	return h.router
}

// RegisterSwagger sets up Swagger UI endpoint.
func (h *HTTPServer) RegisterSwagger() {
	if docs.SwaggerInfo.ReadDoc() == "" {
		h.log.Warn("swagger documentation is empty, skipping registration")
		return
	}

	h.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Run starts the HTTP server and blocks until the context is cancelled.
func (h *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: h.router,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		h.log.Debug("start HTTP server", zap.String("addr", h.config.Addr))

		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}

	case <-ctx.Done():
		h.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			h.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			if closeErr := server.Close(); closeErr != nil {
				h.log.Error("force close HTTP server", zap.Error(closeErr))
			}
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		h.log.Warn("HTTP server stopped")
	}

	return nil
}
