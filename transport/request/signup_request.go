package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type SignUpReq struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (request SignUpReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(
			&request.Username,
			// validation.By(func(value interface{}) error {
			// 	if request.Email == "" && request.PhoneNumber == "" {
			// 		return utils.ErrSignUpMissingEmailOrPhone
			// 	}
			// 	return nil
			// }),
			validation.Length(3, 50),
		),
		validation.Field(
			&request.Email,
			// validation.By(func(value interface{}) error {
			// 	if request.PhoneNumber == "" {
			// 		return utils.ErrSignUpMissingEmailOrPhone
			// 	}
			// 	return nil
			// }),
			is.Email,
		),
		validation.Field(
			&request.PhoneNumber,
			// validation.By(func(value interface{}) error {
			// 	if request.Email == "" {
			// 		return utils.ErrSignUpMissingEmailOrPhone
			// 	}
			// 	return nil
			// }),
			validation.Length(10, 15),
		),
		validation.Field(&request.Password, validation.Required, validation.Length(8, 20)),
	)
}
