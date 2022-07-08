package auth

// Config defines the available configuration options for authentication
type Config struct {
	Endpoint             string
	EnableRegistration   bool
	UrlRedirectOnSuccess string
	FuncUserFindByToken  func(sessionID string) (userID string, err error)
	FuncUserLogin        func(username string, password string) (userID string, err error)
	FuncUserRegister     func(username string, password string, first_name string, last_name string) (err error)
	FuncUserStoreToken   func(sessionID string, userID string) error
	FuncUserLogout       func(username string) (err error)
	UseCookies           bool
	UseLocalStorage      bool
}
