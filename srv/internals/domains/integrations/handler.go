package integrationsdomain

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yomiAdenaike01/support-queue/internals/config"
)

type Handler struct {
	integrations *Integrations
	config       *config.Config
}

func (h *Handler) handleInitOauth(ctx *gin.Context) {
	platform := ctx.Query("platform")
	callbackUrl := fmt.Sprintf("%s?platform=%s", h.config.GetEnvOrFail("OAUTH_CALLBACK_URL"), platform)
	url := fmt.Sprintf("https://slack.com?response_type=code&scope=openid,profile,mail&client_id=%s&redirect_uri=%sstate=RANDOM_SECURITY_STRING", h.config.GetEnvOrFail("SLACK_CLIENT_ID"), callbackUrl)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) handleCallback(ctx *gin.Context) {
	ctx.Status(200)
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/oauth", h.handleInitOauth)
	group.GET("/oauth/callback", h.handleCallback)
}

func NewHandler(integrations *Integrations) *Handler {
	return &Handler{
		integrations: integrations,
	}
}
