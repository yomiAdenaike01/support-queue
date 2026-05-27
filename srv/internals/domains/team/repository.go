package teamdomain

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"maps"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/yomiAdenaike01/support-queue/internals/db/queries"
)

type TeamIntegrations map[string]any
type SlackIntegration map[string]interface{}

func (t TeamIntegrations) GetSlack() *SlackIntegration {
	slackIntegrationRecord, ok := t["slack"]
	if !ok {
		log.Println("Failed to get slack key on integrations")
		return nil
	}

	slackIntegrationMap, ok := slackIntegrationRecord.(map[string]interface{})
	if !ok {
		log.Println("Failed to cast integration to map[string]interface{}")
		return nil
	}

	slackIntegration := SlackIntegration{}
	maps.Copy(slackIntegration, slackIntegrationMap)
	return &slackIntegration
}

func (s SlackIntegration) GetChannelId() *string {
	rawChannelId, ok := s["channel_id"]
	if !ok {
		return nil
	}
	channelId, ok := rawChannelId.(string)
	if !ok {
		log.Println("Failed to cast channel id to string")
		return nil
	}
	return &channelId

}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, input CreateTeamRequest) (TeamResponse, bool, error) {
	roles := make([]string, 0, len(input.Members))
	names := make([]string, 0, len(input.Members))
	phoneNumbers := make([]string, 0, len(input.Members))
	emailAddresses := make([]string, 0, len(input.Members))

	for _, member := range input.Members {
		roles = append(roles, member.Role)
		names = append(names, member.Name)
		phoneNumbers = append(phoneNumbers, member.PhoneNumber)
		emailAddresses = append(emailAddresses, member.EmailAddress)
	}

	var rows []DbCreateTeamRow
	err := r.db.SelectContext(
		ctx,
		&rows,
		queries.CREATE_TEAM_QUERY,
		input.Department,
		pq.Array(roles),
		pq.Array(names),
		pq.Array(phoneNumbers),
		pq.Array(emailAddresses),
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return TeamResponse{}, false, nil
		}
		return TeamResponse{}, false, err
	}
	if len(rows) == 0 {
		return TeamResponse{}, false, nil
	}

	return buildTeamResponse(rows), true, nil
}

func buildTeamResponse(rows []DbCreateTeamRow) TeamResponse {
	if len(rows) == 0 {
		return TeamResponse{}
	}

	first := rows[0]
	response := TeamResponse{
		Id:           first.TeamId,
		Department:   first.Department,
		Integrations: nullableJSONPtr(first.Integrations),
		Members:      make([]TeamMemberResponse, 0, len(rows)),
	}

	for _, row := range rows {
		if !row.MemberId.Valid {
			continue
		}

		response.Members = append(response.Members, TeamMemberResponse{
			Id:           row.MemberId.String,
			Name:         row.Name.String,
			Role:         row.Role.String,
			EmailAddress: nullableStringPtr(row.EmailAddress),
			PhoneNumber:  nullableStringPtr(row.PhoneNumber),
		})
	}

	return response
}

func nullableStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func nullableJSONPtr(value sql.NullString) *json.RawMessage {
	if !value.Valid {
		return nil
	}
	raw := json.RawMessage(value.String)
	return &raw
}

type Filters struct {
	Id          string
	Departments []string
}

func (f Filters) GetId() *string {
	if f.Id == "" {
		return nil
	}
	return &f.Id
}

func (f Filters) GetDepartments() []string {
	departments := make([]string, 0, len(f.Departments))
	for _, department := range f.Departments {
		department = strings.TrimSpace(department)
		if department == "" {
			continue
		}
		departments = append(departments, department)
	}
	if len(departments) == 0 {
		return nil
	}
	return departments
}

func (r *Repository) Find(ctx context.Context, filters Filters) ([]TeamResponse, error) {
	var rawTeams []DbCreateTeamRow
	var teamResponse []TeamResponse
	if err := r.db.SelectContext(ctx, &rawTeams, queries.SEARCH_TEAM_QUERY, filters.GetId(), pq.Array(filters.GetDepartments())); err != nil {
		return teamResponse, err
	}
	if len(rawTeams) == 0 {
		return teamResponse, nil
	}
	if len(rawTeams) > 1 {
		var teamsList = make([]TeamResponse, 0, len(rawTeams))
		for _, rawTeam := range rawTeams {
			team := buildTeamResponse([]DbCreateTeamRow{rawTeam})
			teamsList = append(teamsList, team)
		}
		return teamsList, nil
	}

	return []TeamResponse{buildTeamResponse(rawTeams)}, nil
}
