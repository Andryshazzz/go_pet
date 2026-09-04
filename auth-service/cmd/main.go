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
