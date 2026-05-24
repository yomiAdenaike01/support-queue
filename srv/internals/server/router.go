package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	teamdomain "github.com/yomiAdenaike01/support-queue/internals/domains/team"
	ticketdomain "github.com/yomiAdenaike01/support-queue/internals/domains/ticket"
	redisinfra "github.com/yomiAdenaike01/support-queue/internals/infra/redis"
)

type RouterDependencies struct {
	Context      context.Context
	DB           *sqlx.DB
	StreamClient *redisinfra.StreamClient
}

func NewRouter(deps RouterDependencies) *gin.Engine {
	g := gin.Default()
	v1 := g.Group("/api/v1")

	v1.GET("healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	teamHandler := teamdomain.NewHandler(teamdomain.NewRepository())
	teamHandler.RegisterRoutes(v1.Group("/teams"))

	ticketHandler := ticketdomain.NewHandler(deps.Context, deps.DB, deps.StreamClient)
	ticketHandler.RegisterRoutes(v1)

	return g
}
