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

const (
	MAX_SIGN_UP_RESEND_LIMIT         = 3
	MAX_FORGOT_PASSWORD_RESEND_LIMIT = 3
)

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
	SignupSuccess                  string
	SignupFailed                   string
	SignupExists                   string
	SignupFailedToCreateUser       string
	SignupFailedToCreateProfile    string
	SignupFailedToCreateAuthToken  string
	SignupFailedToUpdateUser       string
	SignupFailedToUpdateProfile    string
	SignupFailedToUpdateAuthToken  string
	SignupFailedToGenerateToken    string
	SignupFailedToSendEmail        string
	SignupFailedEncryptEmail       string
	SignupFailedResendNotAvailable string
	SignupFailedResendLimit        string
}

var SignupMessage = SignupResponseMessage{
	SignupSuccess:                  "Sign up successful",
	SignupFailed:                   "Sign up failed",
	SignupExists:                   "Email/phone number/username already registered",
	SignupFailedToCreateUser:       "Failed to create user",
	SignupFailedToCreateProfile:    "Failed to create profile",
	SignupFailedToCreateAuthToken:  "Failed to create auth token",
	SignupFailedToUpdateUser:       "Failed to update user",
	SignupFailedToUpdateProfile:    "Failed to update profile",
	SignupFailedToUpdateAuthToken:  "Failed to update auth token",
	SignupFailedToGenerateToken:    "Failed to generate access token",
	SignupFailedToSendEmail:        "Failed to send Email",
	SignupFailedEncryptEmail:       "Failed to encrypt Email",
	SignupFailedResendNotAvailable: "Resend not available",
	SignupFailedResendLimit:        "Account verification resend limit reached",
}

const (
	SIGN_UP_LIMITER_REDIS_KEY = "sign-up-limiter"

	SIGN_UP_LIMITER_RESEND_AVAILABLE_TIME          = time.Minute * 1
	SIGN_UP_LIMITER_RESEND_ON_LIMIT_AVAILABLE_TIME = time.Hour * 24
	SIGN_UP_LIMITER_ON_LIMIT_EXPIRES               = time.Hour * 24
)

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
	ResetPasswordSuccess                   string
	ResetPasswordFailed                    string
	ResetPasswordFailedResendNotAvailable  string
	ResetPasswordFailedResendLimit         string
}

var ForgotPasswordMessage = ForgotPasswordResponseMessage{
	ForgotPasswordSuccess:                 "Forgot password successful",
	ForgotPasswordFailed:                  "Forgot password failed",
	ForgotPasswordUserNotFound:            "User not found",
	VerifyResetPasswordSuccess:            "Verify reset password successful",
	VerifyResetPasswordFailed:             "Verify reset password failed",
	VerifyResetPasswordTokenNotValid:      "Token not valid",
	ResetPasswordSuccess:                  "Reset password successful",
	ResetPasswordFailed:                   "Reset password failed",
	ResetPasswordFailedResendNotAvailable: "Resend not available",
	ResetPasswordFailedResendLimit:        "Reset password resend limit reached",
}

const (
	FORGOT_PASSWORD_REDIS_KEY = "forgot-password"

	FORGOT_PASSWORD_RESEND_AVAILABLE_TIME          = time.Minute * 1
	FORGOT_PASSWORD_EXPIRES                        = time.Minute * 5
	FORGOT_PASSWORD_RESEND_ON_LIMIT_AVAILABLE_TIME = time.Hour * 24
	FORGOT_PASSWORD_ON_LIMIT_EXPIRES               = time.Hour * 24
)
