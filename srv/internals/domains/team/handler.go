package team

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
	group.GET("", h.list)
	group.POST("", h.create)
	group.GET("/:slug", h.get)
	group.PUT("/:slug", h.update)
	group.DELETE("/:slug", h.delete)
}

func (h *Handler) list(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, h.repository.List())
}

func (h *Handler) create(ctx *gin.Context) {
	var team Team
	if err := ctx.ShouldBindBodyWithJSON(&team); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	if team.Slug == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Slug must be defined"})
		return
	}
	if team.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Name must be defined"})
		return
	}

	createdTeam, ok := h.repository.Create(team)
	if !ok {
		ctx.JSON(http.StatusConflict, gin.H{"errors": "Team already exists"})
		return
	}
	ctx.JSON(http.StatusCreated, createdTeam)
}

func (h *Handler) get(ctx *gin.Context) {
	team, ok := h.repository.Get(ctx.Param("slug"))
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"errors": "Team not found"})
		return
	}
	ctx.JSON(http.StatusOK, team)
}

func (h *Handler) update(ctx *gin.Context) {
	slug := ctx.Param("slug")
	var team Team
	if err := ctx.ShouldBindBodyWithJSON(&team); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	if team.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Name must be defined"})
		return
	}
	if team.Slug != "" && team.Slug != slug {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Slug cannot differ from path"})
		return
	}

	updatedTeam, ok := h.repository.Update(slug, team)
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"errors": "Team not found"})
		return
	}
	ctx.JSON(http.StatusOK, updatedTeam)
}

func (h *Handler) delete(ctx *gin.Context) {
	if !h.repository.Delete(ctx.Param("slug")) {
		ctx.JSON(http.StatusNotFound, gin.H{"errors": "Team not found"})
		return
	}
	ctx.Status(http.StatusNoContent)
}
