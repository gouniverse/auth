package auth

import (
	"net/http"
	"strings"
)

// Auth defines the structure for the authentication
type Auth struct {
	endpoint           string
	enableRegistration bool

	// urlRedirectOnSuccess the endpoint to return to on success
	urlRedirectOnSuccess string

	funcEmailTemplatePasswordRestore func(userID string, passwordRestoreLink string) string // optional
	funcEmailSend                    func(userID string, emailSubject string, emailBody string) (err error)
	funcStoreTemporaryKey            func(key string, value string, expiresSeconds int) (err error)
	funcUserLogin                    func(username string, password string) (userID string, err error)
	funcUserLogout                   func(username string) (err error)
	funcUserStoreToken               func(token string, userID string) error
	funcUserFindByToken              func(token string) (userID string, err error)
	funcUserRegister                 func(username string, password string, first_name string, last_name string) (err error)
	funcUserFindByUsername           func(username string, first_name string, last_name string) (userID string, err error)
	useCookies                       bool
	useLocalStorage                  bool
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

func (a Auth) LinkApiLogout() string {
	return link(a.endpoint, pathApiLogout)
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

func (a Auth) LinkLogout() string {
	return link(a.endpoint, pathLogout)
}

func (a Auth) LinkPasswordRestore() string {
	return link(a.endpoint, pathPasswordRestore)
}

func (a Auth) LinkPasswordReset(token string) string {
	return link(a.endpoint, pathPasswordReset) + "?t=" + token
}

func (a Auth) LinkRegister() string {
	return link(a.endpoint, pathRegister)
}

func (a Auth) LinkRedirectOnSuccess() string {
	return a.urlRedirectOnSuccess
}

func link(endpoint, uri string) string {
	if strings.HasSuffix(endpoint, "/") {
		return endpoint + uri
	} else {
		return endpoint + "/" + uri
	}
}
