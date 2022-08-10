package auth

import (
	"net/url"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type uiTestSuite struct {
	suite.Suite
}

func (suite *uiTestSuite) SetupTest() {
	// config.SetupTests()
}

func (suite *uiTestSuite) TestPageLogin() {
	endpoint := "http://localhost/auth"
	auth, err := NewAuth(Config{
		Endpoint:                endpoint,
		UrlRedirectOnSuccess:    "http://localhost/dashboard",
		FuncTemporaryKeyGet:     func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:     func(key string, value string, expiresSeconds int) (err error) { return nil },
		FuncUserFindByAuthToken: func(token string) (userID string, err error) { return "", nil },
		FuncUserFindByUsername:  func(username string, firstName string, lastName string) (userID string, err error) { return "", nil },
		FuncUserLogin:           func(username string, password string) (userID string, err error) { return "", nil },
		FuncUserLogout:          func(userID string) (err error) { return nil },
		FuncUserStoreAuthToken:  func(sessionID string, userID string) (err error) { return nil },
		FuncEmailSend:           func(userID string, emailSubject string, emailBody string) (err error) { return nil },
	})

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), auth)
	expected := `<title>Login</title>`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkLogin(), url.Values{}, expected, "%")
}

func (suite *uiTestSuite) TestPageRegister() {
	endpoint := "http://localhost/auth"
	auth, err := NewAuth(Config{
		Endpoint:                endpoint,
		UrlRedirectOnSuccess:    "http://localhost/dashboard",
		FuncTemporaryKeyGet:     func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:     func(key string, value string, expiresSeconds int) (err error) { return nil },
		FuncUserFindByAuthToken: func(token string) (userID string, err error) { return "", nil },
		FuncUserFindByUsername:  func(username string, firstName string, lastName string) (userID string, err error) { return "", nil },
		FuncUserLogin:           func(username string, password string) (userID string, err error) { return "", nil },
		FuncUserLogout:          func(userID string) (err error) { return nil },
		FuncUserRegister:        func(username string, password string, firstName string, lastName string) (err error) { return nil },
		FuncUserStoreAuthToken:  func(sessionID string, userID string) (err error) { return nil },
		FuncEmailSend:           func(userID string, emailSubject string, emailBody string) (err error) { return nil },
		EnableRegistration:      true,
	})

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), auth)
	expected := `<title>Register</title>`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkRegister(), url.Values{}, expected, "%")
}

func (suite *uiTestSuite) TestPageRegisterDisabled() {
	endpoint := "http://localhost/auth"
	auth, err := NewAuth(Config{
		Endpoint:                endpoint,
		UrlRedirectOnSuccess:    "http://localhost/dashboard",
		FuncTemporaryKeyGet:     func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:     func(key string, value string, expiresSeconds int) (err error) { return nil },
		FuncUserFindByAuthToken: func(token string) (userID string, err error) { return "", nil },
		FuncUserFindByUsername:  func(username string, firstName string, lastName string) (userID string, err error) { return "", nil },
		FuncUserLogin:           func(username string, password string) (userID string, err error) { return "", nil },
		FuncUserLogout:          func(userID string) (err error) { return nil },
		FuncUserStoreAuthToken:  func(sessionID string, userID string) (err error) { return nil },
		FuncEmailSend:           func(userID string, emailSubject string, emailBody string) (err error) { return nil },
		EnableRegistration:      false,
	})

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), auth)
	expectedTitle := `<title>Register</title>`
	assert.HTTPBodyNotContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkRegister(), url.Values{}, expectedTitle, "%")
	expectedEmpty := ""
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkRegister(), url.Values{}, expectedEmpty, "%")
}

func (suite *uiTestSuite) TestPagePasswordRestore() {
	endpoint := "http://localhost/auth"
	auth, err := NewAuth(Config{
		Endpoint:                endpoint,
		UrlRedirectOnSuccess:    "http://localhost/dashboard",
		FuncTemporaryKeyGet:     func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:     func(key string, value string, expiresSeconds int) (err error) { return nil },
		FuncUserFindByAuthToken: func(token string) (userID string, err error) { return "", nil },
		FuncUserFindByUsername:  func(username string, firstName string, lastName string) (userID string, err error) { return "", nil },
		FuncUserLogin:           func(username string, password string) (userID string, err error) { return "", nil },
		FuncUserLogout:          func(userID string) (err error) { return nil },
		FuncUserRegister:        func(username string, password string, firstName string, lastName string) (err error) { return nil },
		FuncUserStoreAuthToken:  func(sessionID string, userID string) (err error) { return nil },
		FuncEmailSend:           func(userID string, emailSubject string, emailBody string) (err error) { return nil },
		EnableRegistration:      true,
	})

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), auth)
	expected := `<title>Restore Forgotten Password</title>`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkPasswordRestore(), url.Values{}, expected, "%")
}

func (suite *uiTestSuite) TestPagePasswordReset() {
	endpoint := "http://localhost/auth"
	auth, err := NewAuth(Config{
		Endpoint:                endpoint,
		UrlRedirectOnSuccess:    "http://localhost/dashboard",
		FuncTemporaryKeyGet:     func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:     func(key string, value string, expiresSeconds int) (err error) { return nil },
		FuncUserFindByAuthToken: func(token string) (userID string, err error) { return "", nil },
		FuncUserFindByUsername:  func(username string, firstName string, lastName string) (userID string, err error) { return "", nil },
		FuncUserLogin:           func(username string, password string) (userID string, err error) { return "", nil },
		FuncUserLogout:          func(userID string) (err error) { return nil },
		FuncUserRegister:        func(username string, password string, firstName string, lastName string) (err error) { return nil },
		FuncUserStoreAuthToken:  func(sessionID string, userID string) (err error) { return nil },
		FuncEmailSend:           func(userID string, emailSubject string, emailBody string) (err error) { return nil },
		EnableRegistration:      true,
	})

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), auth)
	expected := `<title>Reset Password</title>`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkPasswordReset("testtoken"), url.Values{}, expected, "%")
}
