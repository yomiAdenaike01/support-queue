package knowledgebase

import (
	"context"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/pgvector/pgvector-go"
	"github.com/yomiAdenaike01/support-queue/internals/db/queries"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, input CreateRequest) (DbCreateKnowledgeBaseRow, error) {
	var created DbCreateKnowledgeBaseRow
	metadata := input.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage("null")
	}

	err := r.db.GetContext(
		ctx,
		&created,
		queries.INSERT_KNOWLEDGE_BASE,
		input.SourceId,
		input.SourceEventType,
		metadata,
		input.Content,
		pgvector.NewVector(input.Embedding),
	)
	return created, err
}
