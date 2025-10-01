package request

import validation "github.com/go-ozzo/ozzo-validation"

type CompleteProfileReq struct {
	BirthDate     string `json:"birth_date"`
	Gender        string `json:"gender"`
	InstitutionID string `json:"institution_id"`
}

func (request CompleteProfileReq) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(
			&request.Gender,
			validation.Required.Error("gender is required"),
			validation.In("M", "F").Error("gender not valid (M/F)"),
		),
	)
}
