package integrationsdomain

import (
	"context"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/yomiAdenaike01/support-queue/internals/config"
	inputsourcesdomain "github.com/yomiAdenaike01/support-queue/internals/domains/inputsources"
	"github.com/yomiAdenaike01/support-queue/internals/domains/integrations/oauth"
)

type Platform string

const (
	PLATFORM_SLACK Platform = "SLACK"
)

type Integrations struct {
	Notifications       *Notifications
	inputSourcesRepo    inputsourcesdomain.Repository
	gmailClient         *GmailClient
	authIntegrationRepo oauth.AuthRepository
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

type Deps struct {
	Db               *sqlx.DB
	Config           *config.Config
	InputSourcesRepo inputsourcesdomain.Repository
	IngestionChanel  chan<- TicketPayload
	Providers        *oauth.Providers
}

func New(deps Deps) *Integrations {
	notifications := &Notifications{
		config: deps.Config,
	}
	authIntegrationRepo := oauth.NewCredentialsRepository(deps.Db)

	emailIntegrationDep := emailIntegrationDeps{
		repository:     deps.InputSourcesRepo,
		onNewEmail:     deps.IngestionChanel,
		oauthProviders: deps.Providers,
	}

	integrations := &Integrations{
		inputSourcesRepo:    deps.InputSourcesRepo,
		Notifications:       notifications,
		authIntegrationRepo: authIntegrationRepo,
		gmailClient:         newGmailClient(context.TODO(), emailIntegrationDep),
	}
	return integrations
}
