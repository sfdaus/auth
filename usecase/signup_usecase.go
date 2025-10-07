package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"prakarsa-app/config"
	"prakarsa-app/repository/redis"
	"prakarsa-app/transport/response"
	"time"

	"prakarsa-app/config/constant"
	"prakarsa-app/domain"
	service "prakarsa-app/service/mail"
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
	emailService   service.EmailService
	redisRepo      redis.RedisRepository
}

func SignUpUsecase(userRepo domain.UserRepository, authTokenRepo domain.AuthTokenRepository, profileRepo domain.ProfileRepository,
	cryptoSvc crypto.CryptoService, jwtSvc jwt.JWTService, contextTimeout time.Duration, emailService service.EmailService,
	redisRepo redis.RedisRepository) *signupUsecase {
	return &signupUsecase{
		userRepo:       userRepo,
		authTokenRepo:  authTokenRepo,
		profileRepo:    profileRepo,
		cryptoSvc:      cryptoSvc,
		jwtSvc:         jwtSvc,
		contextTimeout: contextTimeout,
		emailService:   emailService,
		redisRepo:      redisRepo,
	}
}

func (u *signupUsecase) SignUp(c context.Context, request *request.SignUpReq) (res response.SignUpResponseData, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	encryptedMail, err := utils.EncryptDeterministic(request.Email)
	if err != nil {
		err = utils.NewBadRequestError(constant.SignupMessage.SignupFailedEncryptEmail)
		return
	}

	user, err := u.userRepo.GetByEmail(ctx, encryptedMail)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if user.ID != "" && user.IsVerified {
		err = utils.NewBadRequestError(constant.SignupMessage.SignupExists)
		return
	}

	passwordHash, err := u.cryptoSvc.CreatePasswordHash(ctx, request.Password)
	if err != nil {
		return
	}

	tokenVersion := utils.GenerateTokenVersion()
	newUserID := utils.GenerateUUID()
	nameAlias := utils.GenerateNameAlias(request.Name)

	// response builder
	if user.ID != "" {
		res.UserID = user.ID
	} else {
		res.UserID = newUserID
	}
	res.Email = request.Email
	res.Name = request.Name

	var empty = ""
	var userPayload = domain.User{
		ID:           newUserID,
		Email:        encryptedMail,
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
		IsVerified:   false,
	}

	var profilePayload = domain.Profile{
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
	}

	var authTokenPayload = domain.AuthToken{
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
	}

	if user.ID != "" {
		// Update user
		err = u.userRepo.Update(ctx, &userPayload)
		if err != nil {
			errorInfo := utils.ErrorInfo{
				Message: constant.SignupMessage.SignupFailedToUpdateUser,
				Details: err.Error(),
			}
			err = utils.NewBadRequestError(errorInfo)
			return
		}

		// Update profile
		updateProfilePayload := domain.UpdateProfile{
			Name:      request.Name,
			NameAlias: nameAlias,
			UpdatedAt: time.Now().Unix(),
			UpdatedBy: "",
		}
		err = u.profileRepo.UpdateProfile(ctx, user.ID, &updateProfilePayload)
		if err != nil {
			errorInfo := utils.ErrorInfo{
				Message: constant.SignupMessage.SignupFailedToUpdateProfile,
				Details: err.Error(),
			}
			err = utils.NewBadRequestError(errorInfo)
			return
		}

		// Update auth token for the new user
		err = u.userRepo.UpdateTokenVersionByID(ctx, tokenVersion, user.ID)
		if err != nil {
			errorInfo := utils.ErrorInfo{
				Message: constant.SignupMessage.SignupFailedToUpdateAuthToken,
				Details: err.Error(),
			}
			err = utils.NewBadRequestError(errorInfo)
			return
		}
	} else {
		// Create new user
		err = u.userRepo.Create(ctx, &userPayload)
		if err != nil {
			errorInfo := utils.ErrorInfo{
				Message: constant.SignupMessage.SignupFailedToCreateUser,
				Details: err.Error(),
			}
			err = utils.NewBadRequestError(errorInfo)
			return
		}

		// Create profile for the new user
		err = u.profileRepo.Create(ctx, &profilePayload)
		if err != nil {
			errorInfo := utils.ErrorInfo{
				Message: constant.SignupMessage.SignupFailedToCreateProfile,
				Details: err.Error(),
			}
			err = utils.NewBadRequestError(errorInfo)
			return
		}

		// Create auth token for the new user
		err = u.authTokenRepo.Create(ctx, &authTokenPayload)
		if err != nil {
			errorInfo := utils.ErrorInfo{
				Message: constant.SignupMessage.SignupFailedToCreateAuthToken,
				Details: err.Error(),
			}
			err = utils.NewBadRequestError(errorInfo)
			return
		}
	}

	// Check Limiter on Redis
	redisValue, _ := u.redisRepo.Get(fmt.Sprintf("%s:%s", constant.SIGN_UP_LIMITER_REDIS_KEY, res.UserID))

	// Create Limiter for resend
	resendAvailableAt := time.Now().Add(constant.SIGN_UP_LIMITER_RESEND_AVAILABLE_TIME)
	resendCount := int64(0)

	res.AvailableAt = resendAvailableAt
	res.ResendCount = resendCount

	if redisValue != "" {
		var p response.SignUpLimiterData
		if err := json.Unmarshal([]byte(redisValue), &p); err != nil {
			return res, fmt.Errorf("unmarshal redis value: %w", err)
		}

		resendCount = p.ResendCount + 1

		res.AvailableAt = p.AvailableAt
		res.ResendCount = resendCount

		if p.ResendCount >= constant.MAX_SIGN_UP_RESEND_LIMIT {
			res.ResendCount = p.ResendCount
			return res, utils.NewTooManyRequestError(constant.SignupMessage.SignupFailedResendLimit)
		}

		//if p.AvailableAt.After(time.Now()) {
		//	return res, utils.NewTooManyRequestError(constant.SignupMessage.SignupFailedResendNotAvailable)
		//}

		if p.ResendCount == (constant.MAX_SIGN_UP_RESEND_LIMIT - 1) {
			resendAvailableAt = time.Now().Add(constant.SIGN_UP_LIMITER_ON_LIMIT_EXPIRES)
		} else {
			resendAvailableAt = time.Now().Add(constant.SIGN_UP_LIMITER_RESEND_ON_LIMIT_AVAILABLE_TIME)
		}
	}

	b, _ := json.Marshal(map[string]interface{}{
		"user_id":      res.UserID,
		"available_at": resendAvailableAt,
		"resend_count": resendCount,
	})

	err = u.redisRepo.Set(fmt.Sprintf("%s:%s", constant.SIGN_UP_LIMITER_REDIS_KEY, res.UserID),
		b,
		constant.SIGN_UP_LIMITER_ON_LIMIT_EXPIRES)
	if err != nil {
		return
	}

	// Generate JWT untuk verif dan send email
	verificationToken, err := jwt.GenerateShortJWT(newUserID)
	htmlBody, err := utils.RenderTemplate("templates/verify_email.html", map[string]interface{}{
		"Name":      request.Name,
		"VerifyURL": fmt.Sprintf("%s/%s%s", config.LoadConfig().BaseURLApp, "auth/verify-account?token=", verificationToken),
	})
	if err != nil {
		err = utils.NewBadRequestError("Err")
		return
	}
	err = u.emailService.SendEmail(
		ctx,
		request.Email,
		"Verifikasi Email",
		htmlBody,
	)
	if err != nil {
		err = utils.NewBadRequestError(constant.SignupMessage.SignupFailedToSendEmail)
	}

	return
}

