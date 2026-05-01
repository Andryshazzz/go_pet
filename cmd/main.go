package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Andryshazzz/go_pet/internal/core/logger"
	core_postgres_pool "github.com/Andryshazzz/go_pet/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/Andryshazzz/go_pet/internal/core/transport/http/middleware"
	core_http_server "github.com/Andryshazzz/go_pet/internal/core/transport/http/server"
	users_postgres_repository "github.com/Andryshazzz/go_pet/internal/features/auth/repository/postgres"
	users_service "github.com/Andryshazzz/go_pet/internal/features/auth/service"
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

	logger.Debug("start init app")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("Starting Application")

	usersrepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersrepository)
	usersTransportHTTP := auth_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("init HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)
	apiVersionsRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionsRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionsRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
