package teamdomain

import (
	"database/sql"
	"encoding/json"
)

type CreateTeamMemberRequest struct {
	Name         string `json:"name"`
	Role         string `json:"role"`
	EmailAddress string `json:"email_address"`
	PhoneNumber  string `json:"phone_number"`
}

type CreateTeamRequest struct {
	Department string                    `json:"department"`
	Members    []CreateTeamMemberRequest `json:"members"`
}

type DbCreateTeamRow struct {
	TeamId       string         `db:"team_id"`
	TeamPk       int            `db:"team_pk"`
	Department   string         `db:"department"`
	Integrations sql.NullString `db:"integrations"`
	MemberPk     sql.NullInt64  `db:"member_pk"`
	MemberId     sql.NullString `db:"member_id"`
	EmailAddress sql.NullString `db:"email_address"`
	PhoneNumber  sql.NullString `db:"phone_number"`
	Name         sql.NullString `db:"name"`
	Role         sql.NullString `db:"role"`
}

type TeamMemberResponse struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Role         string  `json:"role"`
	EmailAddress *string `json:"email_address,omitempty"`
	PhoneNumber  *string `json:"phone_number,omitempty"`
}

type TeamResponse struct {
	Id           string               `json:"id"`
	Department   string               `json:"department"`
	Integrations *json.RawMessage     `json:"integrations,omitempty"`
	Members      []TeamMemberResponse `json:"members"`
}
