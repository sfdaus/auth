package http

import (
	"net/http"

	"prakarsa-app/config/constant"
	"prakarsa-app/delivery/middleware"
	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
	"prakarsa-app/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	SignUpUC         domain.SignUpUsecase
	SignInUC         domain.SignInUsecase
	ProfileUC        domain.ProfileUsecase
	ForgotPasswordUC domain.ForgotPassword
}

// NewAuthHandler will initialize the auth resources endpoint
func NewAuthHandler(e *echo.Echo, middleware *middleware.Middleware, signUpUC domain.SignUpUsecase, signInUC domain.SignInUsecase,
	profileUC domain.ProfileUsecase, forgotPasswordUC domain.ForgotPassword) {
	handler := &AuthHandler{
		SignUpUC:         signUpUC,
		SignInUC:         signInUC,
		ProfileUC:        profileUC,
		ForgotPasswordUC: forgotPasswordUC,
	}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/auth/signup", handler.SignUp)
	apiV1.POST("/auth/signin", handler.SignIn)
	apiV1.POST("/auth/complete-profile", handler.CompleteProfile)
	apiV1.GET("/auth/profile-completion", handler.ProfileCompletion)
	apiV1.GET("/auth/user-profile", handler.UserProfile)
	apiV1.PUT("/auth/update-profile", handler.UpdateProfile)
	apiV1.GET("/auth/:public_id/user-profile", handler.UserProfileByID)
	apiV1.POST("/auth/forgot-password", handler.ForgotPassword)
	apiV1.POST("/auth/reset-password/verify", handler.VerifyResetPassword)
	apiV1.POST("/auth/reset-password", handler.ResetPassword)
	apiV1.POST("/auth/verify-account", handler.VerifyAccount)
	apiV1.POST("/auth/verify-account/resend", handler.VerifyAccountResend)
}

// SignUp godoc
// @Summary SignUp
// @Description SignUp
// @Tags Auth
// @Accept json
// @Produce json
// @Param signup body request.SignUpReq true "SignUp user"
// @Success 200
// @Router /api/v1/auth/signup [post]
func (h *AuthHandler) SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.SignUpReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.SignupMessage.SignupFailed,
			Error:   err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.SignupMessage.SignupFailed,
			Error:   errVal,
		})
	}

	userID, err := h.SignUpUC.SignUp(ctx, &req)
	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.SignupMessage.SignupFailed))
	}

	return c.JSON(http.StatusOK, response.SignUpResponse{
		BasicResponse: response.BasicResponse{
			Status:  constant.Status.Success,
			Message: constant.SignupMessage.SignupSuccess,
		},
		Data: response.SignUpResponseData{
			UserID: userID,
		},
	})
}

// Verify Account Resend godoc
// @Summary Verify Account Resend
// @Description Verify Account Resend
// @Tags Auth
// @Accept json
// @Produce json
// @Param token body request.VerifyAccountReq true "Verify account resend"
// @Success 200
// @Router /api/v1/auth/verify-account/resend [post]
func (h *AuthHandler) VerifyAccountResend(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.VerifyAccountResendReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.SignupMessage.SignupFailed,
			Error:   err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.SignupMessage.SignupFailed,
			Error:   errVal,
		})
	}

	userID, err := h.SignUpUC.VerifyAccountResend(ctx, &req)
	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.SignupMessage.SignupFailed))
	}

	return c.JSON(http.StatusOK, response.SignUpResponse{
		BasicResponse: response.BasicResponse{
			Status:  constant.Status.Success,
			Message: constant.SignupMessage.SignupSuccess,
		},
		Data: response.SignUpResponseData{
			UserID: userID,
		},
	})
}

// Verify Account godoc
// @Summary Verify Account
// @Description Verify Account
// @Tags Auth
// @Accept json
// @Produce json
// @Param token body request.VerifyAccountReq true "Verify account"
// @Success 200
// @Router /api/v1/auth/verify-account [post]
func (h *AuthHandler) VerifyAccount(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.VerifyAccountReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.SignupMessage.SignupFailed,
			Error:   err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.SignupMessage.SignupFailed,
			Error:   errVal,
		})
	}

	accessToken, userID, err := h.SignUpUC.VerifyAccount(ctx, &req)
	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.SignupMessage.SignupFailed))
	}

	return c.JSON(http.StatusOK, response.SignUpVerificationResponse{
		BasicResponse: response.BasicResponse{
			Status:  constant.Status.Success,
			Message: constant.SignupMessage.SignupSuccess,
		},
		Data: response.SignUpVerificationResponseData{
			UserID:      userID,
			AccessToken: accessToken,
		},
	})
}

