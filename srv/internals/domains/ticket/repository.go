package ticket

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yomiAdenaike01/support-queue/internals/db/queries"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

type PaginationInput struct {
	Offset int
	Limit  int
}

func (p PaginationInput) GetOffset() int {
	return p.Offset
}

func (p PaginationInput) GetLimit() int {
	if p.Limit == 0 {
		return 5
	}
	return p.Limit
}

type IPaginationInput interface {
	GetOffset() int
	GetLimit() int
}

type FindByIdInput struct {
	TicketId   string
	Pagination IPaginationInput
}

type DbCreateTicketInput struct {
	CustomerEmail  string `db:"customer_email"`
	Subject        string `db:"subject"`
	MessageContent string `db:"message_content"`
}

type DbInsertMessageInput struct {
	MessageContent string `db:"content"`
	TicketId       string `db:"ticket_id"`
	CustomerEmail  string `db:"customer_email"`
	Role           string `db:"role"`
}
type DbCreateMessageResult struct {
	Id       string `db:"id" json:"id"`
	Content  string `db:"content" json:"content"`
	TicketId string `db:"ticket_id" json:"ticketId"`
	Role     string `db:"role" json:"role"`
	TicketPk int    `db:"ticket_pk" json:"-"`
}

type TicketEventType string

type FindFilters struct {
	Search   string
	Status   string
	Priority string
	Category string
	Limit    int
}

type CreateEventInput struct {
	TicketId string
	Status   *string         `json:"status" binding:"required"`
	Payload  json.RawMessage `json:"payload"`
}

type DbEventResponse struct {
	EventId   string         `json:"eventId" db:"event_id"`
	EventPk   int            `json:"-" db:"event_pk"`
	TicketId  string         `json:"ticketId" db:"ticket_id"`
	TicketPK  int            `json:"-" db:"ticket_pk"`
	Payload   sql.NullString `json:"payload" db:"payload"`
	EventType string         `json:"status" db:"event_type"`
	CreatedAt time.Time      `json:"createdAt" db:"created_at"`
}

type EventResponse struct {
	DbEventResponse
	Payload json.RawMessage `json:"payload"`
}

var ErrDuplicateTicketEvent = errors.New("ticket event already exists")

const (
	TICKET_CREATED     TicketEventType = "TICKET_CREATED"
	TICKET_PROCESSING  TicketEventType = "TICKET_PROCESSING"
	TICKET_CLASSIFIED  TicketEventType = "TICKET_CLASSIFIED"
	TICKET_PROCESSED   TicketEventType = "TICKET_PROCESSED"
	TICKET_FAILED      TicketEventType = "TICKET_FAILED"
	TICKET_REPROCESSED TicketEventType = "TICKET_REPROCESSED"
	TICKET_RESOLVED    TicketEventType = "TICKET_RESOLVED"
)

func nullableStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func emptyStrOrNil(s string) *string {
	if len(s) == 0 || s == "ALL" {
		return nil
	}
	return &s
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (f FindFilters) ToArgs() []any {
	return []any{
		emptyStrOrNil(f.Search),
		emptyStrOrNil(f.Status),
		emptyStrOrNil(f.Priority),
		emptyStrOrNil(f.Category),
		f.Limit,
	}
}

func nullableTimePtr(nullableTime sql.NullTime) *time.Time {
	if nullableTime.Valid {
		return &nullableTime.Time
	}
	return nil
}

type UpdateTicketInput struct {
	SuggestedResponse *string
	AverageSentiment  *float32
	Priority          *string
	Category          *string
	RequiresUrgency   *bool
	UrgencyReason     *string
	AssignedTeam      *string
}

func (r *Repository) FindAndUpdate(ctx context.Context, id string, update UpdateTicketInput) error {
	res, err := r.db.ExecContext(ctx, `UPDATE tickets 
	 SET suggested_response = COALESCE($1::TEXT,suggested_response),
	 	average_sentiment = COALESCE($2::INT, average_sentiment),
		priority = COALESCE($3::TEXT, priority),
		category = COALESCE($4::TEXT, category),
		requires_urgency = COALESCE($5, requires_urgency),
		urgency_reason = COALESCE($6::TEXT, urgency_reason)
		assigned_team = COALESCE($7::TEXT, assigned_team)
	WHERE id = $8
	RETURNING *
	 `, update.SuggestedResponse, update.AverageSentiment, update.Priority, update.Category, update.RequiresUrgency, update.UrgencyReason, update.AssignedTeam, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("Failed to update ticket data")
	}
	return nil
}

func (r *Repository) Find(ctx context.Context, filters FindFilters) ([]TicketResponse, error) {
	var dbTickets []DbFindTicketByIdRow
	var ticketResponse []TicketResponse
	if err := r.db.SelectContext(ctx, &dbTickets, queries.FIND_TICKETS, filters.ToArgs()...); err != nil {
		return ticketResponse, err
	}
	ticketResponse = make([]TicketResponse, 0, len(dbTickets))
	for _, ticket := range dbTickets {
		ticketResponse = append(ticketResponse, TicketResponse{
			Id:                ticket.TicketId,
			CustomerEmail:     ticket.CustomerEmail,
			Subject:           ticket.Subject,
			Status:            ticket.Status,
			Pk:                ticket.TicketPk,
			Priority:          nullableStringPtr(ticket.Priority),
			Category:          nullableStringPtr(ticket.Category),
			AssignedTeam:      nullableStringPtr(ticket.AssignedTeam),
			SuggestedResponse: nullableStringPtr(ticket.SuggestedResponse),
			CreatedAt:         nullableTimePtr(ticket.CreatedAt),
		})
	}
	return ticketResponse, nil
}

func (r *Repository) InsertMessage(ctx context.Context, input DbInsertMessageInput) (DbCreateMessageResult, error) {
	var ins DbCreateMessageResult
	log.Println(input)
	err := r.db.GetContext(ctx, &ins, queries.INSERT_MESSAGE_QUERY, input.MessageContent, input.Role, input.TicketId, input.CustomerEmail)
	if err != nil {
		return ins, err
	}
	return ins, err
}

func (r *Repository) Create(ctx context.Context, input DbCreateTicketInput) (DbCreateTicketResult, error) {
	var createdTicket DbCreateTicketResult
	if err := r.db.Get(&createdTicket, queries.CREATE_TICKET_QUERY, input.CustomerEmail, input.Subject, input.MessageContent); err != nil {
		return createdTicket, err
	}
	return createdTicket, nil
}

func (r *Repository) FindById(ctx context.Context, input FindByIdInput) (TicketResponse, bool, error) {
	var rows []DbFindTicketByIdRow
	ticketId, limit, offset := input.TicketId, input.Pagination.GetLimit(), input.Pagination.GetOffset()
	if err := r.db.SelectContext(ctx, &rows, queries.FIND_TICKET_BY_ID_QUERY, ticketId, offset, limit); err != nil {
		return TicketResponse{}, false, err
	}
	if len(rows) == 0 {
		return TicketResponse{}, false, nil
	}

	first := rows[0]
	response := TicketResponse{
		Id:                first.TicketId,
		CustomerEmail:     first.CustomerEmail,
		Subject:           first.Subject,
		Status:            first.Status,
		Pk:                first.TicketPk,
		Priority:          nullableStringPtr(first.Priority),
		Category:          nullableStringPtr(first.Category),
		AssignedTeam:      nullableStringPtr(first.AssignedTeam),
		SuggestedResponse: nullableStringPtr(first.SuggestedResponse),
		Messages:          make([]MessageResponse, 0, len(rows)),
		Events:            first.Events,
		CreatedAt:         nil,
	}

	for _, row := range rows {
		if !row.MessageId.Valid {
			continue
		}
		response.Messages = append(response.Messages, MessageResponse{
			Id:       row.MessageId.String,
			Content:  row.MessageContent.String,
			Role:     row.MessageRole.String,
			TicketId: &first.TicketId,
		})
	}

	return response, true, nil
}

func (e DbEventResponse) ToResponse() EventResponse {
	d, _ := json.Marshal(e)
	var r DbEventResponse
	if err := json.Unmarshal(d, &r); err != nil {
		panic(err)
	}
	response := EventResponse{
		DbEventResponse: r,
		Payload:         nil,
	}
	payloadPtr := nullableStringPtr(e.Payload)
	if payloadPtr != nil {
		response.Payload = json.RawMessage(*payloadPtr)
	}
	return response
}

func (r *Repository) CreateEvent(input CreateEventInput) (EventResponse, error) {
	var createdEvent DbEventResponse
	if err := r.db.Get(&createdEvent, queries.INSERT_TICKET_EVENT, input.Payload, input.Status, input.TicketId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return EventResponse{}, ErrDuplicateTicketEvent
		}
		return createdEvent.ToResponse(), err
	}
	return createdEvent.ToResponse(), nil

}
