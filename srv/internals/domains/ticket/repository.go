package ticket

import (
	"context"
	"database/sql"
	"log"

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
	query, args, err := sqlx.Named(queries.CREATE_TICKET_QUERY, input)
	if err != nil {
		return createdTicket, err
	}
	query = r.db.Rebind(query)

	if err := r.db.Get(&createdTicket, query, args...); err != nil {
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

func nullableStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