// SignIn godoc
// @Summary SignIn
// @Description SignIn
// @Tags Auth
// @Accept json
// @Produce json
// @Param signin body request.SignInReq true "SignIn user"
// @Success 200
// @Router /api/v1/auth/signin [post]
func (h *AuthHandler) SignIn(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.SignInReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.SigninMessage.SignininFailed,
			Error:   err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.SigninMessage.SignininFailed,
			Error:   errVal,
		})
	}

	accessToken, userID, err := h.SignInUC.SignIn(ctx, &req)

	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.SigninMessage.SignininFailed))
	}

	return c.JSON(http.StatusOK, response.SignInResponse{
		BasicResponse: response.BasicResponse{
			Status:  constant.Status.Success,
			Message: constant.SigninMessage.SigninSuccess,
		},
		Data: response.SignInResponseData{
			UserID:      userID,
			AccessToken: accessToken,
		},
	})
}

// CompleteProfile godoc
// @Summary CompleteProfile
// @Description CompleteProfile
// @Tags Auth
// @Accept json
// @Produce json
// @Param x-user-id header string true "User ID from Gateway"
// @Param completeprofile body request.CompleteProfileReq true "CompleteProfile user"
// @Success 200
// @Router /api/v1/auth/complete-profile [post]
func (h *AuthHandler) CompleteProfile(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.CompleteProfileReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.CompleteProfileMessage.CompleteProfileFailed,
			Error:   err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.CompleteProfileMessage.CompleteProfileFailed,
			Error:   errVal,
		})
	}

	userID := c.Request().Header.Get("x-user-id")
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, response.BasicResponse{
			Status:  constant.Status.Error,
			Message: constant.CompleteProfileMessage.CompleteProfileFailed,
			Error:   constant.AuthorizationMessage.AuthorizationXUserIDMissing,
		})
	}
	err := h.ProfileUC.CompleteProfile(ctx, userID, &req)
	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.CompleteProfileMessage.CompleteProfileFailed))
	}

	return c.JSON(http.StatusOK, response.BasicResponse{
		Status:  constant.Status.Success,
		Message: constant.CompleteProfileMessage.CompleteProfileSuccess,
	})
}

// ProfileCompletion godoc
// @Summary ProfileCompletion
// @Description ProfileCompletion
// @Tags Auth
// @Accept json
// @Produce json
// @Param x-user-id header string true "User ID from Gateway"
// @Success 200
// @Router /api/v1/auth/profile-completion [get]
func (h *AuthHandler) ProfileCompletion(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Request().Header.Get("x-user-id")
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, response.BasicResponse{
			Status:  constant.Status.Error,
			Message: constant.ProfileCompletionMessage.ProfileCompletionFailed,
			Error:   constant.AuthorizationMessage.AuthorizationXUserIDMissing,
		})
	}
	profileCompletion, err := h.ProfileUC.ProfileCompletion(ctx, userID)
	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.ProfileCompletionMessage.ProfileCompletionFailed))
	}

	return c.JSON(http.StatusOK, response.ProfileCompletionResponse{
		BasicResponse: response.BasicResponse{
			Status:  constant.Status.Success,
			Message: constant.ProfileCompletionMessage.ProfileCompletionSuccess,
		},
		Data: response.ProfileCompletionResponseData{
			CompletionStatus: profileCompletion,
		},
	})
}

// UserProfile godoc
// @Summary UserProfile
// @Description UserProfile
// @Tags Auth
// @Accept json
// @Produce json
// @Param x-user-id header string true "User ID from Gateway"
// @Success 200
// @Router /api/v1/auth/user-profile [get]
func (h *AuthHandler) UserProfile(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Request().Header.Get("x-user-id")
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, response.BasicResponse{
			Status:  constant.Status.Error,
			Message: constant.UserProfileMessage.UserProfileFailed,
			Error:   constant.AuthorizationMessage.AuthorizationXUserIDMissing,
		})
	}
	userProfile, err := h.ProfileUC.UserProfile(ctx, userID)
	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.UserProfileMessage.UserProfileFailed))
	}

	return c.JSON(http.StatusOK, response.UserProfileResponse{
		BasicResponse: response.BasicResponse{
			Status:  constant.Status.Success,
			Message: constant.UserProfileMessage.UserProfileSuccess,
		},
		Data: userProfile,
	})
}

