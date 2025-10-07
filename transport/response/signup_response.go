package response

import "time"

type SignUpResponse struct {
	BasicResponse
	Data SignUpResponseData `json:"data,omitempty"`
}

type SignUpResponseData struct {
	UserID      string    `json:"user_id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	AvailableAt time.Time `json:"available_at"`
	ResendCount int64     `json:"resend_count"`
}

type SignUpVerificationResponse struct {
	BasicResponse
	Data SignUpVerificationResponseData `json:"data,omitempty"`
}

type SignUpVerificationResponseData struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}

type SignUpLimiterData struct {
	UserID      string    `json:"user_id"`
	AvailableAt time.Time `json:"available_at"`
	ResendCount int64     `json:"resend_count"`
}
