package usecase

import (
	"context"
	"database/sql"
	"prakarsa-app/config/constant"
	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/utils"
	"time"
)

type profileUsecase struct {
	profileRepo    domain.ProfileRepository
	contextTimeout time.Duration
}

func ProfileUsecase(profileRepo domain.ProfileRepository, contextTimeout time.Duration) *profileUsecase {
	return &profileUsecase{
		profileRepo:    profileRepo,
		contextTimeout: contextTimeout,
	}
}

func (u *profileUsecase) CompleteProfile(ctx context.Context, userID string, request *request.CompleteProfileReq) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	profile, err := u.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewBadRequestError(constant.ProfileCompletionMessage.ProfileCompletionUserNotFound)
			return
		}
		return
	}

	birthDate, err := utils.ParseDateYYYYMMDD(request.BirthDate)
	if err != nil {
		return
	}
	err = u.profileRepo.CompleteProfile(ctx, userID, &domain.CompleteProfile{
		BirthDate:     birthDate,
		Gender:        request.Gender,
		InstitutionID: request.InstitutionID,
		UpdatedAt:     time.Now().Unix(),
		UpdatedBy:     profile.UserID,
	})
	if err != nil {
		return
	}

	return nil
}

func (u *profileUsecase) ProfileCompletion(ctx context.Context, userID string) (isComplete bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	profile, err := u.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewBadRequestError(constant.ProfileCompletionMessage.ProfileCompletionUserNotFound)
			return
		}
		return
	}

	isComplete = false

	if profile.BirthDate.IsZero() || profile.Gender == "" || profile.InstitutionID == "" {
		isComplete = false
	} else {
		isComplete = true
	}

	return isComplete, nil
}

func (u *profileUsecase) UserProfile(ctx context.Context, userID string) (domain.UserProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	userProfile, err := u.profileRepo.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewBadRequestError(constant.UserProfileMessage.UserProfileUserNotFound)
			return domain.UserProfile{}, err
		}
		return domain.UserProfile{}, err
	}

	return userProfile, nil
}

func (u *profileUsecase) UpdateProfile(ctx context.Context, userID string, request *request.UpdateProfileReq) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	profile, err := u.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewBadRequestError(constant.UserProfileMessage.UserProfileUserNotFound)
			return
		}
		return
	}

	var birthDate time.Time
	if request.BirthDate != "" {
		birthDate, err = utils.ParseDateYYYYMMDD(request.BirthDate)
		if err != nil {
			return
		}
	}

	nameAlias := ""
	if request.Name != "" {
		nameAlias = utils.GenerateNameAlias(request.Name)
	}

	err = u.profileRepo.UpdateProfile(ctx, userID, &domain.UpdateProfile{
		Name:          request.Name,
		NameAlias:     nameAlias,
		BirthDate:     birthDate,
		Gender:        request.Gender,
		InstitutionID: request.InstitutionID,
		UpdatedAt:     time.Now().Unix(),
		UpdatedBy:     profile.UserID,
	})
	if err != nil {
		return
	}
	return nil
}
