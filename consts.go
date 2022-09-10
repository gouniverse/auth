package auth

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

type AuthenticatedUserID struct{}

const (
	keyEndpoint = contextKey("endpoint")

	// PathApiLogin contains the path to api login endpoint
	PathApiLogin string = "api/login"

	// PathApiLogin contains the path to api login endpoint
	PathApiLoginCodeVerify string = "api/login-code-verify"

	// PathApiLogout contains the path to api logout endpoint
	PathApiLogout string = "api/logout"

	// PathApiRegister contains the path to api register endpoint
	PathApiRegister string = "api/register"

	// PathApiRestorePassword contains the path to api restore password endpoint
	PathApiRestorePassword string = "api/restore-password"

	// PathApiResetPassword contains the path to api reset password endpoint
	PathApiResetPassword string = "api/reset-password"

	// PathLogin contains the path to login page
	PathLogin string = "login"

	// PathLogin contains the path to login page
	PathLoginCodeVerify string = "login-code-verify"

	// PathLogout contains the path to logout page
	PathLogout string = "logout"

	// PathRegister contains the path to logout page
	PathRegister string = "register"

	// PathRestore contains the path to password restore page
	PathPasswordRestore string = "password-restore"

	// PathReset contains the path to password reset page
	PathPasswordReset string = "password-reset"

	// LoginCodeLength specified the length of the login code
	LoginCodeLength int = 8

	// LoginCodeGamma specifies the characters to be used for building the login code
	LoginCodeGamma string = "BCDFGHJKLMNPQRSTVXYZ"
)
