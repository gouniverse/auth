package auth

import (
	"net/http"
	"strings"
)

// Auth defines the structure for the authentication
type Auth struct {
	endpoint string

	// enableRegistration enables the registration page and endpoint
	enableRegistration bool

	// urlRedirectOnSuccess the endpoint to return to on success
	urlRedirectOnSuccess string

	funcEmailTemplatePasswordRestore func(userID string, passwordRestoreLink string) string // optional
	funcEmailSend                    func(userID string, emailSubject string, emailBody string) (err error)

	funcLayout              func(content string) string
	funcTemporaryKeyGet     func(key string) (value string, err error)
	funcTemporaryKeySet     func(key string, value string, expiresSeconds int) (err error)
	funcUserLogin           func(username string, password string) (userID string, err error)
	funcUserLogout          func(username string) (err error)
	funcUserStoreAuthToken  func(token string, userID string) error
	funcUserFindByAuthToken func(token string) (userID string, err error)
	funcUserPasswordChange  func(username string, newPassword string) (err error)
	funcUserRegister        func(username string, password string, first_name string, last_name string) (err error)
	funcUserFindByUsername  func(username string, first_name string, last_name string) (userID string, err error)

	// labelUsername   string
	useCookies      bool
	useLocalStorage bool
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

// LinkPasswordReset - returns the password reset URL
func (a Auth) LinkPasswordReset(token string) string {
	return link(a.endpoint, pathPasswordReset) + "?t=" + token
}

// LinkRegister - returns the regsitration URL
func (a Auth) LinkRegister() string {
	return link(a.endpoint, pathRegister)
}

// LinkRedirectOnSuccess - returns the URL to where the user will be redirected after successful registration
func (a Auth) LinkRedirectOnSuccess() string {
	return a.urlRedirectOnSuccess
}

// link - creates the final URL by combining the provided endpoint with the provided URL
func link(endpoint, uri string) string {
	if strings.HasSuffix(endpoint, "/") {
		return endpoint + uri
	} else {
		return endpoint + "/" + uri
	}
}

// RegistrationEnable - enables registration
func (a Auth) RegistrationEnable() {
	a.enableRegistration = true
}


// RegistrationDisable - disables registration
func (a Auth) RegistrationDisable() {
	a.enableRegistration = false
}
