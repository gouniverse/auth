package auth

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

type AuthenticatedUserID struct{}

const (
	keyEndpoint = contextKey("endpoint")

	// PathApiLogin contains the path to api login endpoint
	pathApiLogin string = "api/login"

	// PathApiRegister contains the path to api register endpoint
	pathApiRegister string = "api/register"

	// PathApiRestorePassword contains the path to api restore password endpoint
	pathApiRestorePassword string = "api/restore-password"

	// PathApiResetPassword contains the path to api reset password endpoint
	pathApiResetPassword string = "api/reset-password"

	// PathLogin contains the path to login page
	pathLogin string = "login"

	// PathLogout contains the path to logout page
	pathLogout string = "logout"

	// PathRegister contains the path to logout page
	pathRegister string = "register"

	// PathRestore contains the path to password restore page
	pathPasswordRestore string = "password-restore"

	// PathReset contains the path to password reset page
	pathPasswordReset string = "password-reset"
)
