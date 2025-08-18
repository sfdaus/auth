package response

type SignUpResponse struct {
	BasicResponse
	Data SignUpResponseData `json:"data,omitempty"`
}

type SignUpResponseData struct {
	AccessToken string `json:"access_token"`
}
