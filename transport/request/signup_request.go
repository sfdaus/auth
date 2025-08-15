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
	return validation.ValidateStruct(
		validation.Field(
			&request.Email,
			validation.By(func(value interface{}) error {
				if request.Email == "" {
					return utils.ErrSignUpMissingEmailOrPhone
				}
				return nil
			}),
			is.Email,
		),
		validation.Field(&request.Password, validation.Required, validation.Length(8, 20)),
	)
}
