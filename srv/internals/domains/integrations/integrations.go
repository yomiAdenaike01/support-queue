package integrationsdomain

import (
	"context"
	"encoding/json"

	"github.com/yomiAdenaike01/support-queue/internals/config"
)

type Platform string

const (
	PLATFORM_SLACK Platform = "SLACK"
)

type Integrations struct {
	Notifications *Notifications
}

type Notifications struct {
	config *config.Config
}

type NotificationInput struct {
	TargetId string
	Content  json.RawMessage
	Platform Platform
}

func (n *Notifications) SendNotification(ctx context.Context, input NotificationInput) error {
	switch input.Platform {
	case PLATFORM_SLACK:
		return sendSlackNotification(ctx, n.config, input)
	}
	return nil

}

func NewIntegrations(config *config.Config) *Integrations {
	return &Integrations{
		Notifications: &Notifications{
			config: config,
		},
	}
}
