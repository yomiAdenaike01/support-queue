package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	integrationsdomain "github.com/yomiAdenaike01/support-queue/internals/domains/integrations"
	teamdomain "github.com/yomiAdenaike01/support-queue/internals/domains/team"
	ticketdomain "github.com/yomiAdenaike01/support-queue/internals/domains/ticket"
	workerdomain "github.com/yomiAdenaike01/support-queue/internals/domains/worker"
	redisinfra "github.com/yomiAdenaike01/support-queue/internals/infra/redis"
)

type RouterDependencies struct {
	Context      context.Context
	DB           *sqlx.DB
	StreamClient *redisinfra.StreamClient
	Integrations *integrationsdomain.Integrations
}

func NewRouter(deps RouterDependencies) *gin.Engine {
	g := gin.Default()
	v1 := g.Group("/api/v1")

	v1.GET("healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	teamRepository := teamdomain.NewRepository(deps.DB)
	teamHandler := teamdomain.NewHandler(teamRepository)
	workerHandler := workerdomain.NewHandler(deps.Context, deps.Integrations, teamRepository)
	integrationsHandler := integrationsdomain.NewHandler(deps.Integrations)

	integrationsHandler.RegisterRoutes(v1.Group("/integrations"))

	workerHandler.RegisterRoutes(v1.Group("/worker"))
	teamHandler.RegisterRoutes(v1.Group("/team"))

	ticketHandler := ticketdomain.NewHandler(deps.Context, deps.DB, deps.StreamClient)
	ticketHandler.RegisterRoutes(v1)

	return g
}
