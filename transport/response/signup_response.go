package response

type SignUpResponse struct {
	BasicResponse
	Data SignUpResponseData `json:"data,omitempty"`
}

type SignUpResponseData struct {
	UserID string `json:"user_id"`
}

type SignUpVerificationResponse struct {
	BasicResponse
	Data SignUpVerificationResponseData `json:"data,omitempty"`
}

type SignUpVerificationResponseData struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}
