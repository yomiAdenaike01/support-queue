package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/yomiAdenaike01/support-queue/internals/config"
	"github.com/yomiAdenaike01/support-queue/internals/infra/database"
	redisinfra "github.com/yomiAdenaike01/support-queue/internals/infra/redis"
	"github.com/yomiAdenaike01/support-queue/internals/server"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}
	cfg := config.New()
	db := database.NewClient(cfg)

	shutdownCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	streamClient, err := redisinfra.NewStreamClient(shutdownCtx, cfg)
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(shutdownCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	router := server.NewRouter(server.RouterDependencies{
		Context:      shutdownCtx,
		DB:           db,
		StreamClient: streamClient,
	})
	server.Run(ctx, cfg, router)
}
