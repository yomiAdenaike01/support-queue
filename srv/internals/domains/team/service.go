package teamdomain

import (
	"context"

	integrationsdomain "github.com/yomiAdenaike01/support-queue/internals/domains/integrations"
)

type Service struct {
	repository   *Repository
	integrations *integrationsdomain.Integrations
}

func (h *Service) FindTeams(ctx context.Context, filters *Filters) ([]TeamResponse, error) {
	if filters == nil {
		filters = &Filters{}
	}
	return h.repository.Find(ctx, *filters)
}

func NewService(repository *Repository, integrations *integrationsdomain.Integrations) *Service {
	return &Service{
		repository:   repository,
		integrations: integrations,
	}
}
