package pgsql

import (
	"context"
	"database/sql"
	"prakarsa-app/domain"
)

type pgsqlProfileRepository struct {
	db *sql.DB
}

func NewPgsqlProfileRepository(db *sql.DB) *pgsqlProfileRepository {
	return &pgsqlProfileRepository{
		db: db,
	}
}

func (r *pgsqlProfileRepository) Create(ctx context.Context, profile *domain.Profile) (err error) {
	query := "INSERT INTO profile (user_id, name, name_alias, avatar, gender, birth_date, slug_name, about_me, institution_id, is_active, created_at, created_by, updated_at, updated_by, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)"
	_, err = r.db.ExecContext(ctx, query, profile.UserID, profile.Name, profile.NameAlias, profile.Avatar, profile.Gender, profile.BirthDate, profile.SlugName, profile.AboutMe, profile.InstitutionID, profile.IsActive, profile.CreatedAt, profile.CreatedBy, profile.UpdatedAt, profile.UpdatedBy, profile.DeletedAt)
	return
}

func (r *pgsqlProfileRepository) GetByUserID(ctx context.Context, userID string) (profile domain.Profile, err error) {
	query := "SELECT user_id, name, name_alias, avatar, gender, birth_date, slug_name, about_me, institution_id, is_active, created_at, created_by, updated_at, updated_by, deleted_at FROM profile WHERE user_id = $1"
	err = r.db.QueryRowContext(ctx, query, userID).Scan(&profile.UserID, &profile.Name, &profile.NameAlias, &profile.Avatar, &profile.Gender, &profile.BirthDate, &profile.SlugName, &profile.AboutMe, &profile.InstitutionID, &profile.IsActive, &profile.CreatedAt, &profile.CreatedBy, &profile.UpdatedAt, &profile.UpdatedBy, &profile.DeletedAt)
	return
}

func (r *pgsqlProfileRepository) CompleteProfile(ctx context.Context, userID string, completeProfile *domain.CompleteProfile) (err error) {
	query := "UPDATE profile SET birth_date = $1, gender = $2, institution_id = $3, updated_at = $4, updated_by = $5 WHERE user_id = $6"
	_, err = r.db.ExecContext(ctx, query, completeProfile.BirthDate, completeProfile.Gender, completeProfile.InstitutionID, completeProfile.UpdatedAt, completeProfile.UpdatedBy, userID)
	return
}
