package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"prakarsa-app/config/constant"
	"prakarsa-app/domain"
	"prakarsa-app/repository/redis"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
	"prakarsa-app/utils"
	"prakarsa-app/utils/crypto"
	"prakarsa-app/utils/jwt"
	"time"
)

type forgotpasswordUsecase struct {
	userRepo       domain.UserRepository
	redisRepo      redis.RedisRepository
	cryptoSvc      crypto.CryptoService
	jwtSvc         jwt.JWTService
	contextTimeout time.Duration
}

func ForgotPasswordUsecase(userRepo domain.UserRepository, redisRepo redis.RedisRepository, cryptoSvc crypto.CryptoService, jwtSvc jwt.JWTService, contextTimeout time.Duration) *forgotpasswordUsecase {
	return &forgotpasswordUsecase{
		userRepo:       userRepo,
		redisRepo:      redisRepo,
		cryptoSvc:      cryptoSvc,
		jwtSvc:         jwtSvc,
		contextTimeout: contextTimeout,
	}
}

func (u *forgotpasswordUsecase) ForgotPassword(c context.Context, request *request.ForgotPasswordReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return
	}

	_, jti, err := u.jwtSvc.GenerateForgotPasswordToken(ctx, user.ID)

	// Write on redis
	b, _ := json.Marshal(map[string]string{
		"user_id": user.ID,
	})

	err = u.redisRepo.Set(fmt.Sprintf("%s:token:%s", constant.FORGOT_PASSWORD_REDIS_KEY, jti),
		b,
		constant.FORGOT_PASSWORD_EXPIRES)
	if err != nil {
		return
	}

	return
}

func (u *forgotpasswordUsecase) VerifyResetPassword(c context.Context, request *request.VerifyResetPasswordReq) (valid bool, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Decode JWT Token
	claim, err := u.jwtSvc.ParseForgotPasswordToken(c, request.Token)
	if err != nil {
		return false, err
	}

	// Check redis
	redisValue, err := u.redisRepo.Get(fmt.Sprintf("%s:token:%s", constant.FORGOT_PASSWORD_REDIS_KEY, claim.Id))
	if err != nil {
		return false, utils.NewNotFoundError(constant.ForgotPasswordMessage.VerifyResetPasswordTokenNotValid)
	}

	var p response.ResetPayload
	if err := json.Unmarshal([]byte(redisValue), &p); err != nil {
		return false, fmt.Errorf("unmarshal redis value: %w", err)
	}

	if p.UserID != claim.UserID {
		return false, utils.NewUnauthorizedError("Token is not valid: User ID mismatch")
	}

	// Find user from token
	_, err = u.userRepo.GetByUserID(ctx, claim.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, utils.NewNotFoundError(constant.ForgotPasswordMessage.VerifyResetPasswordTokenNotValid)
		}
		return false, err
	}

	valid = true

	return
}

func (u *forgotpasswordUsecase) ResetPassword(c context.Context, request *request.ResetPasswordReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Decode JWT Token
	claim, err := u.jwtSvc.ParseForgotPasswordToken(c, request.Token)
	if err != nil {
		return err
	}

	// Check redis
	redisValue, err := u.redisRepo.Get(fmt.Sprintf("%s:token:%s", constant.FORGOT_PASSWORD_REDIS_KEY, claim.Id))
	if err != nil {
		return utils.NewNotFoundError(constant.ForgotPasswordMessage.VerifyResetPasswordTokenNotValid)
	}

	var p response.ResetPayload
	if err := json.Unmarshal([]byte(redisValue), &p); err != nil {
		return fmt.Errorf("unmarshal redis value: %w", err)
	}

	if p.UserID != claim.UserID {
		return utils.NewUnauthorizedError("Token is not valid: User ID mismatch")
	}

	// Check password and verify_password valid
	if request.Password != request.VerifyPassword {
		return utils.NewBadRequestError("Passwords do not match")
	}

	passwordHash, err := u.cryptoSvc.CreatePasswordHash(ctx, request.Password)
	if err != nil {
		return
	}

	// Update password
	err = u.userRepo.UpdatePasswordByUserID(ctx, claim.UserID, passwordHash)
	if err != nil {
		return
	}

	// Delete redis
	_ = u.redisRepo.Del(fmt.Sprintf("%s:token:%s", constant.FORGOT_PASSWORD_REDIS_KEY, claim.Id))

	return
}
