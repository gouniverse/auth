package auth

// Config defines the available configuration options for authentication
type Config struct {
	Endpoint                         string
	EnableRegistration               bool
	UrlRedirectOnSuccess             string
	FuncEmailTemplatePasswordRestore func(userID string, passwordRestoreLink string) string // optional
	FuncEmailSend                    func(userID string, emailSubject string, emailBody string) (err error)
	FuncTemporaryKeyGet              func(key string) (value string, err error)
	FuncTemporaryKeySet              func(key string, value string, expiresSeconds int) (err error)
	FuncUserFindByToken              func(sessionID string) (userID string, err error)
	FuncUserFindByUsername           func(username string, firstName string, lastName string) (userID string, err error)
	FuncUserLogin                    func(username string, password string) (userID string, err error)
	FuncUserPasswordChange           func(username string, newPassword string) (err error)
	FuncUserRegister                 func(username string, password string, first_name string, last_name string) (err error)
	FuncUserStoreToken               func(sessionID string, userID string) error
	FuncUserLogout                   func(username string) (err error)
	UseCookies                       bool
	UseLocalStorage                  bool
}
