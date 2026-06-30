package teamdomain

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	integrationsdomain "github.com/yomiAdenaike01/support-queue/internals/domains/integrations"
)

type Handler struct {
	service      *Service
	repository   *Repository
	integrations *integrationsdomain.Integrations
}

func NewHandler(service *Service, integration *integrationsdomain.Integrations, repository *Repository) *Handler {
	return &Handler{service: service, integrations: integration, repository: repository}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("", h.list)
	group.POST("/", h.create)
	group.GET("/search", h.search)
}

func (h *Handler) list(ctx *gin.Context) {
	teams, err := h.service.FindTeams(ctx.Request.Context(), nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, teams)
}

func (h *Handler) search(ctx *gin.Context) {
	id := ctx.Query("id")
	departmentsAsString := strings.ReplaceAll(ctx.Query("departments"), " ", "")
	if id == "" && departmentsAsString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No filters provided",
		})
		return
	}
	deps, err := url.PathUnescape(departmentsAsString)
	if err != nil {
		panic(err)
	}
	var departments []string
	if deps != "" {
		departments = strings.Split(deps, ",")
	}

	team, err := h.service.FindTeams(ctx.Request.Context(), &Filters{
		Id:          id,
		Departments: departments,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusFound, team)
}

func (h *Handler) create(ctx *gin.Context) {
	var input CreateTeamRequest
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	if input.Department == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Department must be defined"})
		return
	}
	if len(input.Members) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "At least one team member must be defined"})
		return
	}
	for _, member := range input.Members {
		if member.Name == "" || member.Role == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Team member name and role must be defined"})
			return
		}
	}

	createdTeam, created, err := h.repository.Create(ctx.Request.Context(), input)
	if err != nil {
		panic(err)
	}
	if !created {
		ctx.JSON(http.StatusConflict, gin.H{"errors": "Team already exists"})
		return
	}

	ctx.JSON(http.StatusCreated, createdTeam)
}
