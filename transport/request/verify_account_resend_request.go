package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"prakarsa-app/utils"
)

type VerifyAccountResendReq struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (request VerifyAccountResendReq) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(
			&request.Email,
			validation.Required.Error(utils.ErrSignUpMissingEmailOrPhone.Error()),
			is.Email.Error("invalid email format"),
		),
		validation.Field(
			&request.Name,
			validation.Required.Error("name is required"),
			validation.Length(3, 50).Error("name must be between 3-50 characters"),
		),
	)
}
