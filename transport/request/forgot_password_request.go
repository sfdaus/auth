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
