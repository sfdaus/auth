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

type signinUsecase struct {
	userRepo       domain.UserRepository
	cryptoSvc      crypto.CryptoService
	jwtSvc         jwt.JWTService
	contextTimeout time.Duration
}

func SignInUsecase(userRepo domain.UserRepository, cryptoSvc crypto.CryptoService, jwtSvc jwt.JWTService, contextTimeout time.Duration) *signinUsecase {
	return &signinUsecase{
		userRepo:       userRepo,
		cryptoSvc:      cryptoSvc,
		jwtSvc:         jwtSvc,
		contextTimeout: contextTimeout,
	}
}

func (u *signinUsecase) SignIn(c context.Context, request *request.SignInReq) (accessToken string, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewBadRequestError(constant.SigninMessage.SigninUserNotFound)
			return
		}
		return
	}

	if !u.cryptoSvc.ValidatePassword(ctx, user.Password, request.Password) {
		err = utils.NewBadRequestError(constant.SigninMessage.SigninEmailPasswordNotMatch)
		return
	}

	tokenVersion := utils.GenerateTokenVersion()
	err = u.userRepo.UpdateTokenVersionByID(ctx, tokenVersion, user.ID)
	if err != nil {
		errorInfo := utils.ErrorInfo{
			Message: constant.SigninMessage.SigninUpdateTokenVersionFailed,
			Details: err.Error(),
		}
		err = utils.NewBadRequestError(errorInfo)
		return
	}

	accessToken, err = u.jwtSvc.GenerateToken(ctx, user.ID, tokenVersion)
	return
}
