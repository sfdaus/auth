package usecase

import (
	"context"
	"database/sql"
	"time"

	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/utils"
	"prakarsa-app/utils/crypto"
)

type signupUsecase struct {
	userRepo       domain.UserRepository
	cryptoSvc      crypto.CryptoService
	contextTimeout time.Duration
}

func SignUpUsecase(userRepo domain.UserRepository, cryptoSvc crypto.CryptoService, contextTimeout time.Duration) *signupUsecase {
	return &signupUsecase{
		userRepo:       userRepo,
		cryptoSvc:      cryptoSvc,
		contextTimeout: contextTimeout,
	}
}

func (u *signupUsecase) SignUp(c context.Context, request *request.SignUpReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, request.Email)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if user.ID != "" {
		err = utils.NewBadRequestError("email already registered")
		return
	}

	passwordHash, err := u.cryptoSvc.CreatePasswordHash(ctx, request.Password)
	if err != nil {
		return
	}

	err = u.userRepo.Create(ctx, &domain.User{
		ID:           utils.GenerateUUID(),
		Email:        request.Email,
		PhoneNumber:  request.PhoneNumber,
		Username:     request.Username,
		TokenVersion: utils.GenerateTokenVersion(),
		IsActive:     true,
		Password:     passwordHash,
		CreatedAt:    time.Now().Unix(),
		CreatedBy:    "",
		UpdatedAt:    0,
		UpdatedBy:    "",
		DeletedAt:    0,
	})
	return
}
