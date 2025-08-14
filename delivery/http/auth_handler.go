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
	SignUpUC domain.SignUpUsecase
	SignInUC domain.SignInUsecase
}

// NewAuthHandler will initialize the auth resources endpoint
func NewAuthHandler(e *echo.Echo, middleware *middleware.Middleware, signUpUC domain.SignUpUsecase, signInUC domain.SignInUsecase) {
	handler := &AuthHandler{
		SignUpUC: signUpUC,
		SignInUC: signInUC,
	}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/auth/signup", handler.SignUp)
	apiV1.POST("/auth/signin", handler.SignIn)
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

	if err := h.SignUpUC.SignUp(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpErrorToBasicResponse(err, constant.SignupMessage.SignupFailed))
	}

	return c.JSON(http.StatusOK, response.BasicResponse{
		Status:  constant.Status.Success,
		Message: constant.SignupMessage.SignupSuccess,
		Data:    nil,
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

	return c.JSON(http.StatusOK, response.BasicResponse{
		Status:  constant.Status.Success,
		Message: constant.SigninMessage.SigninSuccess,
		Data: map[string]interface{}{
			"access_token": accessToken,
		},
	})
}
