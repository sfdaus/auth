package constant

import "time"

type ResponseStatus struct {
	Success string
	Failed  string
	Error   string
}

var Status = ResponseStatus{
	Success: "success",
	Failed:  "failed",
	Error:   "error",
}

type SigninResponseMessage struct {
	SigninSuccess                  string
	SignininFailed                 string
	SigninUpdateTokenVersionFailed string
	SigninUserNotFound             string
	SigninEmailPasswordNotMatch    string
}

var SigninMessage = SigninResponseMessage{
	SigninSuccess:                  "Sign in successful",
	SignininFailed:                 "Sign in failed",
	SigninUpdateTokenVersionFailed: "Failed to update token version",
	SigninUserNotFound:             "User not found",
	SigninEmailPasswordNotMatch:    "Email and password do not match",
}

type SignupResponseMessage struct {
	SignupSuccess                 string
	SignupFailed                  string
	SignupExists                  string
	SignupFailedToCreateUser      string
	SignupFailedToCreateProfile   string
	SignupFailedToCreateAuthToken string
	SignupFailedToGenerateToken   string
}

var SignupMessage = SignupResponseMessage{
	SignupSuccess:                 "Sign up successful",
	SignupFailed:                  "Sign up failed",
	SignupExists:                  "Email/phone number/username already registered",
	SignupFailedToCreateUser:      "Failed to create user",
	SignupFailedToCreateProfile:   "Failed to create profile",
	SignupFailedToCreateAuthToken: "Failed to create auth token",
	SignupFailedToGenerateToken:   "Failed to generate access token",
}

type CompleteProfileResponseMessage struct {
	CompleteProfileSuccess string
	CompleteProfileFailed  string
}

var CompleteProfileMessage = CompleteProfileResponseMessage{
	CompleteProfileSuccess: "Complete profile successful",
	CompleteProfileFailed:  "Complete profile failed",
}

type ProfileCompletionResponseMessage struct {
	ProfileCompletionSuccess      string
	ProfileCompletionFailed       string
	ProfileCompletionUserNotFound string
}

var ProfileCompletionMessage = ProfileCompletionResponseMessage{
	ProfileCompletionSuccess:      "Profile completion successful",
	ProfileCompletionFailed:       "Profile completion failed",
	ProfileCompletionUserNotFound: "User not found",
}

type AuthorizationResponseMessage struct {
	AuthorizationSuccess        string
	AuthorizationFailed         string
	AuthorizationXUserIDMissing string
}

var AuthorizationMessage = AuthorizationResponseMessage{
	AuthorizationSuccess:        "Authorization successful",
	AuthorizationFailed:         "Authorization failed",
	AuthorizationXUserIDMissing: "Missing x-user-id header",
}

type UserProfileResponseMessage struct {
	UserProfileSuccess      string
	UserProfileFailed       string
	UserProfileUserNotFound string
}

var UserProfileMessage = UserProfileResponseMessage{
	UserProfileSuccess:      "User profile retrieved successfully",
	UserProfileFailed:       "Failed to retrieve user profile",
	UserProfileUserNotFound: "User not found",
}

type UpdateProfileResponseMessage struct {
	UpdateProfileSuccess      string
	UpdateProfileFailed       string
	UpdateProfileUserNotFound string
}

var UpdateProfileMessage = UpdateProfileResponseMessage{
	UpdateProfileSuccess:      "Update profile successful",
	UpdateProfileFailed:       "Update profile failed",
	UpdateProfileUserNotFound: "User not found",
}

type ForgotPasswordResponseMessage struct {
	ForgotPasswordSuccess                  string
	ForgotPasswordFailed                   string
	ForgotPasswordUpdateTokenVersionFailed string
	ForgotPasswordUserNotFound             string
	ForgotPasswordEmailPasswordNotMatch    string
	VerifyResetPasswordSuccess             string
	VerifyResetPasswordFailed              string
	VerifyResetPasswordTokenNotValid       string
}

var ForgotPasswordMessage = ForgotPasswordResponseMessage{
	ForgotPasswordSuccess:            "Forgot password successful",
	ForgotPasswordFailed:             "Forgot password failed",
	ForgotPasswordUserNotFound:       "User not found",
	VerifyResetPasswordSuccess:       "Verify reset password successful",
	VerifyResetPasswordFailed:        "Verify reset password failed",
	VerifyResetPasswordTokenNotValid: "Token not valid",
}

const (
	FORGOT_PASSWORD_REDIS_KEY = "forgot-password"

	FORGOT_PASSWORD_EXPIRES = time.Minute * 10
)
