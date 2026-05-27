package ticket

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	redisinfra "github.com/yomiAdenaike01/support-queue/internals/infra/redis"
)

type Handler struct {
	db           *sqlx.DB
	repository   *Repository
	streamClient *redisinfra.StreamClient
	ctx          context.Context
}

func NewHandler(ctx context.Context, db *sqlx.DB, streamClient *redisinfra.StreamClient) *Handler {
	return &Handler{
		db:           db,
		repository:   NewRepository(db),
		streamClient: streamClient,
		ctx:          ctx,
	}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/ticket", h.create)
	group.GET("/ticket/:id", h.findById)
	group.POST("/ticket/:id/message", h.insertMessage)
}

func (h *Handler) pushToStream(ticketId string, message string) error {
	return h.streamClient.Push(h.ctx, redisinfra.PushEvent{
		EventType: redisinfra.PUSHEVENT_TICKET_SUBMITTED,
		Values: map[string]interface{}{
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
