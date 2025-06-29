package main

import (
	"firstmod/internal/config"
	"firstmod/internal/handlers"
	"firstmod/internal/repository"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", ".env", "configuration file")
	flag.Parse()

	cfg := config.MustLoadCfg(configPath)

	log := mustMakeLogger(cfg.LogLevel)

	log.Info("starting server")

	log.Debug("debug messages are enabled")

	// db

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	storage, err := repository.New(log, dsn)
	if err != nil {
		log.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}
	if err := storage.Migrate(); err != nil {
		log.Error("failed to migrate db", "error", err)
		os.Exit(1)
	}

	log.Info("successfully connected to database")

	mux := http.NewServeMux()

	mux.Handle("POST /order", handlers.CreateOrderHandler(log, storage))
	mux.Handle("DELETE /order/{orderID}", handlers.DeleteOrderHandler(log, storage))
	mux.Handle("GET /order/{orderID}", handlers.GetOrderByIDHandler(log, storage))
}

func mustMakeLogger(logLevel string) *slog.Logger {
	return slog.Default()
}
