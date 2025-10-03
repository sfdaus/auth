package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"prakarsa-app/config"
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
}

func SignUpUsecase(userRepo domain.UserRepository, authTokenRepo domain.AuthTokenRepository, profileRepo domain.ProfileRepository, cryptoSvc crypto.CryptoService, jwtSvc jwt.JWTService, contextTimeout time.Duration, emailService service.EmailService) *signupUsecase {
	return &signupUsecase{
		userRepo:       userRepo,
		authTokenRepo:  authTokenRepo,
		profileRepo:    profileRepo,
		cryptoSvc:      cryptoSvc,
		jwtSvc:         jwtSvc,
		contextTimeout: contextTimeout,
		emailService:   emailService,
	}
}

func (u *signupUsecase) SignUp(c context.Context, request *request.SignUpReq) (accessToken string, userID string, err error) {
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

	verificationToken, err := jwt.GenerateShortJWT(userID)
	htmlBody, err := utils.RenderTemplate("templates/verify_email.html", map[string]interface{}{
		"Name":      request.Name,
		"VerifyURL": fmt.Sprintf("%s/%s/%s", config.LoadConfig().BaseURLApp, "api/v1/auth/verify-account", verificationToken),
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

	userID, err := jwt.ValidateShortJWT(request.Token)
	if err != nil {
		return
	}
	err = u.userRepo.VerifyAccountByUserID(ctx, userID)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	tokenVersion := utils.GenerateTokenVersion()
	accessToken, err = u.jwtSvc.GenerateToken(ctx, userID, tokenVersion)
	if err != nil {
		err = utils.NewBadRequestError(constant.SignupMessage.SignupFailedToGenerateToken)
		return
	}
	return
}
