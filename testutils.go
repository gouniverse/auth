package auth

// testSetupUsernameAndPasswordAuth creates a new Auth for testing
func testSetupUsernameAndPasswordAuth() (*Auth, error) {
	endpoint := "http://localhost/auth"
	return NewUsernameAndPasswordAuth(ConfigUsernameAndPassword{
		Endpoint:                endpoint,
		UrlRedirectOnSuccess:    "http://localhost/dashboard",
		FuncTemporaryKeyGet:     func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:     func(key string, value string, expiresSeconds int) (err error) { return nil },
		FuncUserFindByAuthToken: func(token string, options UserAuthOptions) (userID string, err error) { return "", nil },
		FuncUserFindByUsername: func(username string, firstName string, lastName string, options UserAuthOptions) (userID string, err error) {
			return "", nil
		},
		FuncUserLogin: func(username string, password string, options UserAuthOptions) (userID string, err error) {
			return "", nil
		},
		FuncUserLogout: func(userID string, options UserAuthOptions) (err error) { return nil },
		FuncUserStoreAuthToken: func(sessionID string, userID string, options UserAuthOptions) (err error) {
			return nil
		},
		FuncEmailSend: func(userID string, emailSubject string, emailBody string) (err error) { return nil },
		UseCookies:    true,
	})
}

// testSetupUsernameAndPasswordAuth creates a new Auth for testing
func testSetupPasswordlessAuth() (*Auth, error) {
	endpoint := "http://localhost/auth"
	return NewPasswordlessAuth(ConfigPasswordless{
		Endpoint:                endpoint,
		UrlRedirectOnSuccess:    "http://localhost/dashboard",
		FuncTemporaryKeyGet:     func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:     func(key string, value string, expiresSeconds int) (err error) { return nil },
		FuncUserFindByAuthToken: func(token string, options UserAuthOptions) (userID string, err error) { return "111", nil },
		FuncUserFindByEmail:     func(email string, options UserAuthOptions) (userID string, err error) { return "111", nil },
		FuncUserLogout:          func(userID string, options UserAuthOptions) (err error) { return nil },
		FuncUserStoreAuthToken:  func(sessionID string, userID string, options UserAuthOptions) (err error) { return nil },
		FuncEmailSend:           func(email string, emailSubject string, emailBody string) (err error) { return nil },
		UseCookies:              true,
	})
}
