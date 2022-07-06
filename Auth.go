package auth

import (
	"net/http"
	"strings"
)

// Auth defines the structure for the authentication
type Auth struct {
	endpoint           string
	enableRegistration bool

	// successRedirectEndpoint the endpoint to return to on success
	successRedirectEndpoint string
	funcUserLogin           func(username string, password string) (userID string, err error)
	funcUserStoreToken      func(sessionID string, userID string) error
	funcUserFindByToken     func(sessionID string) (userID string, err error)
	funcUserRegister        func(username string, password string, first_name string, last_name string) (err error)
	useCookies              bool
	useLocalStorage         bool
}

func (a Auth) GetCurrentUserID(r *http.Request) string {
	authenticatedUserID := r.Context().Value(AuthenticatedUserID{})
	if authenticatedUserID == nil {
		return ""
	}
	return authenticatedUserID.(string)
}

func (a Auth) LinkApiLogin() string {
	return link(a.endpoint, pathApiLogin)
}

func (a Auth) LinkApiRegister() string {
	return link(a.endpoint, pathApiRegister)
}

func (a Auth) LinkApiPasswordRestore() string {
	return link(a.endpoint, pathApiRestorePassword)
}

func (a Auth) LinkApiPasswordReset() string {
	return link(a.endpoint, pathApiResetPassword)
}

func (a Auth) LinkLogin() string {
	return link(a.endpoint, pathLogin)
}

func (a Auth) LinkPasswordRestore() string {
	return link(a.endpoint, pathPasswordRestore)
}

func (a Auth) LinkPasswordReset() string {
	return link(a.endpoint, pathPasswordReset)
}

func (a Auth) LinkRegister() string {
	return link(a.endpoint, pathRegister)
}

func link(endpoint, uri string) string {
	if strings.HasSuffix(endpoint, "/") {
		return endpoint + uri
	} else {
		return endpoint + "/" + uri
	}
}