// UpdateProfile godoc
// @Summary UpdateProfile
// @Description UpdateProfile
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param x-user-id header string true "User ID from Gateway"
// @Param name formData string false "Name"
// @Param avatar formData file false "Avatar file"
// @Param gender formData string false "Gender"
// @Param birth_date formData string false "Birth Date"
// @Param slug_name formData string false "Slug Name"
// @Param about_me formData string false "About Me"
// @Param institution_id formData string false "Institution ID"
// @Param linkedin formData string false "Linkedin URL"
// @Success 200
// @Router /api/v1/auth/update-profile [put]
func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.UpdateProfileReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.UpdateProfileMessage.UpdateProfileFailed,
			Error:   err.Error(),
		})
	}

	fileHeader, err := c.FormFile("avatar")
	if err == nil { // kalau ada file
		req.Avatar = fileHeader
	}

	userID := c.Request().Header.Get("x-user-id")
	if userID == "" {
		return c.JSON(http.StatusUnauthorized, response.BasicResponse{
			Status:  constant.Status.Error,
			Message: constant.UpdateProfileMessage.UpdateProfileFailed,
			Error:   constant.AuthorizationMessage.AuthorizationXUserIDMissing,
		})
	}
	err = h.ProfileUC.UpdateProfile(ctx, userID, &req)
	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.UpdateProfileMessage.UpdateProfileFailed))
	}

	return c.JSON(http.StatusOK, response.BasicResponse{
		Status:  constant.Status.Success,
		Message: constant.UpdateProfileMessage.UpdateProfileSuccess,
	})
}

// UserProfileByID godoc
// @Summary UserProfileByID
// @Description UserProfileByID
// @Tags Auth
// @Accept json
// @Produce json
// @Param x-user-id header string true "User ID from Gateway"
// @Success 200
// @Router /api/v1/auth/:id/user-profile [get]
func (h *AuthHandler) UserProfileByID(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.UserProfileByIDReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	req.UserID = c.Request().Header.Get("x-user-id")

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	userProfile, err := h.ProfileUC.UserProfileByID(ctx, &req)
	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.UserProfileMessage.UserProfileFailed))
	}

	return c.JSON(http.StatusOK, response.UserProfileResponse{
		BasicResponse: response.BasicResponse{
			Status:  constant.Status.Success,
			Message: constant.UserProfileMessage.UserProfileSuccess,
		},
		Data: userProfile,
	})
}

// ForgotPassword godoc
// @Summary ForgotPassword
// @Description ForgotPassword
// @Tags Auth
// @Accept json
// @Produce json
// @Param forgot_password body request.ForgotPasswordReq true "Forgot Password user"
// @Success 202
// @Router /api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var req request.ForgotPasswordReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.ForgotPasswordMessage.ForgotPasswordFailed,
			Error:   err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.ForgotPasswordMessage.ForgotPasswordFailed,
			Error:   errVal,
		})
	}

	err = h.ForgotPasswordUC.ForgotPassword(ctx, &req)

	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.ForgotPasswordMessage.ForgotPasswordFailed))
	}

	return c.JSON(http.StatusAccepted, response.ForgotPasswordResponse{
		BasicResponse: response.BasicResponse{
			Status:  constant.Status.Success,
			Message: constant.ForgotPasswordMessage.ForgotPasswordSuccess,
		},
	})
}

// VerifyResetPassword godoc
// @Summary VerifyResetPassword
// @Description VerifyResetPassword
// @Tags Auth
// @Accept json
// @Produce json
// @Param verify_reset_password body request.VerifyResetPasswordReq true "Verify Reset Password user"
// @Success 200
// @Router /api/v1/auth/reset-password/verify [post]
func (h *AuthHandler) VerifyResetPassword(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var req request.VerifyResetPasswordReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.ForgotPasswordMessage.VerifyResetPasswordFailed,
			Error:   err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.ForgotPasswordMessage.VerifyResetPasswordFailed,
			Error:   errVal,
		})
	}

	valid, err := h.ForgotPasswordUC.VerifyResetPassword(ctx, &req)

	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.ForgotPasswordMessage.VerifyResetPasswordFailed))
	}

	return c.JSON(http.StatusOK, response.VerifyResetPasswordResponse{
		BasicResponse: response.BasicResponse{
			Status:  constant.Status.Success,
			Message: constant.ForgotPasswordMessage.VerifyResetPasswordSuccess,
		},
		Data: response.VerifyResetPasswordResponseData{
			Valid: valid,
		},
	})
}

// ResetPassword godoc
// @Summary ResetPassword
// @Description ResetPassword
// @Tags Auth
// @Accept json
// @Produce json
// @Param reset_password body request.ResetPasswordReq true "Reset Password user"
// @Success 200
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var req request.ResetPasswordReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.ForgotPasswordMessage.ResetPasswordFailed,
			Error:   err.Error(),
		})
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, response.BasicResponse{
			Status:  constant.Status.Failed,
			Message: constant.ForgotPasswordMessage.ResetPasswordFailed,
			Error:   errVal,
		})
	}

	err = h.ForgotPasswordUC.ResetPassword(ctx, &req)

	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.ForgotPasswordMessage.ResetPasswordFailed))
	}

	return c.JSON(http.StatusOK, response.ResetPasswordResponse{
		BasicResponse: response.BasicResponse{
			Status:  constant.Status.Success,
			Message: constant.ForgotPasswordMessage.ResetPasswordSuccess,
		},
	})
}
