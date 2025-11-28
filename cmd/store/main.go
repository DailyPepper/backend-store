package main

import (
	"backend-store/config"
	"backend-store/internal/app"
	"backend-store/internal/auth"
	"backend-store/internal/server"
	"backend-store/pkg/logger"
)

var log logger.Log

func main() {
	cfg := config.Load()
	log = logger.New(cfg.LogLevel)

	log.Info("Starting application", "mode", cfg.Environment)
	log.Info("Server will start on", "address", cfg.ServerHost+":"+cfg.ServerPort)

	authClient, err := auth.NewAuthClient(cfg.AuthServiceURL)
	if err != nil {
		log.Fatal("Failed to connect to auth service:", err)
	}
	defer authClient.Close()
	log.Info("✅ Successfully connected to auth service")

	authHandlers := auth.NewAuthHandlers(authClient)

	application, err := app.New(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}
	defer application.Close()

	router := server.SetupRouter(cfg, application.Handlers, authHandlers, authClient)

	srv := server.New(cfg, log, router)
	if err := srv.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	srv.WaitForShutdown()
}
