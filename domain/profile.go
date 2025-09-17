package domain

import (
	"context"
	"prakarsa-app/transport/request"
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
	Linkedin      string    `json:"linkedin"`
}

type UserProfile struct {
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	PhoneNumber   string `json:"phone_number"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	NameAlias     string `json:"name_alias"`
	Avatar        string `json:"avatar"`
	Gender        string `json:"gender"`
	BirthDate     string `json:"birth_date"`
	SlugName      string `json:"slug_name"`
	AboutMe       string `json:"about_me"`
	InstitutionID string `json:"institution_id"`
	JoinedAt      int64  `json:"joined_at"`
	Linkedin      string `json:"linkedin"`
}

type CompleteProfile struct {
	BirthDate     time.Time `json:"birth_date"`
	Gender        string    `json:"gender"`
	InstitutionID string    `json:"institution_id"`
	UpdatedAt     int64     `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
}

type UpdateProfile struct {
	Name          string    `json:"name"`
	NameAlias     string    `json:"name_alias"`
	Avatar        string    `json:"avatar"`
	Gender        string    `json:"gender"`
	BirthDate     time.Time `json:"birth_date"`
	SlugName      string    `json:"slug_name"`
	AboutMe       string    `json:"about_me"`
	InstitutionID string    `json:"institution_id"`
	UpdatedAt     int64     `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	Linkedin      string    `json:"linkedin"`
}

type ProfileRepository interface {
	Create(ctx context.Context, profile *Profile) error
	GetByUserID(ctx context.Context, userID string) (Profile, error)
	CompleteProfile(ctx context.Context, userID string, completeProfile *CompleteProfile) error
	GetUserProfileByUserID(ctx context.Context, userID string) (UserProfile, error)
	UpdateProfile(ctx context.Context, userID string, updateProfile *UpdateProfile) error
}

type ProfileUsecase interface {
	CompleteProfile(ctx context.Context, userID string, request *request.CompleteProfileReq) error
	ProfileCompletion(ctx context.Context, userID string) (bool, error)
	UserProfile(ctx context.Context, userID string) (UserProfile, error)
	UpdateProfile(ctx context.Context, userID string, request *request.UpdateProfileReq) error
}
