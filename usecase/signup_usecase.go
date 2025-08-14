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
	cryptoSvc      crypto.CryptoService
	jwtSvc         jwt.JWTService
	contextTimeout time.Duration
}

func SignUpUsecase(userRepo domain.UserRepository, cryptoSvc crypto.CryptoService, jwtSvc jwt.JWTService, contextTimeout time.Duration) *signupUsecase {
	return &signupUsecase{
		userRepo:       userRepo,
		cryptoSvc:      cryptoSvc,
		jwtSvc:         jwtSvc,
		contextTimeout: contextTimeout,
	}
}

func (u *signupUsecase) SignUp(c context.Context, request *request.SignUpReq) (accessToken string, err error) {
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

	err = u.userRepo.Create(ctx, &domain.User{
		ID:           newUserID,
		Email:        request.Email,
		PhoneNumber:  request.PhoneNumber,
		Username:     request.Username,
		TokenVersion: tokenVersion,
		IsActive:     true,
		Password:     passwordHash,
		CreatedAt:    time.Now().Unix(),
		CreatedBy:    "",
		UpdatedAt:    0,
		UpdatedBy:    "",
		DeletedAt:    0,
	})

	accessToken, err = u.jwtSvc.GenerateToken(ctx, newUserID, tokenVersion)
	return
}
