package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"prakarsa-app/utils"
)

type ForgotPasswordReq struct {
	Email string `json:"email"`
}

func (request ForgotPasswordReq) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(
			&request.Email,
			validation.Required.Error(utils.ErrForgotPasswordMissingEmail.Error()),
			is.Email.Error("invalid email format"),
		),
	)
}

type VerifyResetPasswordReq struct {
	Token string `json:"token"`
}

func (request VerifyResetPasswordReq) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(
			&request.Token,
			validation.Required.Error(utils.ErrForgotPasswordMissingToken.Error()),
		),
	)
}
