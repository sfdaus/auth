package request

import validation "github.com/go-ozzo/ozzo-validation"

type VerifyAccountReq struct {
	Token string `json:"token"`
}

func (request VerifyAccountReq) Validate() error {
	return validation.ValidateStruct(&request,
		validation.Field(
			&request.Token,
			validation.Required.Error("token is required"),
		),
	)
}
