package main

import (
	"log"

	"github.com/Andryshazzz/go_pet/config"
	_ "github.com/Andryshazzz/go_pet/docs"
	"github.com/Andryshazzz/go_pet/internal/app"
)

// @title        Golang app API
// @version      1.0
// @description  API for Golang app
// @host 		 127.0.0.1:5050
// @BasePath 	 /api/v1
const configDir = "./config/main.yaml"

func main() {
	cfg, err := config.NewConfig(configDir)

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}

// func main() {
// 	ctx, cancel := signal.NotifyContext(
// 		context.Background(),
// 		syscall.SIGINT, syscall.SIGTERM,
// 	)

// 	defer cancel()

// 	logger, err := logger.NewLogger(logger.NewConfigMust())

// 	if err != nil {
// 		fmt.Println("Failed to init application logger:", err)

// 		os.Exit(1)
// 	}

// 	defer logger.Close()

// 	logger.Debug("Start init app")

// 	pool, err := postgrespool.NewConnectionPool(
// 		ctx,
// 		postgrespool.NewConfigMust(),
// 	)

// 	if err != nil {
// 		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
// 	}

// 	defer pool.Close()

// 	jwtConfig := jwt.NewJWTConfigMust()

// 	jwtService := jwt.NewJWTService(
// 		jwtConfig.Secret,
// 		jwtConfig.AccessExpiration,
// 		jwtConfig.RefreshExpiration,
// 	)

// 	logger.Debug("Starting Application")

// 	usersrepository := usersrepository.NewUsersRepository(pool)
// 	usersService := usersservice.NewUsersService(usersrepository, jwtService)
// 	usersTransportHTTP := authtransport.NewUsersHTTPHandler(usersService)

// 	logger.Debug("init HTTP server")

// 	httpServer := httpserver.NewHTTPServer(
// 		httpserver.NewConfigMust(),
// 		logger,
// 		httpmiddleware.CORS(),
// 		httpmiddleware.RequestID(),
// 		httpmiddleware.Logger(logger),
// 		httpmiddleware.Panic(),
// 		httpmiddleware.Trace(),
// 	)

// 	apiVersionsRouter := httpserver.NewAPIVersionRouter(httpserver.ApiVersion1)
// 	apiVersionsRouter.RegisterRoutes(usersTransportHTTP.Routes()...)

// 	httpServer.RegisterAPIRouters(apiVersionsRouter)

// 	httpServer.RegisterSwagger()

// 	if err := httpServer.Run(ctx); err != nil {
// 		logger.Error("HTTP server run error", zap.Error(err))
// 	}
// }
