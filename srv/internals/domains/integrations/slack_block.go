package integrationsdomain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/yomiAdenaike01/support-queue/internals/config"
)

const slackTicketNotification = `[
	{
		"type": "card",
		"title": {
			"type": "mrkdwn",
			"text": "Priority - {{priorityEmoji .Priority}} - {{ .Priority }} {{.Subject}}",
			"verbatim": false
		},
		"subtitle": {
			"type": "mrkdwn",
			"text": "{{.TeamName}} · Ticket #{{.TicketID}}",
			"verbatim": false
		},
		"body": {
			"type": "mrkdwn",
			"text": "We suggest: {{.SuggestedResponse}}",
			"verbatim": false
		},
		"actions": [
			{
				"type": "button",
				"text": {
					"type": "plain_text",
					"text": "View Ticket",
					"emoji": false
				},
				"style": "primary",
				"action_id": "view_ticket_{{.TicketID}}"
			},
			{
				"type": "button",
				"text": {
					"type": "plain_text",
					"text": "Resolve",
					"emoji": false
				},
				"action_id": "resolve_ticket_{{.TicketID}}"
			}
		]
	},
	{
		"type": "context",
		"elements": [
			{
				"type": "mrkdwn",
				"text": "{{.Priority}} · {{.Category}} · Assigned to {{.AssignedTeam}} · <{{.DashboardURL}}/tickets/{{.TicketID}}|View in Dashboard>"
			}
		]
	}
]`

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
	funcMap := template.FuncMap{
		"priorityEmoji": func(p string) string {
			switch p {
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
		},
	}
	tmpl, err := template.New("slack").Funcs(funcMap).Parse(slackTicketNotification)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	var slackNotificationData SlackNotificationData

	if err := json.Unmarshal(input.Content, &slackNotificationData); err != nil {
		return err
	}

	if err := tmpl.Execute(&buf, slackNotificationData); err != nil {
		return err
	}

	body, err := json.Marshal(map[string]any{
		"channel": input.TargetId,
		"blocks":  json.RawMessage(buf.Bytes()),
	})
	if err != nil {
		return err
	}
	reader := bytes.NewReader(body)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://slack.com/api/chat.postMessage", reader)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetEnvOrFail("SLACK_BOT_TOKEN")))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("Failed send notification to url=%s method=%s status:%d", req.URL.String(), req.Method, response.StatusCode)
	}
	return nil
}
