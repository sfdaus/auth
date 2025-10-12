package usecase

import (
	"context"
	"database/sql"
	"prakarsa-app/config/constant"
	"prakarsa-app/domain"
	"prakarsa-app/infrastructure/filestorage"
	"prakarsa-app/transport/request"
	"prakarsa-app/utils"
	"time"

	"github.com/google/uuid"
)

type profileUsecase struct {
	profileRepo    domain.ProfileRepository
	contextTimeout time.Duration
	fileStorage    filestorage.FileStorage
}

func ProfileUsecase(profileRepo domain.ProfileRepository, contextTimeout time.Duration, fileStorage filestorage.FileStorage) *profileUsecase {
	return &profileUsecase{
		profileRepo:    profileRepo,
		contextTimeout: contextTimeout,
		fileStorage:    fileStorage,
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

	if profile.BirthDate.IsZero() || profile.Gender == "" || profile.InstitutionID == nil {
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

	userProfile.Email, err = utils.DecryptDeterministic(userProfile.Email)
	if err != nil {
		return domain.UserProfile{}, err
	}

	if userProfile.Avatar != nil {
		*userProfile.Avatar, err = u.fileStorage.GetURL(ctx, *userProfile.Avatar, time.Hour*24)
		if err != nil {
			return domain.UserProfile{}, err
		}
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
	var empty = ""

	updateProfilePayload := &domain.UpdateProfile{
		Gender:    request.Gender,
		UpdatedAt: time.Now().Unix(),
		UpdatedBy: profile.UserID,
	}

	if request.BirthDate != "" {
		birthDate, err = utils.ParseDateYYYYMMDD(request.BirthDate)
		if err != nil {
			return
		}

		updateProfilePayload.BirthDate = birthDate
	}

	var nameAlias, slugName string
	if request.Name != "" {
		nameAlias = utils.GenerateNameAlias(request.Name)
		slugName = utils.GenerateSlugName(request.Name)
		updateProfilePayload.Name = request.Name
	} else {
		nameAlias = profile.NameAlias
		slugName = profile.SlugName
	}
	updateProfilePayload.SlugName = slugName
	updateProfilePayload.NameAlias = nameAlias
	avatarPath := profile.Avatar
	if request.AvatarDelete != nil && utils.AsBool(request.AvatarDelete) {
		avatarPath = &empty
	} else if request.Avatar != nil {
		fileAvatar, errFile := request.Avatar.Open()
		if errFile != nil {
			return errFile
		}
		defer fileAvatar.Close()
		tempAvatar, err := u.fileStorage.Put(ctx, "avatars/"+uuid.NewString(), fileAvatar)
		avatarPath = &tempAvatar
		if err != nil {
			return err
		}
	}
	updateProfilePayload.Avatar = avatarPath

	if request.AboutMeDelete != nil && utils.AsBool(request.AboutMeDelete) {
		updateProfilePayload.AboutMe = &empty
	} else if request.AboutMe != "" {
		updateProfilePayload.AboutMe = &request.AboutMe
	}

	if request.LinkedinDelete != nil && utils.AsBool(request.LinkedinDelete) {
		updateProfilePayload.Linkedin = &empty
	} else if request.Linkedin != "" {
		updateProfilePayload.Linkedin = &request.Linkedin
	}

	if request.InstitutionDelete != nil && utils.AsBool(request.InstitutionDelete) {
		updateProfilePayload.InstitutionID = &empty
	} else if request.InstitutionID != "" {
		updateProfilePayload.InstitutionID = &request.InstitutionID
	}

	err = u.profileRepo.UpdateProfile(ctx, userID, updateProfilePayload)

	if err != nil {
		return
	}

	if request.Avatar != nil || utils.AsBool(request.AvatarDelete) {
		if profile.Avatar != nil {
			_ = u.fileStorage.Delete(ctx, *profile.Avatar)
		}
	}

	return
}

func (u *profileUsecase) UserProfileByID(ctx context.Context, request *request.UserProfileByIDReq) (userProfile domain.SecureUserProfile, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	userProfile, err = u.profileRepo.GetUserProfileByPublicID(ctx, request.PublicID, request.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewBadRequestError(constant.UserProfileMessage.UserProfileUserNotFound)
			return domain.SecureUserProfile{}, err
		}
		return domain.SecureUserProfile{}, err
	}

	tempAvatar, err := u.fileStorage.GetURL(ctx, *userProfile.Avatar, time.Hour*24)
	if err != nil {
		return domain.SecureUserProfile{}, err
	}

	userProfile.Avatar = &tempAvatar

	return userProfile, nil
}
