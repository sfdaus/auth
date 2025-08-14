package pgsql

import (
	"database/sql"
)

type pgsqlAuthTokenRepository struct {
	db *sql.DB
}

func NewPgsqlAuthTokenRepository(db *sql.DB) *pgsqlAuthTokenRepository {
	return &pgsqlAuthTokenRepository{
		db: db,
	}
}

// func (r *pgsqlAuthTokenRepository) Create(ctx context.Context, authToken *domain.AuthToken) (err error) {
// 	query := "INSERT INTO auth_tokens (id, user_id, user_agent, refresh_token, expires_at, issued_at, is_active, created_at, created_by, updated_at, updated_by, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
// 	_, err = r.db.ExecContext(ctx, query, user.ID, user.Password, user.Username, user.PhoneNumber, user.Email, user.TokenVersion, user.IsActive, user.CreatedAt, user.CreatedBy, user.UpdatedAt, user.UpdatedBy, user.DeletedAt)
// 	return
// }
