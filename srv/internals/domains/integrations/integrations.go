package integrationsdomain

import (
	"context"
	"encoding/json"

	"github.com/yomiAdenaike01/support-queue/internals/config"
	inputsourcesdomain "github.com/yomiAdenaike01/support-queue/internals/domains/inputsources"
)

type Platform string

const (
	PLATFORM_SLACK Platform = "SLACK"
)

type Integrations struct {
	Notifications    *Notifications
	inputSourcesRepo inputsourcesdomain.Repository
	imapClient       *ImapClient
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

func New(config *config.Config, inputSourcesRepo inputsourcesdomain.Repository) *Integrations {
	notifications := &Notifications{
		config: config,
	}

	integrations := &Integrations{
		inputSourcesRepo: inputSourcesRepo,
		Notifications:    notifications,
		imapClient:       newImapClient(context.TODO(), inputSourcesRepo, make(chan any)),
	}
	return integrations
}
