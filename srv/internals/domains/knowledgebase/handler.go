package knowledgebase

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/", h.create)
}

func (h *Handler) create(ctx *gin.Context) {
	var input CreateRequest
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(input.Embedding) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "embedding must be defined"})
		return
	}

	created, err := h.repository.Create(ctx.Request.Context(), input)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusCreated, created)
}
