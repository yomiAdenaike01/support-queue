package ticket

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/yomiAdenaike01/support-queue/internals/db/queries"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindById(ctx context.Context, id string) (TicketResponse, bool, error) {
	var rows []DbFindTicketByIdRow
	if err := r.db.SelectContext(ctx, &rows, queries.FIND_TICKET_BY_ID_QUERY, id); err != nil {
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
			Id:      row.MessageId.String,
			Content: row.MessageContent.String,
			Role:    row.MessageRole.String,
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
