package ticketdomain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	integrationsdomain "github.com/yomiAdenaike01/support-queue/internals/domains/integrations"
	teamdomain "github.com/yomiAdenaike01/support-queue/internals/domains/team"
	redisinfra "github.com/yomiAdenaike01/support-queue/internals/infra/redis"
	"github.com/yomiAdenaike01/support-queue/internals/web/sse"
)

type Handler struct {
	repository   *Repository
	streamClient *redisinfra.StreamClient
	teams        *teamdomain.Repository
	integrations *integrationsdomain.Integrations
	ctx          context.Context
	sseServer    *sse.Server
}

func NewHandler(ctx context.Context, db *sqlx.DB, streamClient *redisinfra.StreamClient, repository *Repository, teams *teamdomain.Repository, integrations *integrationsdomain.Integrations) *Handler {
	sseServer := sse.New(ctx)

	return &Handler{
		sseServer:    sseServer,
		repository:   repository,
		streamClient: streamClient,
		teams:        teams,
		integrations: integrations,
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
	group.POST("/tickets/:id/reprocess", h.reprocessTicket)
	group.PATCH("/tickets/:id", h.updateTicket)
	group.POST("/tickets/:id/pipeline", h.rerunPipeline)
	group.GET("/tickets/:id/events/sse", h.sseServer.Middleware())
}

type IdParam struct {
	Id string `uri:"id" binding:"required"`
}

func (h *Handler) rerunPipeline(ctx *gin.Context) {
	h.queueTicketForClassification(ctx, "Failed to rerun pipeline")
}

func (h *Handler) reprocessTicket(ctx *gin.Context) {
	h.queueTicketForClassification(ctx, "Failed to reprocess ticket")
}

func (h *Handler) queueTicketForClassification(ctx *gin.Context, failureMessage string) {
	var idParam IdParam
	if err := ctx.ShouldBindUri(&idParam); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Id missing from params",
		})
	}
	ticket, success, err := h.repository.FindById(ctx.Request.Context(), FindByIdInput{
		TicketId:   idParam.Id,
		Pagination: PaginationInput{},
	})
	if err != nil || success == false {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to find ticket",
		})
		return
	}

	if err = h.pushToStream(ctx.Request.Context(), idParam.Id, ticket.Messages[0].Content); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": failureMessage,
		})
		fmt.Println(err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}
func (h *Handler) updateTicket(ctx *gin.Context) {
	var params IdParam

	if err := ctx.ShouldBindUri(&params); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update ticket",
		})
		return
	}
	var input UpdateTicketInput
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		panic(err)
	}

	if err := h.repository.FindAndUpdate(ctx.Request.Context(), params.Id, input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update ticket",
		})
		return
	}
	if input.NotifyTeam {
		if err := h.notifyAssignedTeam(ctx.Request.Context(), params.Id); err != nil {
			log.Printf("Failed to notify assigned team ticket_id=%s reason=%s", params.Id, err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ticket updated but failed to notify team",
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{})

}

func (h *Handler) notifyAssignedTeam(ctx context.Context, ticketId string) error {
	if h.teams == nil || h.integrations == nil {
		return errors.New("team notifications are not configured")
	}
	ticket, ok, err := h.repository.FindById(ctx, FindByIdInput{
		TicketId:   ticketId,
		Pagination: PaginationInput{},
	})
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("ticket not found")
	}
	if ticket.AssignedTeam == nil || *ticket.AssignedTeam == "" {
		return errors.New("ticket has no assigned team")
	}
	teams, err := h.teams.Find(ctx, teamdomain.Filters{
		Departments: []string{*ticket.AssignedTeam},
	})
	if err != nil {
		return err
	}
	if len(teams) == 0 {
		return fmt.Errorf("assigned team %q not found", *ticket.AssignedTeam)
	}

	for _, team := range teams {
		if team.Integrations == nil {
			continue
		}
		var teamIntegrations teamdomain.TeamIntegrations
		if err := json.Unmarshal(*team.Integrations, &teamIntegrations); err != nil {
			return err
		}
		slack := teamIntegrations.GetSlack()
		if slack == nil || slack.GetChannelId() == nil {
			continue
		}
		channelId := *slack.GetChannelId()
		payload, err := json.Marshal(integrationsdomain.SlackNotificationData{
			Subject:           ticket.Subject,
			TicketID:          ticket.Id,
			TeamName:          team.Department,
			SuggestedResponse: nullableStringValue(ticket.SuggestedResponse),
			Priority:          nullableStringValue(ticket.Priority),
			Category:          nullableStringValue(ticket.Category),
			AssignedTeam:      nullableStringValue(ticket.AssignedTeam),
		})
		if err != nil {
			return err
		}
		if err := h.integrations.Notifications.SendNotification(ctx, integrationsdomain.NotificationInput{
			Platform: integrationsdomain.PLATFORM_SLACK,
			TargetId: channelId,
			Content:  payload,
		}); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("assigned team %q has no notification channel", *ticket.AssignedTeam)
}

func nullableStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
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

	h.sseServer.Emit(addEventBody.toJSON())
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
		f.Limit = 60
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

func (h *Handler) pushToStream(ctx context.Context, ticketId string, message string) error {
	return h.streamClient.Push(ctx, redisinfra.PushEvent{
		EventType:  redisinfra.PUSHEVENT_TICKET_SUBMITTED,
		StreamName: redisinfra.TICKET_CREATED,
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

	err = h.pushToStream(ctx.Request.Context(), result.TicketId, data.Message)
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

func newPaginationInput(page string) IPaginationInput {
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
	err = h.pushToStream(ctx.Request.Context(), ticketId, msg.Content)
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
	paginationInput := newPaginationInput(ctx.Query("message_page"))
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
