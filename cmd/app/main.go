package main

import (
	"context"
	"errors"
	"firstmod/internal/cache"
	"firstmod/internal/config"
	"firstmod/internal/handlers"
	"firstmod/internal/kafka"
	"firstmod/internal/repository"
	"firstmod/internal/service"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
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

	orderCache := cache.NewCache(log)
	log.Info("in-memory cache initialized")

	kafkaBrokers := strings.Split(cfg.KafkaBrokers, ",")
	kafkaProducer := kafka.NewProducer(log, kafkaBrokers, cfg.KafkaTopic)
	defer kafkaProducer.Close()
	log.Info("Kafka producer initialized")

	orderService := service.NewOrderService(storage, orderCache, log, kafkaProducer)
	log.Info("order service initialized")

	loadCtx, cancelLoad := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelLoad()

	if err := orderService.LoadCacheFromDB(loadCtx); err != nil {
		log.Error("failed to load cache from database", "error", err)
	} else {
		log.Info("cache successfully loaded from database")
	}

	mux := http.NewServeMux()

	mux.Handle("POST /order", handlers.CreateOrderHandler(log, storage))
	mux.Handle("DELETE /order/{orderID}", handlers.DeleteOrderHandler(log, storage))
	mux.Handle("GET /order/{orderID}", handlers.GetOrderByIDHandler(log, storage))
	mux.Handle("GET /orders/", handlers.GetOrdersIDsHandler(log, storage))

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fileServer)

	server := http.Server{
		Addr:        cfg.HttpServerAddress,
		ReadTimeout: cfg.HttpServerTimeout * time.Second,
		Handler:     mux,
	}

	log.Info("server is listening on", "address", cfg.HttpServerAddress)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Debug("shutting down server")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error("erroneous shutdown", "error", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("server closed unexpectedly", "error", err)
			return
		}
	}
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	default:
		panic("unknown log level: " + logLevel)
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level, AddSource: true})
	return slog.New(handler)
}
