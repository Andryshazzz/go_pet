package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Andryshazzz/go_pet/config"
	"github.com/Andryshazzz/go_pet/internal/app/container"
)

// Run initializes the application and manages its lifecycle.
func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	cnt := container.New(cfg)

	if err := cnt.Start(ctx); err != nil {
		log.Printf("start error: %v", err)
		cnt.Stop()
		os.Exit(1)
	}

	<-ctx.Done()
	log.Println("shutting down...")

	cnt.Stop()
	log.Println("application stopped")
}
