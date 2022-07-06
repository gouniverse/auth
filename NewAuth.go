package auth

import "errors"

func NewAuth(config Config) (*Auth, error) {
	auth := &Auth{}

	if config.Endpoint == "" {
		return nil, errors.New("auth: endpoint is required")
	}

	if config.SuccessRedirectEndpoint == "" {
		return nil, errors.New("auth: success endpoint is required")
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

	if config.EnableRegistration && config.FuncUserRegister == nil {
		return nil, errors.New("auth: FuncUserRegister function is required")
	}

	auth.enableRegistration = config.EnableRegistration
	auth.endpoint = config.Endpoint
	auth.successRedirectEndpoint = config.SuccessRedirectEndpoint
	auth.useCookies = config.UseCookies
	auth.useLocalStorage = config.UseLocalStorage
	auth.funcUserLogin = config.FuncUserLogin
	auth.funcUserRegister = config.FuncUserRegister
	auth.funcUserFindByToken = config.FuncUserFindByToken
	auth.funcUserStoreToken = config.FuncUserStoreToken

	return auth, nil
}
