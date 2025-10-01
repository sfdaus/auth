package usecase

import (
	"context"
	"database/sql"
	"time"

	"prakarsa-app/config/constant"
	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/utils"
	"prakarsa-app/utils/crypto"
	"prakarsa-app/utils/jwt"
)

type signupUsecase struct {
	userRepo       domain.UserRepository
	authTokenRepo  domain.AuthTokenRepository
	profileRepo    domain.ProfileRepository
	cryptoSvc      crypto.CryptoService
	jwtSvc         jwt.JWTService
	contextTimeout time.Duration
}

func SignUpUsecase(userRepo domain.UserRepository, authTokenRepo domain.AuthTokenRepository, profileRepo domain.ProfileRepository, cryptoSvc crypto.CryptoService, jwtSvc jwt.JWTService, contextTimeout time.Duration) *signupUsecase {
	return &signupUsecase{
		userRepo:       userRepo,
		authTokenRepo:  authTokenRepo,
		profileRepo:    profileRepo,
		cryptoSvc:      cryptoSvc,
		jwtSvc:         jwtSvc,
		contextTimeout: contextTimeout,
	}
}

func (u *signupUsecase) SignUp(c context.Context, request *request.SignUpReq) (accessToken string, userID string, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, request.Email)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if user.ID != "" {
		err = utils.NewBadRequestError(constant.SignupMessage.SignupExists)
		return
	}

	passwordHash, err := u.cryptoSvc.CreatePasswordHash(ctx, request.Password)
	if err != nil {
		return
	}

	tokenVersion := utils.GenerateTokenVersion()
	newUserID := utils.GenerateUUID()

	// Create new user
	err = u.userRepo.Create(ctx, &domain.User{
		ID:           newUserID,
		Email:        request.Email,
		PhoneNumber:  "",
		Username:     "",
		TokenVersion: tokenVersion,
		IsActive:     true,
		Password:     passwordHash,
		CreatedAt:    time.Now().Unix(),
		CreatedBy:    "",
		UpdatedAt:    0,
		UpdatedBy:    "",
		DeletedAt:    0,
	})
	if err != nil {
		errorInfo := utils.ErrorInfo{
			Message: constant.SignupMessage.SignupFailedToCreateUser,
			Details: err.Error(),
		}
		err = utils.NewBadRequestError(errorInfo)
		return
	}

	nameAlias := utils.GenerateNameAlias(request.Name)
	// Create profile for the new user
	var empty = ""
	err = u.profileRepo.Create(ctx, &domain.Profile{
		UserID:        newUserID,
		Name:          request.Name,
		NameAlias:     nameAlias,
		Avatar:        &empty,
		Gender:        "",
		BirthDate:     time.Time{},
		SlugName:      utils.GenerateSlugName(request.Name),
		AboutMe:       &empty,
		InstitutionID: &empty,
		IsActive:      true,
		CreatedAt:     time.Now().Unix(),
		CreatedBy:     "",
		UpdatedAt:     0,
		UpdatedBy:     "",
		DeletedAt:     0,
	})
	if err != nil {
		errorInfo := utils.ErrorInfo{
			Message: constant.SignupMessage.SignupFailedToCreateProfile,
			Details: err.Error(),
		}
		err = utils.NewBadRequestError(errorInfo)
		return
	}

	// Create auth token for the new user
	err = u.authTokenRepo.Create(ctx, &domain.AuthToken{
		ID:           utils.GenerateUUID(),
		UserID:       newUserID,
		UserAgent:    "",
		RefreshToken: utils.GenerateUUID(),
		ExpiresAt:    time.Now().Add(time.Hour * 24 * 30).Unix(),
		IssuedAt:     time.Now().Unix(),
		IsActive:     true,
		CreatedAt:    time.Now().Unix(),
		CreatedBy:    "",
		UpdatedAt:    0,
		UpdatedBy:    "",
		DeletedAt:    0,
	})
	if err != nil {
		errorInfo := utils.ErrorInfo{
			Message: constant.SignupMessage.SignupFailedToCreateAuthToken,
			Details: err.Error(),
		}
		err = utils.NewBadRequestError(errorInfo)
		return
	}

	userID = newUserID
	// Generate JWT token for the new user
	accessToken, err = u.jwtSvc.GenerateToken(ctx, newUserID, tokenVersion)
	if err != nil {
		err = utils.NewBadRequestError(constant.SignupMessage.SignupFailedToGenerateToken)
		return
	}
	return
}
