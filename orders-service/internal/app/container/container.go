package container

import (
	"context"
	"fmt"

	"github.com/Andryshazzz/go_pet/orders-service/config"
	ordersservice "github.com/Andryshazzz/go_pet/orders-service/internal/domain/services/orders"
	"github.com/Andryshazzz/go_pet/orders-service/internal/endpoint/controller/http/orders"
	kafkainfra "github.com/Andryshazzz/go_pet/orders-service/internal/infrastructure/kafka"
	ordersrepository "github.com/Andryshazzz/go_pet/orders-service/internal/infrastructure/repository/user"
	"github.com/Andryshazzz/go_pet/orders-service/pkg/jwt"
	"github.com/Andryshazzz/go_pet/orders-service/pkg/kafka"
	"github.com/gin-gonic/gin"

	postgrespool "github.com/Andryshazzz/go_pet/orders-service/pkg/database/postgres/pool"
	httpmiddleware "github.com/Andryshazzz/go_pet/orders-service/pkg/httpserver/middleware"
	httpserver "github.com/Andryshazzz/go_pet/orders-service/pkg/httpserver/server"
	"github.com/Andryshazzz/go_pet/orders-service/pkg/logger"
)

// Container holds all application dependencies.
type Container struct {
	cfg         *config.Config
	logger      *logger.Logger
	pool        *postgrespool.ConnectionPool
	server      *httpserver.HTTPServer
	kafkaClient *kafka.Producer
	jwtService  *jwt.JWTService

	orderRepository *ordersrepository.OrdersRepository
	orderService    *ordersservice.OrdersService
	orderHandler    *orders.OrdersHTTPHandler
	orderEventsPub  *kafkainfra.Producer
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
		SSLMode:  c.cfg.Postgres.SSLMode,
		Timeout:  c.cfg.Postgres.Timeout,
	})
	if err != nil {
		return err
	}
	c.pool = pool
	return nil
}

func (c *Container) initHTTPServer() {
	c.server = c.getHTTPServer()
}

func (c *Container) getJWTService() *jwt.JWTService {
	if c.jwtService == nil {
		c.jwtService = jwt.NewJWTService(jwt.Config{
			Secret: c.cfg.JWT.Secret,
		})
	}
	return c.jwtService
}

func (c *Container) getLogger() *logger.Logger {
	if c.logger == nil {
		log, err := logger.NewLogger(logger.Config{
			Level:  c.cfg.Logger.Level,
			Folder: c.cfg.Logger.Folder,
		})
		if err != nil {
			panic(fmt.Errorf("failed to create logger: %w", err))
		}
		c.logger = log
	}
	return c.logger
}

func (c *Container) getPool() *postgrespool.ConnectionPool {
	if c.pool == nil {
		pool, err := postgrespool.NewConnectionPool(context.Background(), postgrespool.Config{
			Host:     c.cfg.Postgres.Host,
			Port:     c.cfg.Postgres.Port,
			User:     c.cfg.Postgres.User,
			Password: c.cfg.Postgres.Password,
			DB:       c.cfg.Postgres.DB,
			SSLMode:  c.cfg.Postgres.SSLMode,
			Timeout:  c.cfg.Postgres.Timeout,
		})
		if err != nil {
			panic(fmt.Errorf("failed to create database pool: %w", err))
		}
		c.pool = pool
	}
	return c.pool
}

func (c *Container) getKafkaClient() *kafka.Producer {
	if c.kafkaClient == nil {
		client, err := kafka.NewProducer(kafka.Config{
			Brokers:  c.cfg.Kafka.Brokers,
			ClientID: c.cfg.Kafka.ClientID,
		})
		if err != nil {
			panic(fmt.Errorf("failed to create kafka client: %w", err))
		}
		c.kafkaClient = client
	}
	return c.kafkaClient
}

func (c *Container) getOrderEventsPublisher() *kafkainfra.Producer {
	if c.orderEventsPub == nil {
		c.orderEventsPub = kafkainfra.NewProducer(
			c.getKafkaClient(),
			kafkainfra.Config{
				OrdersCreatedTopic: c.cfg.Kafka.OrdersCreatedTopic,
				OrderStatusTopic:   c.cfg.Kafka.OrderStatusTopic,
			},
		)
	}
	return c.orderEventsPub
}

func (c *Container) getOrderRepository() *ordersrepository.OrdersRepository {
	if c.orderRepository == nil {
		c.orderRepository = ordersrepository.NewOrdersRepository(c.getPool())
	}
	return c.orderRepository
}

func (c *Container) getOrderService() *ordersservice.OrdersService {
	if c.orderService == nil {
		c.orderService = ordersservice.NewOrdersService(c.getOrderRepository(), c.getOrderEventsPublisher())
	}
	return c.orderService
}

func (c *Container) getOrderHandler() *orders.OrdersHTTPHandler {
	if c.orderHandler == nil {
		c.orderHandler = orders.NewOrdersHTTPHandler(c.getOrderService())
	}
	return c.orderHandler
}

func (c *Container) getHTTPServer() *httpserver.HTTPServer {
	if c.server == nil {
		c.server = httpserver.NewHTTPServer(
			httpserver.Config{
				Addr:            c.cfg.HTTP.Addr,
				ShutdownTimeout: c.cfg.HTTP.ShutdownTimeout,
			},
			c.getLogger(),
		)

		router := c.server.Router()
		router.Use(gin.Logger())
		router.Use(gin.Recovery())

		apiV1 := router.Group("/api/v1")
		apiV1.Use(httpmiddleware.JWTAuth(c.getJWTService()))
		
		c.getOrderHandler().RegisterRoutes(apiV1)
		c.server.RegisterSwagger()
	}
	return c.server
}
