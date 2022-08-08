package auth

import "errors"

func NewAuth(config Config) (*Auth, error) {
	auth := &Auth{}

	if config.Endpoint == "" {
		return nil, errors.New("auth: endpoint is required")
	}

	if config.UrlRedirectOnSuccess == "" {
		return nil, errors.New("auth: url to redirect to on success is required")
	}

	if config.FuncStoreTemporaryKey == nil {
		return nil, errors.New("auth: FuncStoreTemporaryKey function is required")
	}

	if config.FuncUserLogin == nil {
		return nil, errors.New("auth: FuncUserLogin function is required")
	}

	if config.FuncUserStoreToken == nil {
		return nil, errors.New("auth: FuncUserStoreToken function is required")
	}

	if config.FuncUserFindByToken == nil {
		return nil, errors.New("auth: FuncUserFindByToken function is required")
	}

	if config.FuncUserFindByUsername == nil {
		return nil, errors.New("auth: FuncUserFindByUsername function is required")
	}

	if config.FuncUserLogout == nil {
		return nil, errors.New("auth: FuncUserLogout function is required")
	}

	if config.FuncEmailSend == nil {
		return nil, errors.New("auth: FuncEmailSend function is required")
	}

	if config.FuncUserRegister != nil {
		config.EnableRegistration = true
		//return nil, errors.New("auth: FuncUserRegister function is required")
	}

	auth.enableRegistration = config.EnableRegistration
	auth.endpoint = config.Endpoint
	auth.urlRedirectOnSuccess = config.UrlRedirectOnSuccess
	auth.useCookies = config.UseCookies
	auth.useLocalStorage = config.UseLocalStorage
	auth.funcEmailSend = config.FuncEmailSend
	auth.funcEmailTemplatePasswordRestore = config.FuncEmailTemplatePasswordRestore
	auth.funcUserLogin = config.FuncUserLogin
	auth.funcUserLogout = config.FuncUserLogout
	auth.funcUserRegister = config.FuncUserRegister
	auth.funcUserFindByToken = config.FuncUserFindByToken
	auth.funcUserFindByUsername = config.FuncUserFindByUsername
	auth.funcUserStoreToken = config.FuncUserStoreToken
	auth.funcStoreTemporaryKey = config.FuncStoreTemporaryKey

	if auth.funcEmailTemplatePasswordRestore == nil {
		auth.funcEmailTemplatePasswordRestore = emailTemplatePasswordChange
	}

	return auth, nil
}
