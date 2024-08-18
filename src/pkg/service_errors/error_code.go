package service_errors

const (
	// TOKEN
	UnExpectedError = "Expected error"
	ClaimNotFound = "Claim not found"
	TokenRequired = "token required"
	TokenExpired = "token expired"
	TokenInvalid = "token invalid"

	// OTP
	OtpExists = "Otp exists"
	OtpUsed   = "Otp used"
	OtpInvalid   = "Otp invalid"

	// User
	EmailExists = "Email exists"
	UsernameExists = "Username Exists"
)
