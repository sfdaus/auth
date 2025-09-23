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
	query := "INSERT INTO profiles (user_id, name, name_alias, avatar, gender, birth_date, slug_name, about_me, institution_id, is_active, created_at, created_by, updated_at, updated_by, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)"
	_, err = r.db.ExecContext(ctx, query, profile.UserID, profile.Name, profile.NameAlias, profile.Avatar, profile.Gender, profile.BirthDate, profile.SlugName, profile.AboutMe, profile.InstitutionID, profile.IsActive, profile.CreatedAt, profile.CreatedBy, profile.UpdatedAt, profile.UpdatedBy, profile.DeletedAt)
	return
}

func (r *pgsqlProfileRepository) GetByUserID(ctx context.Context, userID string) (profile domain.Profile, err error) {
	query := "SELECT user_id, name, name_alias, avatar, gender, birth_date, slug_name, about_me, institution_id, is_active, created_at, created_by, updated_at, updated_by, deleted_at, linkedin FROM profiles WHERE user_id = $1"
	err = r.db.QueryRowContext(ctx, query, userID).Scan(&profile.UserID, &profile.Name, &profile.NameAlias, &profile.Avatar, &profile.Gender, &profile.BirthDate, &profile.SlugName, &profile.AboutMe, &profile.InstitutionID, &profile.IsActive, &profile.CreatedAt, &profile.CreatedBy, &profile.UpdatedAt, &profile.UpdatedBy, &profile.DeletedAt, &profile.Linkedin)
	return
}

func (r *pgsqlProfileRepository) CompleteProfile(ctx context.Context, userID string, completeProfile *domain.CompleteProfile) (err error) {
	query := "UPDATE profiles SET birth_date = $1, gender = $2, institution_id = $3, updated_at = $4, updated_by = $5 WHERE user_id = $6"
	_, err = r.db.ExecContext(ctx, query, completeProfile.BirthDate, completeProfile.Gender, completeProfile.InstitutionID, completeProfile.UpdatedAt, completeProfile.UpdatedBy, userID)
	return
}

func (r *pgsqlProfileRepository) GetUserProfileByUserID(ctx context.Context, userID string) (userProfile domain.UserProfile, err error) {
	query := `
		SELECT p.user_id, u.username, u.phone_number, u.email, p.name, p.name_alias, p.avatar, p.gender,
		       p.birth_date, p.slug_name, p.about_me, p.institution_id, p.created_at AS joined_at, p.linkedin
		FROM profiles p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = $1
	`
	err = r.db.QueryRowContext(ctx, query, userID).Scan(
		&userProfile.UserID, &userProfile.Username, &userProfile.PhoneNumber, &userProfile.Email,
		&userProfile.Name, &userProfile.NameAlias, &userProfile.Avatar, &userProfile.Gender,
		&userProfile.BirthDate, &userProfile.SlugName, &userProfile.AboutMe, &userProfile.InstitutionID,
		&userProfile.JoinedAt, &userProfile.Linkedin,
	)
	return
}

func (r *pgsqlProfileRepository) UpdateProfile(ctx context.Context, userID string, updateProfile *domain.UpdateProfile) (err error) {
	query := "UPDATE profiles SET name = COALESCE(NULLIF($1, ''), name), name_alias = COALESCE(NULLIF($2, ''), name_alias), avatar = COALESCE(NULLIF($3, ''), avatar), gender = COALESCE(NULLIF($4, ''), gender), birth_date = COALESCE(NULLIF(TO_DATE($5, 'YYYY-MM-DD'), DATE '0001-01-01'), birth_date), slug_name = COALESCE(NULLIF($6, ''), slug_name), about_me = COALESCE(NULLIF($7, ''), about_me), institution_id = COALESCE(NULLIF($8, ''), institution_id), linkedin = COALESCE(NULLIF($9, ''), linkedin), updated_at = $10, updated_by = $11 WHERE user_id = $12"
	_, err = r.db.ExecContext(ctx, query, updateProfile.Name, updateProfile.NameAlias, updateProfile.Avatar, updateProfile.Gender, updateProfile.BirthDate, updateProfile.SlugName, updateProfile.AboutMe, updateProfile.InstitutionID, updateProfile.Linkedin, updateProfile.UpdatedAt, updateProfile.UpdatedBy, userID)
	return
}
