package ticket

import "database/sql"

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
	Id       string  `json:"id"`
	Content  string  `json:"content"`
	Role     string  `json:"role"`
	TicketId *string `json:"ticket_id"`
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
}

type TicketResponse struct {
	Id                string            `json:"id"`
	CustomerEmail     string            `json:"customer_email"`
	Subject           string            `json:"subject"`
	Status            string            `json:"status"`
	Priority          *string           `json:"priority,omitempty"`
	Category          *string           `json:"category,omitempty"`
	AssignedTeam      *string           `json:"assigned_team,omitempty"`
	SuggestedResponse *string           `json:"suggested_response,omitempty"`
	Messages          []MessageResponse `json:"messages"`
	Pk                int               `json:"-"`
}
