package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"prakarsa-app/domain"
	"strings"
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
		       p.birth_date, p.slug_name, p.about_me, p.institution_id, p.created_at AS joined_at, COALESCE(p.linkedin, '') as linkedin
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
	sets := make([]string, 0, 12)
	args := make([]any, 0, 12)
	i := 1

	setVal := func(col string, v any) {
		sets = append(sets, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, v)
		i++
	}
	setNull := func(col string) {
		sets = append(sets, fmt.Sprintf("%s = ''", col))
	}

	// ---------- Name + Alias + Slug ----------
	// Hanya update kalau Name dikirim (non-empty)
	if updateProfile.Name != "" {
		setVal("name", updateProfile.Name)
		// Kamu sudah generate NameAlias & SlugName di usecase; tetap guard di sini
		if updateProfile.NameAlias != "" {
			setVal("name_alias", updateProfile.NameAlias)
		}
		if updateProfile.SlugName != "" {
			setVal("slug_name", updateProfile.SlugName)
		}
	}

	// ---------- Avatar ----------
	// *updateProfile.Avatar: nil = no-change, "" = delete(NULL), non-empty = set
	if updateProfile.Avatar != nil {
		if *updateProfile.Avatar == "" {
			setNull("avatar")
		} else {
			setVal("avatar", *updateProfile.Avatar)
		}
	}

	// ---------- Gender ----------
	if updateProfile.Gender != "" {
		setVal("gender", updateProfile.Gender)
	}

	// ---------- Birth Date ----------
	// Update kalau non-zero
	if !updateProfile.BirthDate.IsZero() {
		setVal("birth_date", updateProfile.BirthDate.UTC())
	}

	// ---------- About Me ----------
	// *updateProfile.AboutMe: nil = no-change, "" = delete(NULL), non-empty = set
	if updateProfile.AboutMe != nil {
		if *updateProfile.AboutMe == "" {
			setNull("about_me")
		} else {
			setVal("about_me", *updateProfile.AboutMe)
		}
	}

	// ---------- Institution ----------
	if updateProfile.InstitutionID != "" {
		setVal("institution_id", updateProfile.InstitutionID)
	}

	// ---------- Linkedin ----------
	// *updateProfile.Linkedin: nil = no-change, "" = delete(NULL), non-empty = set
	if updateProfile.Linkedin != nil {
		if *updateProfile.Linkedin == "" {
			setNull("linkedin")
		} else {
			setVal("linkedin", *updateProfile.Linkedin)
		}
	}

	// ---------- Audit ----------
	setVal("updated_at", updateProfile.UpdatedAt)
	setVal("updated_by", updateProfile.UpdatedBy)

	if len(sets) == 0 {
		// tidak ada perubahan
		return nil
	}

	q := fmt.Sprintf(`UPDATE profiles SET %s WHERE user_id = $%d`, strings.Join(sets, ", "), i)
	args = append(args, userID)

	_, err = r.db.ExecContext(ctx, q, args...)
	return err
}

func (r *pgsqlProfileRepository) GetUserProfileByPublicID(ctx context.Context, publicID string, userID string) (userProfile domain.SecureUserProfile, err error) {
	query := `
		SELECT p.name, p.name_alias, p.avatar, p.slug_name, p.about_me, p.institution_id, p.created_at AS joined_at, p.linkedin,
		       (u.id = $2) AS is_my_profile
		FROM profiles p
		JOIN users u ON p.user_id = u.id
		WHERE u.public_id = $1
	`
	err = r.db.QueryRowContext(ctx, query, publicID, userID).Scan(
		&userProfile.Name, &userProfile.NameAlias, &userProfile.Avatar, &userProfile.SlugName, &userProfile.AboutMe, &userProfile.InstitutionID,
		&userProfile.JoinedAt, &userProfile.Linkedin, &userProfile.IsMyProfile,
	)
	return
}
