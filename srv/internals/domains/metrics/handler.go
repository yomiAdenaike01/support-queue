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
	r.GET("/series/tickets-count", h.ticketsOverTime)
}

func (h *Handler) ticketsOverTime(ctx *gin.Context) {
	var input TicketsOvertimeInput
	if err := ctx.ShouldBindQuery(&input); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request expected start and end date",
		})
		return
	}
	tickets, err := h.repository.GetTicketsOvertime(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch tickets over time",
		})
		return
	}
	ctx.JSON(http.StatusOK, tickets)
}
