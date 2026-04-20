package main

import (
	"context"
	"log"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Andryshazzz/go_pet/internal/core/config"
	"github.com/Andryshazzz/go_pet/internal/core/http"
	"github.com/Andryshazzz/go_pet/internal/core/postgres"
)

// @title GO PET API
// @version 1.0
// @description Pet project with microservices
// @host localhost:8080
// @BasePath /
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.New(
		ctx,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := chi.NewRouter()

	// transport.RegisterPublicRoutes(r)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	server := http.New(cfg.App.PublicHTTPPort, r)

	log.Println("starting server on port", cfg.App.PublicHTTPPort)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
