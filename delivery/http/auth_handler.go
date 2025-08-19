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
	SignUpUC  domain.SignUpUsecase
	SignInUC  domain.SignInUsecase
	ProfileUC domain.ProfileUsecase
}

// NewAuthHandler will initialize the auth resources endpoint
func NewAuthHandler(e *echo.Echo, middleware *middleware.Middleware, signUpUC domain.SignUpUsecase, signInUC domain.SignInUsecase, profileUC domain.ProfileUsecase) {
	handler := &AuthHandler{
		SignUpUC:  signUpUC,
		SignInUC:  signInUC,
		ProfileUC: profileUC,
	}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/auth/signup", handler.SignUp)
	apiV1.POST("/auth/signin", handler.SignIn)
	apiV1.POST("/auth/complete-profile", handler.CompleteProfile)
	apiV1.GET("/auth/profile-completion", handler.ProfileCompletion)
	apiV1.GET("/auth/user-profile", handler.UserProfile)
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

	accessToken, err := h.SignUpUC.SignUp(ctx, &req)
	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.SignupMessage.SignupFailed))
	}

	return c.JSON(http.StatusOK, response.SignUpResponse{
		BasicResponse: response.BasicResponse{
			Status:  constant.Status.Success,
			Message: constant.SignupMessage.SignupSuccess,
		},
		Data: response.SignUpResponseData{
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

	accessToken, err := h.SignInUC.SignIn(ctx, &req)

	if err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.SigninMessage.SignininFailed))
	}

	return c.JSON(http.StatusOK, response.SignInResponse{
		BasicResponse: response.BasicResponse{
			Status:  constant.Status.Success,
			Message: constant.SigninMessage.SigninSuccess,
		},
		Data: response.SignInResponseData{
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
