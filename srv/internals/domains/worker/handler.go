package workerdomain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	integrationsdomain "github.com/yomiAdenaike01/support-queue/internals/domains/integrations"
	teamdomain "github.com/yomiAdenaike01/support-queue/internals/domains/team"
	"github.com/yomiAdenaike01/support-queue/internals/domains/ticket"
)

type Handler struct {
	context          context.Context
	integrations     *integrationsdomain.Integrations
	teamsRepository  *teamdomain.Repository
	ticketRepository *ticket.Repository
}

type WorkerResult struct {
	TicketId          string  `json:"ticket_id"`
	SuggestedResponse string  `json:"suggested_response"`
	AverageSentiment  float32 `json:"average_sentiment"`
	Priority          string  `json:"priority"`
	Category          string  `json:"category"`
	RequiresUrgency   bool    `json:"requires_urgency"`
	UrgencyReason     string  `json:"urgency_reason"`
}

func workerResultToUpdateInput(workerResult WorkerResult) ticket.UpdateTicketInput {
	return ticket.UpdateTicketInput{
		SuggestedResponse: &workerResult.SuggestedResponse,
		Category:          &workerResult.Category,
		RequiresUrgency:   &workerResult.RequiresUrgency,
		UrgencyReason:     &workerResult.UrgencyReason,
		AverageSentiment:  &workerResult.AverageSentiment,
		Priority:          &workerResult.Priority,
	}
}

func (w *Handler) completeWork(ctx *gin.Context) {
	var workerResult WorkerResult
	if err := ctx.ShouldBindBodyWithJSON(&workerResult); err != nil {
		log.Printf("Failed to bind worker result reason=%s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := w.ticketRepository.FindAndUpdate(ctx.Request.Context(), workerResult.TicketId, workerResultToUpdateInput(workerResult))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	log.Printf(
		"Received worker result category=%s priority=%s requires_urgency=%t average_sentiment=%.2f",
		workerResult.Category,
		workerResult.Priority,
		workerResult.RequiresUrgency,
		workerResult.AverageSentiment,
	)

	teamDepartments := fromCategoryToDepartments(workerResult.Category)
	if len(teamDepartments) == 0 {
		log.Printf("No departments mapped for worker result category=%s", workerResult.Category)
	}

	teams, err := w.teamsRepository.Find(ctx.Request.Context(), teamdomain.Filters{
		Departments: teamDepartments,
	})
	if err != nil {
		log.Printf("Failed to find teams for worker result category=%s departments=%v reason=%s", workerResult.Category, teamDepartments, err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Printf("Found teams for worker result category=%s departments=%v team_count=%d", workerResult.Category, teamDepartments, len(teams))

	go func(w *Handler, teams []teamdomain.TeamResponse, workerResult WorkerResult) {
		log.Printf("Starting team notifications category=%s priority=%s team_count=%d", workerResult.Category, workerResult.Priority, len(teams))
		err := w.notifyTeams(w.context, teams, workerResult)
		if len(err) == 0 {
			log.Printf("Completed team notifications category=%s priority=%s team_count=%d", workerResult.Category, workerResult.Priority, len(teams))
			return
		}
		log.Printf("Completed team notifications with errors category=%s priority=%s team_count=%d error_count=%d", workerResult.Category, workerResult.Priority, len(teams), len(err))

	}(w, teams, workerResult)
	ctx.Status(200)

}

func (w *Handler) notifyTeams(ctx context.Context, teams []teamdomain.TeamResponse, message WorkerResult) []error {
	errors := make([]error, 0, len(teams)*3)
	if len(teams) == 0 {
		log.Printf("No teams to notify category=%s priority=%s", message.Category, message.Priority)
		return errors
	}

	for _, team := range teams {
		var teamIntegrations teamdomain.TeamIntegrations

		if team.Integrations == nil {
			log.Printf("Skipping team notification team_id=%s department=%s reason=missing_integrations", team.Id, team.Department)
			continue
		}
		if err := json.Unmarshal(*team.Integrations, &teamIntegrations); err != nil {
			log.Printf("Failed to decode team integrations team_id=%s department=%s reason=%s", team.Id, team.Department, err.Error())
			errors = append(errors, fmt.Errorf("team_id=%s department=%s decode_integrations: %w", team.Id, team.Department, err))
			continue
		}
		slack := teamIntegrations.GetSlack()
		if slack == nil {
			log.Printf("Skipping team notification team_id=%s department=%s reason=missing_slack_integration", team.Id, team.Department)
			continue
		}
		rawMessage, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to encode worker result for notification team_id=%s department=%s reason=%s", team.Id, team.Department, err.Error())
			errors = append(errors, fmt.Errorf("team_id=%s department=%s encode_notification: %w", team.Id, team.Department, err))
			continue
		}
		channelId := slack.GetChannelId()
		if channelId == nil {
			log.Printf("Skipping team notification team_id=%s department=%s reason=missing_slack_channel_id", team.Id, team.Department)
			continue
		}
		log.Printf("Sending team notification team_id=%s department=%s platform=%s target_id=%s", team.Id, team.Department, integrationsdomain.PLATFORM_SLACK, *channelId)
		err = w.integrations.Notifications.SendNotification(ctx, integrationsdomain.NotificationInput{
			Platform: integrationsdomain.PLATFORM_SLACK,
			TargetId: *channelId,
			Content:  rawMessage,
		})
		if err != nil {
			log.Printf("Failed to send team notification team_id=%s department=%s platform=%s target_id=%s reason=%s", team.Id, team.Department, integrationsdomain.PLATFORM_SLACK, *channelId, err.Error())
			errors = append(errors, err)
			continue
		}
		log.Printf("Sent team notification team_id=%s department=%s platform=%s target_id=%s", team.Id, team.Department, integrationsdomain.PLATFORM_SLACK, *channelId)
	}
	return errors

}

func (w *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("", w.completeWork)
	rg.GET("/health", w.serviceHealth)
}

func (w *Handler) serviceHealth(ctx *gin.Context) {

}

var CATEGORY_TO_TEAM = map[string][]string{
	"BILLING":       {"billing", "finance"},
	"CANCELLATION":  {"billing", "retention"},
	"TECHNICAL":     {"tech-support"},
	"DELIVERY":      {"logistics"},
	"SUBSCRIPTIONS": {"retention", "customer-success"},
	"GENERAL":       {"customer-success", "tech-support", "billing"},
}

func fromCategoryToDepartments(category string) []string {
	teams, ok := CATEGORY_TO_TEAM[category]
	if !ok {
		return nil
	}
	return teams
}

func NewHandler(context context.Context, integrations *integrationsdomain.Integrations, teamRepository *teamdomain.Repository, ticketRepository *ticket.Repository) *Handler {
	return &Handler{
		context:          context,
		integrations:     integrations,
		teamsRepository:  teamRepository,
		ticketRepository: ticketRepository,
	}
}
