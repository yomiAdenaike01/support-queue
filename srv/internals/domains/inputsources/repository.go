package inputsourcesdomain

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/yomiAdenaike01/support-queue/internals/db/queries"
)

type repository struct {
	db *sqlx.DB
}

type Repository interface {
	FindMany(ctx context.Context, input FindInput) ([]InputSource, error)
	FindOne(ctx context.Context, input FindInput) (InputSource, error)
	Create(ctx context.Context, input SaveInputSourceInput) (InputSource, error)
	Update(ctx context.Context, input SaveInputSourceInput) (InputSource, error)
}

func (r *repository) FindMany(ctx context.Context, input FindInput) ([]InputSource, error) {
	log.Printf("[inputsources::repository]: Fetching many input soruce by filter=%v", input)

	var inputSources DbInputSources
	if err := r.db.SelectContext(ctx, &inputSources, queries.FIND_INPUT_SOURCES, input.toManyInputArgs()...); err != nil {
		return nil, err
	}
	return inputSources.toInputSources()
}

func (r *repository) FindOne(ctx context.Context, input FindInput) (InputSource, error) {
	log.Printf("[inputsources::repository]: Fetching one input soruce by filter=%v", input)
	var result DbInputSource

	if err := r.db.GetContext(ctx, &result, queries.FIND_INPUT_SOURCES, input.toOneInputArgs()...); err != nil {
		return result.toInputSource(), err
	}
	return result.toInputSource(), nil
}

func (r *repository) Create(ctx context.Context, input SaveInputSourceInput) (InputSource, error) {
	log.Printf("[inputsources::repository]: Creating input source type=%s name=%s", input.SourceType, input.Name)
	var result DbInputSource
	if err := r.db.GetContext(ctx, &result, queries.CREATE_INPUT_SOURCE, input.createArgs()...); err != nil {
		return result.toInputSource(), err
	}
	return result.toInputSource(), nil
}

func (r *repository) Update(ctx context.Context, input SaveInputSourceInput) (InputSource, error) {
	log.Printf("[inputsources::repository]: Updating input source id=%v type=%s name=%s", input.Id, input.SourceType, input.Name)
	var result DbInputSource
	if err := r.db.GetContext(ctx, &result, queries.UPDATE_INPUT_SOURCE, input.updateArgs()...); err != nil {
		return result.toInputSource(), err
	}
	return result.toInputSource(), nil
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}
