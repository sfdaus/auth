package constant

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
