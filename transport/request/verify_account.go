package request

import validation "github.com/go-ozzo/ozzo-validation"

type VerifyAccount struct {
	UserID string `json:"user_id"`
}

func (request VerifyAccount) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(
			&request.UserID,
			validation.Required.Error("userID is required"),
		),
	)
}
