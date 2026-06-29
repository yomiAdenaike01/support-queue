package inputsourcesdomain

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository Repository
}

func NewHandler(repository Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("", h.findMany)
	group.POST("", h.create)
	group.PATCH("/:id", h.update)
}

func (h *Handler) findMany(ctx *gin.Context) {
	sources, err := h.repository.FindMany(ctx.Request.Context(), FindInput{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, sources)
}

func (h *Handler) create(ctx *gin.Context) {
	var input SaveInputSourceRequest
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	source, err := h.repository.Create(ctx.Request.Context(), input.toInput(nil))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, source)
}

func (h *Handler) update(ctx *gin.Context) {
	id := ctx.Param("id")
	var input SaveInputSourceRequest
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	source, err := h.repository.Update(ctx.Request.Context(), input.toInput(&id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, source)
}
