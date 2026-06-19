package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/yomiAdenaike01/support-queue/internals/config"
	integrationsdomain "github.com/yomiAdenaike01/support-queue/internals/domains/integrations"
	"github.com/yomiAdenaike01/support-queue/internals/infra/database"
	redisinfra "github.com/yomiAdenaike01/support-queue/internals/infra/redis"
	"github.com/yomiAdenaike01/support-queue/internals/server"
)

func main() {
	godotenv.Load("../.env") // non-fatal; env vars injected externally (Docker) take precedence
	cfg := config.New()
	db := database.NewClient(cfg)

	shutdownCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	integrations := integrationsdomain.NewIntegrations(cfg)

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
		Integrations: integrations,
	})
	server.Run(ctx, cfg, router)
}
