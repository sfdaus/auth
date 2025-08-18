package request

import (
	"prakarsa-app/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type SignUpReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (request SignUpReq) Validate() error {
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
		validation.Field(
			&request.Password,
			validation.Required.Error("password is required"),
			validation.Length(8, 20).Error("password must be between 8-20 characters"),
		),
	)
}
