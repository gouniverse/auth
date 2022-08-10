package auth

import (
	"net/url"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type apiTestSuite struct {
	suite.Suite
}

func (suite *apiTestSuite) SetupTest() {
	// config.SetupTests()
}

func (suite *apiTestSuite) TestLoginEndpointRequiresEmail() {
	endpoint := "http://localhost"
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

	expectedError := `"status":"error"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiLogin(), url.Values{}, expectedError, "%")

	expectedErrorMessage := `"message":"Email is required field"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiLogin(), url.Values{}, expectedErrorMessage, "%")
}

func (suite *apiTestSuite) TestLoginEndpointRequiresPassword() {
	endpoint := "http://localhost"
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

	expectedError := `"status":"error"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiLogin(), url.Values{
		"email": {"test@test.com"},
	}, expectedError, "%")

	expectedErrorMessage := `"message":"Password is required field"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiLogin(), url.Values{
		"email": {"test@test.com"},
	}, expectedErrorMessage, "%")
}

func (suite *apiTestSuite) TestLoginEndpointRequiresPasswords() {
	endpoint := "http://localhost"
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

	expectedSuccess := `"status":"success"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiLogin(), url.Values{
		"email":    {"test@test.com"},
		"password": {"1234"},
	}, expectedSuccess, "%")

	expectedSuccessMessage := `"message":"login success"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiLogin(), url.Values{
		"email":    {"test@test.com"},
		"password": {"1234"},
	}, expectedSuccessMessage, "%")

	expectedToken := `"token":"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiLogin(), url.Values{
		"email":    {"test@test.com"},
		"password": {"1234"},
	}, expectedToken, "%")
}

func (suite *apiTestSuite) TestRegisterEndpointRequiresFirstName() {
	endpoint := "http://localhost"
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

	expectedError := `"status":"error"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiRegister(), url.Values{}, expectedError, "%")

	expectedErrorMessage := `"message":"First name is required field"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiRegister(), url.Values{}, expectedErrorMessage, "%")
}

func (suite *apiTestSuite) TestRegisterEndpointRequiresLastName() {
	endpoint := "http://localhost"
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

	expectedError := `"status":"error"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiRegister(), url.Values{
		"first_name": {"John"},
	}, expectedError, "%")

	expectedErrorMessage := `"message":"Last name is required field"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiRegister(), url.Values{
		"first_name": {"John"},
	}, expectedErrorMessage, "%")
}

func (suite *apiTestSuite) TestRegisterEndpointRequiresEmail() {
	endpoint := "http://localhost"
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

	expectedError := `"status":"error"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiRegister(), url.Values{
		"first_name": {"John"},
		"last_name":  {"Doe"},
	}, expectedError, "%")

	expectedErrorMessage := `"message":"Email is required field"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiRegister(), url.Values{
		"first_name": {"John"},
		"last_name":  {"Doe"},
	}, expectedErrorMessage, "%")
}

func (suite *apiTestSuite) TestRegisterEndpointRequiresPassword() {
	endpoint := "http://localhost"
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

	expectedError := `"status":"error"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiRegister(), url.Values{
		"first_name": {"John"},
		"last_name":  {"Doe"},
		"email":      {"test@test.com"},
	}, expectedError, "%")

	expectedErrorMessage := `"message":"Password is required field"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiRegister(), url.Values{
		"first_name": {"John"},
		"last_name":  {"Doe"},
		"email":      {"test@test.com"},
	}, expectedErrorMessage, "%")
}

func (suite *apiTestSuite) TestRegisterEndpointRequiresPasswords() {
	endpoint := "http://localhost"
	auth, err := NewAuth(Config{
		Endpoint:                endpoint,
		UrlRedirectOnSuccess:    "http://localhost/dashboard",
		FuncTemporaryKeyGet:     func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:     func(key string, value string, expiresSeconds int) (err error) { return nil },
		FuncUserFindByAuthToken: func(token string) (userID string, err error) { return "", nil },
		FuncUserFindByUsername:  func(username string, firstName string, lastName string) (userID string, err error) { return "", nil },
		FuncUserLogin:           func(username string, password string) (userID string, err error) { return "", nil },
		FuncUserRegister:        func(username string, password string, firstName string, lastName string) (err error) { return nil },
		FuncUserLogout:          func(userID string) (err error) { return nil },
		FuncUserStoreAuthToken:  func(sessionID string, userID string) (err error) { return nil },
		FuncEmailSend:           func(userID string, emailSubject string, emailBody string) (err error) { return nil },
		EnableRegistration:      true,
	})

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), auth)

	expectedSuccess := `"status":"success"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiRegister(), url.Values{
		"first_name": {"John"},
		"last_name":  {"Doe"},
		"email":      {"test@test.com"},
		"password":   {"1234"},
	}, expectedSuccess, "%")

	expectedMessage := `"message":"registration success"`
	assert.HTTPBodyContainsf(suite.T(), auth.Router().ServeHTTP, "POST", auth.LinkApiRegister(), url.Values{
		"first_name": {"John"},
		"last_name":  {"Doe"},
		"email":      {"test@test.com"},
		"password":   {"1234"},
	}, expectedMessage, "%")
}
