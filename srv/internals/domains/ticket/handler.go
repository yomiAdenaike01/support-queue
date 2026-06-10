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
}

type AddEventParam struct {
	Id string `uri:"id" binding:"required"`
}

func (h *Handler) addEvent(ctx *gin.Context) {
	var params AddEventParam
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
		if errors.Is(err, ErrDuplicateTicketEvent) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, event)
}
func (h *Handler) find(ctx *gin.Context) {
	search := ctx.Query("search")
	status := ctx.Query("status")
	priority := ctx.Query("priority")
	category := ctx.Query("category")
	limit := ctx.Query("limit")

	limitint, err := strconv.Atoi(limit)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := h.repository.Find(ctx, FindFilters{
		Search:   search,
		Status:   status,
		Priority: priority,
		Category: category,
		Limit:    limitint,
	})

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
		panic(err)
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
