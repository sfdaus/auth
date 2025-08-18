package usecase

import (
	"context"
	"database/sql"
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
	if err != nil && err != sql.ErrNoRows {
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
