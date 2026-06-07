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
	content := input.Content

	if len(metadata) == 0 {
		metadata = json.RawMessage("null")
	}
	if len(input.Content) == 0 {
		content = json.RawMessage("null")

	}

	err := r.db.GetContext(
		ctx,
		&created,
		queries.INSERT_KNOWLEDGE_BASE,
		input.SourceId,
		input.SourceEventType,
		metadata,
		content,
		pgvector.NewVector(input.Embedding),
	)
	return created, err
}

func (r *Repository) Find(ctx context.Context, search searchInput) ([]DbCreateKnowledgeBaseRow, error) {
	var dest []DbCreateKnowledgeBaseRow
	var embedding pgvector.Vector
	if search.Embedding != nil {
		embedding = pgvector.NewVector(search.Embedding)
	}
	if err := r.db.SelectContext(ctx, &dest, queries.FIND_SIMILAR_TICKETS, search.Id, embedding, search.SourceType, search.getOffset(), search.getLimit()); err != nil {
		return dest, err
	}
	return dest, nil
}
