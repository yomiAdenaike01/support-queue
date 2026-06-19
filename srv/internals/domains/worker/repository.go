package workerdomain

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func (r *Repository) GetWorkerHealth(ctx context.Context) {
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}
