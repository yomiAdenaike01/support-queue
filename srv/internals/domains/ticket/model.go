package ticketdomain

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type TicketEvent struct {
	EventType string          `json:"eventType"`
	CreatedAt string          `json:"createdAt"`
	Payload   json.RawMessage `json:"payload"`
}

type TicketEvents []TicketEvent

func (t *TicketEvents) Scan(value interface{}) error {
	if value == nil {
		*t = TicketEvents{}
		return nil
	}
	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("Failed to cast value to []byte")
	}
	return json.Unmarshal(bytes, t)
}

type DbCreateTicketResult struct {
	TicketId       string `db:"ticket_id"`
	TicketPk       int    `db:"ticket_pk"`
	CustomerEmail  string `db:"customer_email"`
	MessagePk      int    `db:"message_pk"`
	MessageId      string `db:"message_id"`
	MessageContent string `db:"message_content"`
	Role           string `db:"role"`
}

type MessageResponse struct {
	Id        string     `json:"id"`
	Content   string     `json:"content"`
	Role      string     `json:"role"`
	TicketId  *string    `json:"ticket_id"`
	CreatedAt *time.Time `json:"createdAt"`
}

type CreateResponse struct {
	Messages      []MessageResponse `json:"messages"`
	TicketId      string            `json:"id"`
	CustomerEmail string            `json:"customer_email"`
}

type CreateRequest struct {
	Message       string `json:"message"`
	Subject       string `json:"subject"`
	CustomerEmail string `json:"customerEmail"`
}

type DbFindTicketByIdRow struct {
	TicketId          string         `db:"ticket_id"`
	TicketPk          int            `db:"ticket_pk"`
	CustomerEmail     string         `db:"customer_email"`
	Subject           string         `db:"subject"`
	Status            string         `db:"status"`
	Priority          sql.NullString `db:"priority"`
	Category          sql.NullString `db:"category"`
	AssignedTeam      sql.NullString `db:"assigned_team"`
	SuggestedResponse sql.NullString `db:"suggested_response"`
	MessageId         sql.NullString `db:"message_id"`
	MessageContent    sql.NullString `db:"message_content"`
	MessageRole       sql.NullString `db:"message_role"`
	MessageCreatedAt  sql.NullTime   `db:"message_created_at"`
	Events            TicketEvents   `db:"events"`
	CreatedAt         sql.NullTime   `db:"created_at"`
}

type TicketResponse struct {
	Id                string            `json:"id"`
	CustomerEmail     string            `json:"customerEmail"`
	Subject           string            `json:"subject"`
	Status            string            `json:"status"`
	Priority          *string           `json:"priority,omitempty"`
	Category          *string           `json:"category,omitempty"`
	AssignedTeam      *string           `json:"assignedTeam,omitempty"`
	SuggestedResponse *string           `json:"suggestedResponse,omitempty"`
	Messages          []MessageResponse `json:"messages,omitempty"`
	Pk                int               `json:"-"`
	Events            TicketEvents      `json:"events,omitempty"`
	CreatedAt         *time.Time        `json:"createdAt"`
}
