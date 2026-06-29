package ticket

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	redisinfra "github.com/yomiAdenaike01/support-queue/internals/infra/redis"
)

type Handler struct {
	repository   *Repository
	streamClient *redisinfra.StreamClient
	ctx          context.Context
}

func NewHandler(ctx context.Context, db *sqlx.DB, streamClient *redisinfra.StreamClient, repository *Repository) *Handler {
	return &Handler{
		repository:   repository,
		streamClient: streamClient,
		ctx:          ctx,
	}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/tickets", h.create)
	group.GET("/tickets/:id", h.findById)
	group.GET("/tickets", h.find)
	group.POST("/tickets/:id/message", h.insertMessage)
	group.POST("/tickets/:id/events", h.addEvent)
	group.POST("/tickets/:id/resolve", h.resolveTicket)
}

type IdParam struct {
	Id string `uri:"id" binding:"required"`
}

func (h *Handler) resolveTicket(ctx *gin.Context) {
	var params IdParam

	if err := ctx.ShouldBindUri(&params); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to resolve ticket",
		})
		return
	}
	status := "RESOLVED"

	if err := h.repository.FindAndUpdate(ctx.Request.Context(), params.Id, UpdateTicketInput{
		Status: &status,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to resolve ticket",
		})
		return
	}
	if err := h.streamClient.Push(ctx.Request.Context(), redisinfra.PushEvent{
		StreamName: redisinfra.RESOLVED_TICKET_STREAM,
		Values: map[string]any{
			"ticket_id": params.Id,
		},
		EventType: redisinfra.TICKET_RESOLVED,
	}); err != nil {
		log.Printf("Failed to push ticket_id=%s to stream reason=%s", params.Id, err.Error())
	}
	ctx.JSON(http.StatusOK, nil)
}

func (h *Handler) addEvent(ctx *gin.Context) {
	var params IdParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	var addEventBody CreateEventInput
	if err := ctx.ShouldBindBodyWithJSON(&addEventBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	addEventBody.TicketId = params.Id

	event, err := h.repository.CreateEvent(addEventBody)
	if err != nil {
		if !errors.Is(err, ErrDuplicateTicketEvent) {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, event)
}

type findQueryParams struct {
	Search   string `uri:"search"`
	Status   string `uri:"status"`
	Priority string `uri:"priority"`
	Category string `uri:"category"`
	Limit    int    `uri:"limit"`
}

func (f findQueryParams) toRepoFindFilters() FindFilters {
	if f.Limit == 0 {
		f.Limit = 20
	}
	return FindFilters{
		Search:   f.Search,
		Status:   f.Status,
		Priority: f.Priority,
		Category: f.Category,
		Limit:    f.Limit,
	}
}
func (h *Handler) find(ctx *gin.Context) {
	var queryParams findQueryParams
	if err := ctx.ShouldBindUri(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	repoFilters := queryParams.toRepoFindFilters()
	response, err := h.repository.Find(ctx, repoFilters)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) pushToStream(ticketId string, message string) error {
	return h.streamClient.Push(h.ctx, redisinfra.PushEvent{
		EventType:  redisinfra.PUSHEVENT_TICKET_SUBMITTED,
		StreamName: redisinfra.TICKET_CREATED,
		Values: map[string]any{
			"ticket_id": ticketId,
			"message":   message,
		}})
}

func (h *Handler) create(ctx *gin.Context) {
	var data CreateRequest
	if err := ctx.ShouldBindBodyWithJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	if data.Message == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Message must be defined"})
		return
	}

	result, err := h.repository.Create(ctx, DbCreateTicketInput{
		CustomerEmail:  data.CustomerEmail,
		Subject:        data.Subject,
		MessageContent: data.Message,
	})

	if err != nil {
		panic(err)
	}

	err = h.pushToStream(result.TicketId, data.Message)
	if err != nil {
		log.Printf("Failed to push ticketId=%s reason=%s", result.TicketId, err.Error())
	}

	ctx.JSON(http.StatusOK, CreateResponse{
		CustomerEmail: result.CustomerEmail,
		Messages: []MessageResponse{
			{
				Id:      result.MessageId,
				Content: result.MessageContent,
				Role:    result.Role,
			},
		},
		TicketId: result.TicketId,
	})
}

func toPaginationInput(page string) IPaginationInput {
	defaultInput := PaginationInput{
		Offset: 0,
		Limit:  5,
	}
	if page == "" {
		return defaultInput
	}
	pageAsInt, err := strconv.Atoi(page)
	if err != nil {
		log.Printf("Failed to convert page to int reason=%s", err.Error())
		return defaultInput
	}

	return PaginationInput{
		Limit:  pageAsInt,
		Offset: pageAsInt + 5,
	}
}

func (h *Handler) insertMessage(ctx *gin.Context) {

	var data struct {
		MessageContent string `json:"content"`
		CustomerEmail  string `json:"customerEmail"`
		Role           string `json:"role"`
	}
	if err := ctx.ShouldBindBodyWithJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ticketId := uuid.MustParse(ctx.Param("id")).String()

	if data.Role == "" {
		data.Role = "assistant"
	}
	dbInsertInput := DbInsertMessageInput{
		CustomerEmail:  data.CustomerEmail,
		TicketId:       ticketId,
		Role:           data.Role,
		MessageContent: data.MessageContent,
	}
	msg, err := h.repository.InsertMessage(ctx.Request.Context(), dbInsertInput)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = h.pushToStream(ticketId, msg.Content)
	if err != nil {
		log.Println(err.Error())
	}
	ctx.JSON(http.StatusOK, msg)
}

func (h *Handler) findById(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Invalid ticket id"})
		return
	}
	paginationInput := toPaginationInput(ctx.Query("message_page"))
	ticket, ok, err := h.repository.FindById(ctx.Request.Context(), FindByIdInput{
		TicketId:   id,
		Pagination: paginationInput,
	})
	if err != nil {
		panic(err)
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"errors": "Ticket not found"})
		return
	}

	ctx.JSON(http.StatusOK, ticket)
}
