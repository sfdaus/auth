package response

type SignInResponse struct {
	BasicResponse
	Data SignInResponseData `json:"data,omitempty"`
}

type SignInResponseData struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}
