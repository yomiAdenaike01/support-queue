package metricsdomain

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository *Repository
}

func (h *Handler) GetSummary(ctx *gin.Context) {
	summary, err := h.repository.GetDashboardSummary()
	if err != nil {
		log.Printf("[GetSummary]: Failed to fetch dashboard summary reason=%s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, summary)
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/summary", h.GetSummary)
}
