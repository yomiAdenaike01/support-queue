package web

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/yomiAdenaike01/support-queue/internals/config"
	inputsourcesdomain "github.com/yomiAdenaike01/support-queue/internals/domains/inputsources"
	integrationsdomain "github.com/yomiAdenaike01/support-queue/internals/domains/integrations"
	"github.com/yomiAdenaike01/support-queue/internals/domains/integrations/oauth"
	knowledgebasedomain "github.com/yomiAdenaike01/support-queue/internals/domains/knowledgebase"
	metricsdomain "github.com/yomiAdenaike01/support-queue/internals/domains/metrics"
	teamdomain "github.com/yomiAdenaike01/support-queue/internals/domains/team"
	ticketdomain "github.com/yomiAdenaike01/support-queue/internals/domains/ticket"
	workerdomain "github.com/yomiAdenaike01/support-queue/internals/domains/worker"
	redisinfra "github.com/yomiAdenaike01/support-queue/internals/infra/redis"
)

type Dependencies struct {
	Context      context.Context
	DB           *sqlx.DB
	StreamClient *redisinfra.StreamClient
	Config       *config.Config
}

func New(deps Dependencies) *gin.Engine {
	g := gin.Default()
	g.ContextWithFallback = true
	v1 := g.Group("/api/v1")

	v1.GET("healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	inputSourcesRepository := inputsourcesdomain.NewRepository(deps.DB)
	inputSourcesHandler := inputsourcesdomain.NewHandler(inputSourcesRepository)
	inputSourcesHandler.RegisterRoutes(v1.Group("/input-sources"))

	ticketRepository := ticketdomain.NewRepository(deps.DB)
	oauthCredsRepo := oauth.NewCredentialsRepository(deps.DB)
	oauthProviders := oauth.NewAuthProviders(deps.Config, oauthCredsRepo)
	integrations := integrationsdomain.New(integrationsdomain.Deps{
		Db:               deps.DB,
		Config:           deps.Config,
		InputSourcesRepo: inputSourcesRepository,
		IngestionChanel:  ticketRepository.GetIngestionChannel(),
		Providers:        oauthProviders,
	})

	teamRepository := teamdomain.NewRepository(deps.DB)
	teamHandler := teamdomain.NewHandler(teamdomain.NewService(teamRepository, integrations), integrations, teamRepository)
	teamHandler.RegisterRoutes(v1.Group("/team"))

	metricsRepository := metricsdomain.NewRepository(deps.DB)
	metricsHandler := metricsdomain.NewHandler(metricsRepository)
	metricsHandler.RegisterRoutes(v1.Group("/metrics"))

	// workerRepository := workerdomain.NewRepository(deps.DB)
	workerHandler := workerdomain.NewHandler(deps.Context, integrations, teamRepository, ticketRepository)
	workerHandler.RegisterRoutes(v1.Group("/worker"))

	integrationsHandler := integrationsdomain.NewHandler(deps.Config, integrations, oauthProviders)
	integrationsHandler.RegisterRoutes(v1.Group("/integrations"))

	knowledgeBaseRepository := knowledgebasedomain.NewRepository(deps.DB)
	knowledgeBaseHandler := knowledgebasedomain.NewHandler(knowledgeBaseRepository)
	knowledgeBaseHandler.RegisterRoutes(v1.Group("/knowledge"))

	ticketHandler := ticketdomain.NewHandler(deps.Context, deps.DB, deps.StreamClient, ticketRepository, teamRepository, integrations)
	ticketHandler.RegisterRoutes(v1)

	return g
}
