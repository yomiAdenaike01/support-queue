package integrationsdomain

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yomiAdenaike01/support-queue/internals/config"
	"github.com/yomiAdenaike01/support-queue/internals/domains/integrations/oauth"
)

type Handler struct {
	integrations   *Integrations
	config         *config.Config
	oauthProviders *oauth.Providers
}

func (h *Handler) handleInitOauth(ctx *gin.Context) {
	provider := oauth.ProviderName(ctx.Query("provider"))
	ctx.Redirect(http.StatusTemporaryRedirect, h.oauthProviders.GetAuthorisationEndpoint(provider))
}

func (h *Handler) handleCallback(ctx *gin.Context) {
	var callbackQueryParams oauth.Callback
	if err := ctx.ShouldBindQuery(&callbackQueryParams); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	go func() {
		err := h.oauthProviders.ExchangeAndSaveTokens(callbackQueryParams.Provider, callbackQueryParams.Code)
		if err != nil {
			log.Println("Failed to save tokens")
		}
	}()

	ctx.Status(http.StatusOK)

}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/oauth", h.handleInitOauth)
	group.GET("/oauth/callback", h.handleCallback)
}

func NewHandler(config *config.Config, integrations *Integrations, oauthProviders *oauth.Providers) *Handler {
	return &Handler{
		config:         config,
		oauthProviders: oauthProviders,
		integrations:   integrations,
	}
}