func (u *signupUsecase) VerifyAccount(c context.Context, request *request.VerifyAccountReq) (accessToken string, userId string, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	userId, err = jwt.ValidateShortJWT(request.Token)
	if err != nil {
		return
	}
	err = u.userRepo.VerifyAccountByUserID(ctx, userId)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	tokenVersion := utils.GenerateTokenVersion()
	accessToken, err = u.jwtSvc.GenerateToken(ctx, userId, tokenVersion)
	if err != nil {
		err = utils.NewBadRequestError(constant.SignupMessage.SignupFailedToGenerateToken)
		return
	}
	return
}

func (u *signupUsecase) VerifyAccountResend(c context.Context, request *request.VerifyAccountResendReq) (res response.SignUpResponseData, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Get User
	encryptedMail, err := utils.EncryptDeterministic(request.Email)
	if err != nil {
		err = utils.NewBadRequestError(constant.SignupMessage.SignupFailedEncryptEmail)
		return
	}

	user, err := u.userRepo.GetByEmail(ctx, encryptedMail)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if user.ID == "" {
		err = utils.NewNotFoundError(constant.SignupMessage.SignupFailed)
		return
	}

	if user.IsVerified {
		err = utils.NewNotFoundError(constant.SignupMessage.SignupFailed)
		return
	}

	// response builder
	res.UserID = user.ID
	res.Email = request.Email
	res.Name = request.Name

	// Check Limiter on Redis
	redisValue, err := u.redisRepo.Get(fmt.Sprintf("%s:%s", constant.SIGN_UP_LIMITER_REDIS_KEY, user.ID))
	if err != nil {
		return res, utils.NewNotFoundError(constant.SignupMessage.SignupFailed)
	}

	var p response.SignUpLimiterData
	if err := json.Unmarshal([]byte(redisValue), &p); err != nil {
		return res, fmt.Errorf("unmarshal redis value: %w", err)
	}

	resendCount := p.ResendCount + 1

	// Handler limiter
	res.ResendCount = resendCount
	res.AvailableAt = p.AvailableAt

	if p.ResendCount >= constant.MAX_SIGN_UP_RESEND_LIMIT {
		res.ResendCount = p.ResendCount
		return res, utils.NewTooManyRequestError(constant.SignupMessage.SignupFailedResendLimit)
	}

	//if p.AvailableAt.After(time.Now()) {
	//	return res, utils.NewTooManyRequestError(constant.SignupMessage.SignupFailedResendNotAvailable)
	//}

	// Create Limiter for resend
	var resendAvailableAt = time.Now().Add(constant.SIGN_UP_LIMITER_RESEND_AVAILABLE_TIME)

	if p.ResendCount == (constant.MAX_SIGN_UP_RESEND_LIMIT - 1) {
		resendAvailableAt = time.Now().Add(constant.SIGN_UP_LIMITER_ON_LIMIT_EXPIRES)
	} else {
		resendAvailableAt = time.Now().Add(constant.SIGN_UP_LIMITER_RESEND_ON_LIMIT_AVAILABLE_TIME)
	}

	b, _ := json.Marshal(map[string]interface{}{
		"user_id":      user.ID,
		"available_at": resendAvailableAt,
		"resend_count": resendCount,
	})

	err = u.redisRepo.Set(fmt.Sprintf("%s:%s", constant.SIGN_UP_LIMITER_REDIS_KEY, user.ID),
		b,
		constant.SIGN_UP_LIMITER_ON_LIMIT_EXPIRES)
	if err != nil {
		return
	}

	// Generate Verifikasi Token
	verificationToken, err := jwt.GenerateShortJWT(user.ID)
	htmlBody, err := utils.RenderTemplate("templates/verify_email.html", map[string]interface{}{
		"Name":      request.Name,
		"VerifyURL": fmt.Sprintf("%s/%s%s", config.LoadConfig().BaseURLApp, "auth/verify-account?token=", verificationToken),
	})
	if err != nil {
		err = utils.NewBadRequestError("Err")
		return
	}

	// Send Email
	err = u.emailService.SendEmail(
		ctx,
		request.Email,
		"Verifikasi Email",
		htmlBody,
	)
	if err != nil {
		err = utils.NewBadRequestError(constant.SignupMessage.SignupFailedToSendEmail)
	}

	return
}
