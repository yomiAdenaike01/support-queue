package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	teamdomain "github.com/yomiAdenaike01/support-queue/internals/domains/team"
	ticketdomain "github.com/yomiAdenaike01/support-queue/internals/domains/ticket"
)

type RouterDependencies struct {
	Context     context.Context
	DB          *sqlx.DB
	RedisClient *redis.Client
	StreamName  string
}

func NewRouter(deps RouterDependencies) *gin.Engine {
	g := gin.Default()
	v1 := g.Group("/api/v1")

	teamHandler := teamdomain.NewHandler(teamdomain.NewRepository())
	teamHandler.RegisterRoutes(g.Group("/teams"))
	teamHandler.RegisterRoutes(v1.Group("/teams"))

	ticketHandler := ticketdomain.NewHandler(deps.Context, deps.DB, deps.RedisClient, deps.StreamName)
	ticketHandler.RegisterRoutes(v1)

	return g
}
