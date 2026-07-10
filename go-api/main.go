package main

import (
	"log"

	"go-api/config"
	"go-api/server"
)

func main() {
	cfg := config.Load()
	srv := server.New(cfg)

	if err := srv.Start(); err != nil {
		log.Fatalf("Error durante la ejecución del servidor: %v", err)
	}
}
