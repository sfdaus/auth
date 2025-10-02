package response

type ForgotPasswordResponse struct {
	BasicResponse
	Data ForgotPasswordResponseData `json:"data,omitempty"`
}

type ForgotPasswordResponseData struct {
}
