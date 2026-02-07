package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"online-subscription-service/internal/connections"
	"online-subscription-service/internal/handler"
	"online-subscription-service/internal/repository"
	"online-subscription-service/internal/usecase"
	l "online-subscription-service/pkg/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := gin.Default()
	r.Use(cors.Default())

	textHandler := slog.NewTextHandler(os.Stdout, nil)
	logger := l.NewSlogWrapper(slog.New(textHandler))

	if err := godotenv.Load(); err != nil {
		logger.Error("failed loading .env file")
	}

	connection := connections.NewConnectionPool(logger)
	pool := connection.InitPool(ctx)

	subscriptionRepository := repository.NewSubscriptionRepository(pool)
	subscriptionService := usecase.NewSubscriptionService(subscriptionRepository, logger)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	subscriptionHandler.RegisterRoute(r)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("server error: ", "error", err.Error())
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	logger.Info("shutting down server...")
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error: ", "error", err.Error())
	}

	defer pool.Close()
}
