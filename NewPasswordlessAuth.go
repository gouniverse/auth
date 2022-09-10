package auth

import "errors"

type PasswordlessConfig struct {
	EnableRegistration bool

	Endpoint string

	// shared options
	FuncTemporaryKeyGet func(key string) (value string, err error)
	FuncTemporaryKeySet func(key string, value string, expiresSeconds int) (err error)

	FuncUserFindByEmail        func(email string) (userID string, err error)
	FuncEmailTemplateLoginCode func(email string, passwordRestoreLink string) string // optional
	FuncEmailSend              func(email string, emailSubject string, emailBody string) (err error)
	FuncUserRegister           func(email string, firstName string, lastName string) (err error)

	UrlRedirectOnSuccess string

	UseCookies      bool
	UseLocalStorage bool
}

func NewPasswordlessAuth(config PasswordlessConfig) (*Auth, error) {
	auth := &Auth{}

	if config.Endpoint == "" {
		return nil, errors.New("auth: endpoint is required")
	}

	if config.UrlRedirectOnSuccess == "" {
		return nil, errors.New("auth: url to redirect to on success is required")
	}

	if config.FuncTemporaryKeyGet == nil {
		return nil, errors.New("auth: FuncTemporaryKeyGet function is required")
	}

	if config.FuncTemporaryKeySet == nil {
		return nil, errors.New("auth: FuncTemporaryKeySet function is required")
	}

	// if config.FuncUserFindByAuthToken == nil {
	// 	return nil, errors.New("auth: FuncUserFindByToken function is required")
	// }

	// if config.FuncUserFindByUsername == nil {
	// 	return nil, errors.New("auth: FuncUserFindByUsername function is required")
	// }

	// if config.FuncUserLogin == nil {
	// 	return nil, errors.New("auth: FuncUserLogin function is required")
	// }

	// if config.FuncUserLogout == nil {
	// 	return nil, errors.New("auth: FuncUserLogout function is required")
	// }

	// if config.EnableRegistration && config.FuncUserRegister == nil {
	// 	return nil, errors.New("auth: FuncUserRegister function is required")
	// }

	// if config.FuncUserStoreAuthToken == nil {
	// 	return nil, errors.New("auth: FuncUserStoreToken function is required")
	// }

	if config.FuncEmailSend == nil {
		return nil, errors.New("auth: FuncEmailSend function is required")
	}

	//if config.FuncUserRegister != nil {
	//	config.EnableRegistration = true
	//}

	if config.FuncLayout == nil {
		config.FuncLayout = auth.layout
	}

	auth.enableRegistration = config.EnableRegistration
	auth.endpoint = config.Endpoint
	auth.passwordless = true
	auth.urlRedirectOnSuccess = config.UrlRedirectOnSuccess
	// auth.useCookies = config.UseCookies
	// auth.useLocalStorage = config.UseLocalStorage
	// auth.funcEmailSend = config.FuncEmailSend
	// auth.funcEmailTemplatePasswordRestore = config.FuncEmailTemplatePasswordRestore
	// auth.funcLayout = config.FuncLayout
	// auth.funcTemporaryKeyGet = config.FuncTemporaryKeyGet
	// auth.funcTemporaryKeySet = config.FuncTemporaryKeySet
	// auth.funcUserLogin = config.FuncUserLogin
	// auth.funcUserLogout = config.FuncUserLogout
	// auth.funcUserPasswordChange = config.FuncUserPasswordChange
	// auth.funcUserRegister = config.FuncUserRegister
	// auth.funcUserFindByAuthToken = config.FuncUserFindByAuthToken
	// auth.funcUserFindByUsername = config.FuncUserFindByUsername
	// auth.funcUserStoreAuthToken = config.FuncUserStoreAuthToken

	// If no user defined email template is set, use default
	if auth.funcEmailTemplatePasswordRestore == nil {
		auth.funcEmailTemplatePasswordRestore = emailTemplatePasswordChange
	}

	return auth, nil
}
