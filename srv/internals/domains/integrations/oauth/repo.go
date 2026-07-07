package oauth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type tokens struct {
	AccessToken      string    `db:"access_token"`
	RefreshToken     string    `db:"refresh_token"`
	AccessExpiresAt  time.Time `db:"access_expires_at"`
	RefreshExpiresAt time.Time `db:"refresh_expires_at"`
}

func (t tokens) ShouldRefreshAccessToken() bool {
	diff := t.AccessExpiresAt.Sub(time.Now())
	return diff.Seconds() <= float64(time.Second*30)
}

type saveTokenInput struct {
	provider         ProviderName
	accessToken      string
	refreshToken     *string
	accessExpiresAt  time.Time
	refreshExpiresAt *time.Time
}

type AuthRepository interface {
	saveTokens(ctx context.Context, input saveTokenInput) error
	GetTokens(ctx context.Context, provider ProviderName) (tokens, error)
}

type repository struct {
	db *sqlx.DB
}

func (r *repository) saveTokens(ctx context.Context, input saveTokenInput) error {
	_, err := r.db.ExecContext(ctx, `
	with existing AS (
		UPDATE auth_credentials
		SET 
			access_token = COALESCE($1::TEXT, auth_credentials.access_token),
			refresh_token = COALESCE($2::TEXT, auth_credentials.refresh_token),
			refresh_expires_at = COALESCE($3::TIMESTAMPTZ, auth_credentials.refresh_expires_at),
			access_expires_at = COALESCE($5::TIMESTAMPTZ, auth_credentials.access_expires_at)
		WHERE provider = $4
		RETURNING 1
	)
	INSERT INTO auth_credentials (
		access_token, 
		refresh_token, 
		access_expires_at, 
		refresh_expires_at
	)
	SELECT 
		$1::TEXT, 
		$2::TEXT,
		$3::TIMESTAMPTZ,
		$5::TIMESTAMPTZ
	WHERE NOT EXISTS (SELECT 1 FROM existing)
	`, input.accessToken, input.refreshToken, input.refreshExpiresAt, input.provider, input.accessExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetTokens(ctx context.Context, provider ProviderName) (tokens, error) {
	var tokens tokens
	if err := r.db.GetContext(ctx, &tokens, `
		SELECT 
			access_token, 
			refresh_token,
			refresh_expires_at,
			access_expires_at
		FROM auth_credentials 
		WHERE provider = $1
	`, provider); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tokens, nil
		}
		return tokens, err
	}
	return tokens, nil
}

func NewCredentialsRepository(db *sqlx.DB) AuthRepository {
	return &repository{
		db: db,
	}
}
