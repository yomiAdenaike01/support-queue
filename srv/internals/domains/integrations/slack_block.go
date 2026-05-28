package integrationsdomain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yomiAdenaike01/support-queue/internals/config"
)

type SlackResponse struct {
	Ok               bool                  `json:"ok"`
	Error            string                `json:"error"`
	Errors           []any                 `json:"errors"`
	Warning          string                `json:"warning"`
	ResponseMetadata SlackResponseMetadata `json:"response_metadata"`
}

type SlackResponseMetadata struct {
	Messages []string `json:"messages"`
	Warnings []string `json:"warnings"`
}

type SlackNotificationData struct {
	Emoji             string `json:"emoji"` // 🔴 🟡 🟢 based on priority
	Subject           string `json:"subject"`
	TicketID          string `json:"ticket_id"`
	TeamName          string `json:"team_name"`
	SuggestedResponse string `json:"suggested_response"`
	Priority          string `json:"priority"`
	Category          string `json:"category"`
	AssignedTeam      string `json:"assigned_team"`
	DashboardURL      string `json:"dashboard_url"`
}

func sendSlackNotification(ctx context.Context, config *config.Config, input NotificationInput) error {
	var slackNotificationData SlackNotificationData

	if err := json.Unmarshal(input.Content, &slackNotificationData); err != nil {
		return err
	}

	body, err := json.Marshal(map[string]any{
		"channel": input.TargetId,
		"blocks":  buildSlackBlocks(slackNotificationData),
	})
	if err != nil {
		return err
	}
	reader := bytes.NewReader(body)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://slack.com/api/chat.postMessage", reader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetEnvOrFail("SLACK_BOT_TOKEN")))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var responseData SlackResponse
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("failed send notification to url=%s method=%s status=%d slack_error=%s", req.URL.String(), req.Method, response.StatusCode, responseData.Error)
	}
	if !responseData.Ok {
		return fmt.Errorf("slack api error=%s errors=%v messages=%v warnings=%v", responseData.Error, responseData.Errors, responseData.ResponseMetadata.Messages, responseData.ResponseMetadata.Warnings)
	}
	return nil
}

func buildSlackBlocks(data SlackNotificationData) []map[string]any {
	headerText := fmt.Sprintf("*Priority - %s - %s*", priorityEmoji(data.Priority), data.Priority)
	if data.Subject != "" {
		headerText = fmt.Sprintf("%s %s", headerText, data.Subject)
	}

	contextText := fmt.Sprintf("%s · %s", data.Priority, data.Category)
	if data.AssignedTeam != "" {
		contextText = fmt.Sprintf("%s · Assigned to %s", contextText, data.AssignedTeam)
	}
	if data.DashboardURL != "" && data.TicketID != "" {
		contextText = fmt.Sprintf("%s · <%s/tickets/%s|View in Dashboard>", contextText, data.DashboardURL, data.TicketID)
	}

	return []map[string]any{
		{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": headerText,
			},
		},
		{
			"type": "section",
			"text": map[string]any{
				"type": "mrkdwn",
				"text": fmt.Sprintf("We suggest: %s", data.SuggestedResponse),
			},
		},
		{
			"type": "actions",
			"elements": []map[string]any{
				{
					"type": "button",
					"text": map[string]any{
						"type":  "plain_text",
						"text":  "View Ticket",
						"emoji": false,
					},
					"style":     "primary",
					"action_id": fmt.Sprintf("view_ticket_%s", data.TicketID),
				},
				{
					"type": "button",
					"text": map[string]any{
						"type":  "plain_text",
						"text":  "Resolve",
						"emoji": false,
					},
					"action_id": fmt.Sprintf("resolve_ticket_%s", data.TicketID),
				},
			},
		},
		{
			"type": "context",
			"elements": []map[string]any{
				{
					"type": "mrkdwn",
					"text": contextText,
				},
			},
		},
	}
}

func priorityEmoji(priority string) string {
	switch priority {
	case "URGENT":
		return "🔴"
	case "HIGH":
		return "🔴"
	case "MEDIUM":
		return "🟡"
	case "LOW":
		return "🟢"
	}
	return ""
}
