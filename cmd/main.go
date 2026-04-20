package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Andryshazzz/go_pet/internal/core/logger"
	core_http_server "github.com/Andryshazzz/go_pet/internal/core/transport/http/server"
	auth_transport_http "github.com/Andryshazzz/go_pet/internal/features/auth/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("Starting Application")

	usersTransportHTTP := auth_transport_http.NewUsersHTTPHandler(nil)
	usersRoutes := usersTransportHTTP.Routes()

	apiVersionsRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionsRouter.RegisterRoutes(usersRoutes...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
	)
	httpServer.RegisterAPIRouters(apiVersionsRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
