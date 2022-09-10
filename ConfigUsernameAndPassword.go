package auth

// Config defines the available configuration options for authentication
type ConfigUsernameAndPassword struct {
	// ===== START: shared by all implementations
	EnableRegistration      bool
	Endpoint                string
	FuncLayout              func(content string) string
	FuncTemporaryKeyGet     func(key string) (value string, err error)
	FuncTemporaryKeySet     func(key string, value string, expiresSeconds int) (err error)
	FuncUserStoreAuthToken  func(sessionID string, userID string) error
	FuncUserFindByAuthToken func(sessionID string) (userID string, err error)
	UrlRedirectOnSuccess    string
	UseCookies              bool
	UseLocalStorage         bool
	// ===== END: shared by all implementations

	// ===== START: username(email) and password options
	FuncEmailTemplatePasswordRestore func(userID string, passwordRestoreLink string) string // optional
	FuncEmailSend                    func(userID string, emailSubject string, emailBody string) (err error)
	FuncUserFindByUsername           func(username string, firstName string, lastName string) (userID string, err error)
	FuncUserLogin                    func(username string, password string) (userID string, err error)
	FuncUserLogout                   func(username string) (err error)
	FuncUserPasswordChange           func(username string, newPassword string) (err error)
	FuncUserRegister                 func(username string, password string, first_name string, last_name string) (err error)
	LabelUsername                    string
	// ===== END: username(email) and password options
}
