package response

type ForgotPasswordResponse struct {
	BasicResponse
	Data ForgotPasswordResponseData `json:"data,omitempty"`
}

type ForgotPasswordResponseData struct {
}

type VerifyResetPasswordResponse struct {
	BasicResponse
	Data VerifyResetPasswordResponseData `json:"data,omitempty"`
}

type VerifyResetPasswordResponseData struct {
	Valid bool `json:"valid"`
}

type ResetPayload struct {
	UserID string `json:"user_id"`
}

type ResetPasswordResponse struct {
	BasicResponse
	Data ResetPasswordResponseData `json:"data,omitempty"`
}

type ResetPasswordResponseData struct {
}
