package response

type SignInResponse struct {
	BasicResponse
	Data SignInResponseData `json:"data,omitempty"`
}

type SignInResponseData struct {
	AccessToken string `json:"access_token"`
}
