package response

type UserProfileResponse struct {
	BasicResponse
	Data any `json:"data,omitempty"`
}
