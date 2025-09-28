package request

import validation "github.com/go-ozzo/ozzo-validation"

type UserProfileByIDReq struct {
	PublicID string `param:"public_id"`
	UserID   string
}

func (request UserProfileByIDReq) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(
			&request.UserID,
			validation.Required.Error("User ID is required"),
		),
		validation.Field(
			&request.PublicID,
			validation.Required.Error("Public ID is required"),
		),
	)
}
