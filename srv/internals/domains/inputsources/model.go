package inputsourcesdomain

import (
	"database/sql"
	"encoding/json"
	"time"
)

type InputSourceType string

const (
	InputSourceTypeEmail InputSourceType = "EMAIL"
)

type DbInputSource struct {
	CreatedAt       time.Time       `db:"created_at" json:"createdAt"`
	Id              string          `db:"id" json:"id"`
	Name            string          `db:"name" json:"name"`
	Pk              int             `db:"pk" json:"-"`
	Enabled         bool            `db:"enabled" json:"enabled"`
	Config          json.RawMessage `db:"config" json:"config"`
	SourceType      string          `db:"source_type" json:"sourceType"`
	TeamId          sql.NullString  `db:"team_id" json:"teamId"`
	TeamPk          sql.NullInt32   `db:"team_pk" json:"-"`
	Department      sql.NullString  `db:"department" json:"department"`
	ConnectionValue string          `db:"connection_value" json:"connectionValue"`
	Status          string          `db:"status" json:"status"`
}

type SaveInputSourceRequest struct {
	ConnectionValue string          `json:"connectionValue"`
	Status          string          `json:"status"`
	SourceType      string          `json:"sourceType"`
	Name            string          `json:"name"`
	Config          json.RawMessage `json:"config"`
	Enabled         bool            `json:"enabled"`
	TeamId          *string         `json:"teamId"`
}

type SaveInputSourceInput struct {
	Id              *string
	ConnectionValue string
	Status          string
	SourceType      string
	Name            string
	Config          json.RawMessage
	Enabled         bool
	TeamId          *string
}

func (s SaveInputSourceRequest) toInput(id *string) SaveInputSourceInput {
	return SaveInputSourceInput{
		Id:              id,
		ConnectionValue: s.ConnectionValue,
		Status:          s.Status,
		SourceType:      s.SourceType,
		Name:            s.Name,
		Config:          s.Config,
		Enabled:         s.Enabled,
		TeamId:          s.TeamId,
	}
}

type InputSource struct {
	CreatedAt       time.Time       `json:"createdAt"`
	Id              string          `json:"id"`
	Name            string          `json:"name"`
	Enabled         bool            `json:"enabled"`
	Config          json.RawMessage `json:"config"`
	SourceType      string          `json:"sourceType"`
	TeamId          *string         `json:"teamId"`
	Department      *string         `json:"department"`
	ConnectionValue string          `json:"connectionValue"`
	Status          string          `json:"status"`
}

func toStringPtr(str sql.NullString) *string {
	if str.Valid {
		return &str.String
	}
	return nil
}

type DbInputSources []DbInputSource

func (d DbInputSources) toInputSources() ([]InputSource, error) {
	sources := make([]InputSource, 0, len(d))
	for _, src := range d {
		sources = append(sources, src.toInputSource())
	}
	return sources, nil
}
func (d DbInputSource) toInputSource() InputSource {
	return InputSource{
		CreatedAt:       d.CreatedAt,
		Id:              d.Id,
		Name:            d.Name,
		Enabled:         d.Enabled,
		Config:          d.Config,
		SourceType:      d.SourceType,
		TeamId:          toStringPtr(d.TeamId),
		Department:      toStringPtr(d.Department),
		ConnectionValue: d.ConnectionValue,
		Status:          d.Status,
	}
}

type Credentials struct {
	Principal string `json:"principal"`
	Password  string `json:"password"`
}

func (i InputSource) GetCredentials() (Credentials, error) {
	var credentials Credentials
	return credentials, json.Unmarshal(i.Config, &credentials)
}

type FindInput struct {
	SourceType *InputSourceType
	Id         *string
	Department *string
	Limit      *int
}

func (f FindInput) toManyInputArgs() []any {
	return []any{f.Department, f.SourceType, f.Id, 10}
}

func (f FindInput) toOneInputArgs() []any {
	return []any{f.Department, f.SourceType, f.Id, 1}
}

func (s SaveInputSourceInput) createArgs() []any {
	return []any{s.ConnectionValue, s.Status, s.SourceType, s.Name, s.Config, s.TeamId}
}

func (s SaveInputSourceInput) updateArgs() []any {
	return []any{s.Id, s.ConnectionValue, s.Status, s.SourceType, s.Name, s.Config, s.Enabled, s.TeamId}
}
