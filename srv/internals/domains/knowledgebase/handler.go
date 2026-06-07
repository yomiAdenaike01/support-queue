package knowledgebase

import (
	"encoding/base64"
	"encoding/json"
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
	group.GET("/search/", h.search)
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

type searchFilter struct {
	Id         *string `uri:"id"`
	Embedding  *string `uri:"embedding"`
	Limit      *int    `uri:"limit"`
	SourceType *string `uri:"source_type"`
}

func (s searchFilter) toSearchInput() (searchInput, error) {
	var search searchInput
	embedding, err := s.getEmbedding()
	if err != nil {
		return search, err
	}
	return searchInput{
			Embedding:  embedding,
			Limit:      *s.Limit,
			SourceType: s.SourceType,
			Id:         s.Id,
		},
		nil
}

func (f searchFilter) getEmbedding() ([]float32, error) {
	if len(*f.Embedding) == 0 {
		return nil, nil
	}
	rawStr, err := base64.StdEncoding.DecodeString(*f.Embedding)
	if err != nil {
		return nil, err
	}
	var embedding []float32
	if err := json.Unmarshal(rawStr, &embedding); err != nil {
		return nil, err
	}
	return embedding, nil

}

func (h *Handler) search(ctx *gin.Context) {
	var filter searchFilter
	if err := ctx.ShouldBindUri(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	searchInput, err := filter.toSearchInput()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse embedding",
		})
		return
	}

	sources, err := h.repository.Find(ctx.Request.Context(), searchInput)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusFound, sources)
}
