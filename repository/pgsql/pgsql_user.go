package pgsql

import (
	"context"
	"database/sql"

	"prakarsa-app/domain"
)

type pgsqlUserRepository struct {
	db *sql.DB
}

func NewPgsqlUserRepository(db *sql.DB) *pgsqlUserRepository {
	return &pgsqlUserRepository{
		db: db,
	}
}

func (r *pgsqlUserRepository) Create(ctx context.Context, user *domain.User) (err error) {
	query := "INSERT INTO users (id, password, username, phone_number, email, token_version, is_active, created_at, created_by, updated_at, updated_by, deleted_at, is_verified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)"
	_, err = r.db.ExecContext(ctx, query, user.ID, user.Password, user.Username, user.PhoneNumber, user.Email, user.TokenVersion, user.IsActive, user.CreatedAt, user.CreatedBy, user.UpdatedAt, user.UpdatedBy, user.DeletedAt, user.IsVerified)
	return
}

func (r *pgsqlUserRepository) UpdateTokenVersionByID(ctx context.Context, tokenVersion string, id string) (err error) {
	query := "UPDATE users SET token_version = $1 WHERE id = $2"
	_, err = r.db.ExecContext(ctx, query, tokenVersion, id)
	return
}

func (r *pgsqlUserRepository) GetByEmail(ctx context.Context, email string) (user domain.User, err error) {
	query := "SELECT id, email, password, is_active, deleted_at FROM users WHERE email = $1"
	err = r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.IsActive, &user.DeletedAt)
	return
}

func (r *pgsqlUserRepository) GetByUsername(ctx context.Context, username string) (user domain.User, err error) {
	query := "SELECT id, username, password, is_active, deleted_at FROM users WHERE email = $1"
	err = r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.IsActive, &user.DeletedAt)
	return
}

func (r *pgsqlUserRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (user domain.User, err error) {
	query := "SELECT id, phone_number, password, is_active, deleted_at FROM users WHERE email = $1"
	err = r.db.QueryRowContext(ctx, query, phoneNumber).Scan(&user.ID, &user.PhoneNumber, &user.Password, &user.IsActive, &user.DeletedAt)
	return
}

func (r *pgsqlUserRepository) GetByUserID(ctx context.Context, userID string) (user domain.User, err error) {
	query := "SELECT id, email, password, is_active, deleted_at FROM users WHERE id = $1"
	err = r.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Email, &user.Password, &user.IsActive, &user.DeletedAt)
	return
}

func (r *pgsqlUserRepository) UpdatePasswordByUserID(ctx context.Context, userID string, newPassword string) (err error) {
	query := "UPDATE users SET password = $1 WHERE id = $2 and is_active = true"
	_, err = r.db.ExecContext(ctx, query, newPassword, userID)
	return
}

func (r *pgsqlUserRepository) VerifyAccountByUserID(ctx context.Context, userID string) (err error) {
	query := "UPDATE users SET is_verified = true WHERE id = $1 and is_active = true"
	_, err = r.db.ExecContext(ctx, query, userID)
	return
}
