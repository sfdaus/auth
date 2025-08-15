package domain

import (
	"context"
	"time"
)

// Profile ...
type Profile struct {
	UserID        string    `json:"user_id"`
	Name          string    `json:"name"`
	NameAlias     string    `json:"name_alias"`
	Avatar        string    `json:"avatar"`
	Gender        string    `json:"gender"`
	BirthDate     time.Time `json:"birth_date"`
	SlugName      string    `json:"slug_name"`
	AboutMe       string    `json:"about_me"`
	InstitutionID string    `json:"institution_id"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     int64     `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     int64     `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	DeletedAt     int64     `json:"deleted_at"`
}

type CompleteProfile struct {
	BirthDate     time.Time `json:"birth_date"`
	Gender        string    `json:"gender"`
	InstitutionID string    `json:"institution_id"`
	UpdatedAt     int64     `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
}

type ProfileRepository interface {
	Create(ctx context.Context, profile *Profile) error
	GetByUserID(ctx context.Context, userID string) (Profile, error)
	CompleteProfile(ctx context.Context, userID string, completeProfile *CompleteProfile) error
}
