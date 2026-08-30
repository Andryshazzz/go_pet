package container

import (
	"context"
	"fmt"

	"github.com/Andryshazzz/go_pet/config"
	usersservice "github.com/Andryshazzz/go_pet/internal/domain/services/user"
	"github.com/Andryshazzz/go_pet/internal/endpoint/controller/http/auth"
	usersrepository "github.com/Andryshazzz/go_pet/internal/infrastructure/repository/user"
	postgrespool "github.com/Andryshazzz/go_pet/pkg/database/postgres/pool"
	httpmiddleware "github.com/Andryshazzz/go_pet/pkg/httpserver/middleware"
	httpserver "github.com/Andryshazzz/go_pet/pkg/httpserver/server"
	"github.com/Andryshazzz/go_pet/pkg/jwt"
	"github.com/Andryshazzz/go_pet/pkg/logger"
)

// Container holds all application dependencies.
type Container struct {
	cfg    *config.Config
	logger *logger.Logger
	pool   *postgrespool.ConnectionPool
	server *httpserver.HTTPServer

	userRepository *usersrepository.UsersRepository
	userService    *usersservice.UsersService
	userHandler    *auth.UsersHTTPHandler
	jwtService     *jwt.JWTService
}

// New creates a new Container with configuration.
func New(cfg *config.Config) *Container {
	return &Container{cfg: cfg}
}

// Start initializes all dependencies and starts the server.
func (c *Container) Start(ctx context.Context) error {
	if err := c.initLogger(); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	if err := c.initDatabase(ctx); err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	c.initJWT()
	c.initLayers()
	c.initHTTPServer()

	return c.server.Run(ctx)
}

// Stop gracefully shuts down all dependencies.
func (c *Container) Stop() {
	if c.pool != nil {
		c.pool.Close()
	}
	if c.logger != nil {
		c.logger.Close()
	}
}

func (c *Container) initLogger() error {
	log, err := logger.NewLogger(logger.Config{
		Level:  c.cfg.Logger.Level,
		Folder: c.cfg.Logger.Folder,
	})
	if err != nil {
		return err
	}
	c.logger = log
	return nil
}

func (c *Container) initDatabase(ctx context.Context) error {
	pool, err := postgrespool.NewConnectionPool(ctx, postgrespool.Config{
		Host:     c.cfg.Postgres.Host,
		Port:     c.cfg.Postgres.Port,
		User:     c.cfg.Postgres.User,
		Password: c.cfg.Postgres.Password,
		DB:       c.cfg.Postgres.DB,
		Timeout:  c.cfg.Postgres.Timeout,
	})
	if err != nil {
		return err
	}
	c.pool = pool
	return nil
}

func (c *Container) initJWT() {
	c.jwtService = jwt.NewJWTService(jwt.Config{
		Secret:            c.cfg.JWT.Secret,
		AccessExpiration:  c.cfg.JWT.AccessExpiration,
		RefreshExpiration: c.cfg.JWT.RefreshExpiration,
	})
}

func (c *Container) initLayers() {
	c.userRepository = usersrepository.NewUsersRepository(c.pool)
	c.userService = usersservice.NewUsersService(c.userRepository, c.jwtService)
	c.userHandler = auth.NewUsersHTTPHandler(c.userService)
}

func (c *Container) initHTTPServer() {
	c.server = httpserver.NewHTTPServer(
		httpserver.Config{
			Addr:            c.cfg.HTTP.Addr,
			ShutdownTimeout: c.cfg.HTTP.ShutdownTimeout,
		},
		c.logger,
		httpmiddleware.CORS(),
		httpmiddleware.RequestID(),
		httpmiddleware.Logger(c.logger),
		httpmiddleware.Panic(),
		httpmiddleware.Trace(),
	)

	apiV1 := httpserver.NewAPIVersionRouter(httpserver.ApiVersion1)
	apiV1.RegisterRoutes(c.userHandler.Routes()...)
	c.server.RegisterAPIRouters(apiV1)
	c.server.RegisterSwagger()
}
