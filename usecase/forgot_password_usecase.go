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

	accessToken, jti, err := u.jwtSvc.GenerateForgotPasswordToken(ctx, user.ID)

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

	fmt.Println(accessToken, jti)

	return
}
