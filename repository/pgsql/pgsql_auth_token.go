package pgsql

import (
	"context"
	"database/sql"
	"prakarsa-app/domain"
)

type pgsqlAuthTokenRepository struct {
	db *sql.DB
}

func NewPgsqlAuthTokenRepository(db *sql.DB) *pgsqlAuthTokenRepository {
	return &pgsqlAuthTokenRepository{
		db: db,
	}
}

func (r *pgsqlAuthTokenRepository) Create(ctx context.Context, authToken *domain.AuthToken) (err error) {
	query := "INSERT INTO auth_tokens (id, user_id, user_agent, refresh_token, expires_at, issued_at, is_active, created_at, created_by, updated_at, updated_by, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
	_, err = r.db.ExecContext(ctx, query, authToken.ID, authToken.UserID, authToken.UserAgent, authToken.RefreshToken, authToken.ExpiresAt, authToken.IssuedAt, authToken.IsActive, authToken.CreatedAt, authToken.CreatedBy, authToken.UpdatedAt, authToken.UpdatedBy, authToken.DeletedAt)
	return
}

func (r *pgsqlAuthTokenRepository) GetByUserID(ctx context.Context, userID string) (authToken domain.AuthToken, err error) {
	query := "SELECT id, user_id, user_agent, refresh_token, expires_at, issued_at, is_active, created_at, created_by, updated_at, updated_by, deleted_at FROM auth_tokens WHERE user_id = $1"
	err = r.db.QueryRowContext(ctx, query, userID).Scan(
		&authToken.ID, &authToken.UserID, &authToken.UserAgent, &authToken.RefreshToken,
		&authToken.ExpiresAt, &authToken.IssuedAt, &authToken.IsActive,
		&authToken.CreatedAt, &authToken.CreatedBy, &authToken.UpdatedAt,
		&authToken.UpdatedBy, &authToken.DeletedAt,
	)
	return
}

func (r *pgsqlAuthTokenRepository) UpdateRefreshToken(ctx context.Context, userID string, updateAuthToken domain.UpdateAuthToken) (err error) {
	query := "UPDATE auth_tokens SET user_agent = $1, refresh_token = $2, expires_at = $3, issued_at = $4, updated_at = $5, updated_by = $6 WHERE user_id = $7"
	_, err = r.db.ExecContext(ctx, query, updateAuthToken.UserAgent, updateAuthToken.RefreshToken,
		updateAuthToken.ExpiresAt, updateAuthToken.IssuedAt, updateAuthToken.UpdatedAt,
		updateAuthToken.UpdatedBy, userID)
	return
}
