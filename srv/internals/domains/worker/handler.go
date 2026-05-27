package workerdomain

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	integrationsdomain "github.com/yomiAdenaike01/support-queue/internals/domains/integrations"
	teamdomain "github.com/yomiAdenaike01/support-queue/internals/domains/team"
)

type Handler struct {
	integrations    *integrationsdomain.Integrations
	teamsRepository *teamdomain.Repository
}

type WorkerResult struct {
	SuggestedResponse string  `json:"suggested_response"`
	AverageSentiment  float32 `json:"average_sentiment"`
	Priority          string  `json:"priority"`
	Category          string  `json:"category"`
	RequiresUrgency   bool    `json:"requires_urgency"`
	UrgencyReason     string  `json:"urgency_reason"`
}

func (w *Handler) completeWork(ctx *gin.Context) {
	var workerResult WorkerResult
	if err := ctx.ShouldBindBodyWithJSON(&workerResult); err != nil {
		panic(err)
	}
	teamDepartments := fromCategoryToDepartments(workerResult.Category)
	teams, err := w.teamsRepository.Find(ctx.Request.Context(), teamdomain.Filters{
		Departments: teamDepartments,
	})
	if err != nil {
		panic(err)
	}
	w.notifyTeams(ctx.Request.Context(), teams, workerResult)
}

func (w *Handler) notifyTeams(ctx context.Context, teams []teamdomain.TeamResponse, message WorkerResult) []error {
	errors := make([]error, 0, len(teams)*3)
	for _, team := range teams {
		var teamIntegrations teamdomain.TeamIntegrations

		if err := json.Unmarshal(*team.Integrations, &teamIntegrations); err != nil {
			panic(err)
		}
		slack := teamIntegrations.GetSlack()
		if slack == nil {
			return nil
		}
		rawMessage, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}
		err = w.integrations.Notifications.SendNotification(ctx, integrationsdomain.NotificationInput{
			Platform: integrationsdomain.PLATFORM_SLACK,
			TargetId: *slack.GetChannelId(),
			Content:  rawMessage,
		})
		if err != nil {
			log.Println(err.Error())
			errors = append(errors, err)
		}
	}
	return errors

}

func (w *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", w.completeWork)
}

var CATEGORY_TO_TEAM = map[string][]string{
	"BILLING":       {"billing", "finance"},
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

func NewHandler(integrations *integrationsdomain.Integrations, teamRepository *teamdomain.Repository) *Handler {
	return &Handler{
		integrations:    integrations,
		teamsRepository: teamRepository,
	}
}
