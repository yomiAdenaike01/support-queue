package ticket

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/yomiAdenaike01/support-queue/internals/db/queries"
)

type Handler struct {
	db         *sqlx.DB
	repository *Repository
	redis      *redis.Client
	streamName string
	ctx        context.Context
}

func NewHandler(ctx context.Context, db *sqlx.DB, redisClient *redis.Client, streamName string) *Handler {
	return &Handler{
		db:         db,
		repository: NewRepository(db),
		redis:      redisClient,
		streamName: streamName,
		ctx:        ctx,
	}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/support/ticket", h.create)
	group.GET("/support/ticket/:id", h.findById)
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

	var createdTicket DbCreateTicketResult
	query, args, err := sqlx.Named(queries.CREATE_TICKET_QUERY, map[string]interface{}{
		"customer_email":  data.CustomerEmail,
		"subject":         data.Subject,
		"message_content": data.Message,
	})
	if err != nil {
		panic(err)
	}
	query = h.db.Rebind(query)

	if err := h.db.Get(&createdTicket, query, args...); err != nil {
		panic(err)
	}

	result, err := h.redis.XAdd(h.ctx, &redis.XAddArgs{
		Stream: h.streamName,
		Values: map[string]any{
			"ticket_id": 3,
			"message":   data.Message,
		},
	}).Result()
	if err != nil {
		panic(err)
	}
	log.Println(result)

	ctx.JSON(http.StatusOK, CreateResponse{
		CustomerEmail: createdTicket.CustomerEmail,
		Messages: []MessageResponse{
			{
				Id:      createdTicket.MessageId,
				Content: createdTicket.MessageContent,
				Role:    createdTicket.Role,
			},
		},
		TicketId: createdTicket.TicketId,
	})
}

func (h *Handler) findById(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": "Invalid ticket id"})
		return
	}

	ticket, ok, err := h.repository.FindById(ctx.Request.Context(), id)
	if err != nil {
		panic(err)
	}
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"errors": "Ticket not found"})
		return
	}

	ctx.JSON(http.StatusOK, ticket)
}
