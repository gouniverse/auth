package auth

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type initTestSuite struct {
	suite.Suite
}

func (suite *initTestSuite) SetupTest() {
	// config.SetupTests()
}

func (suite *initTestSuite) TestEndpointIsRequired() {
	// expected := `<title>Home | Rem.land</title>`
	// assert.HTTPBodyContainsf(suite.T(), routes.Router().ServeHTTP, "POST", links.Home(), url.Values{}, expected, "%")
	_, err := NewAuth(Config{})
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "auth: endpoint is required", err.Error())
}

func (suite *initTestSuite) TesUrlRedirectOnSuccessIsRequired() {
	// expected := `<title>Home | Rem.land</title>`
	// assert.HTTPBodyContainsf(suite.T(), routes.Router().ServeHTTP, "POST", links.Home(), url.Values{}, expected, "%")
	_, err := NewAuth(Config{
		Endpoint: "http://localhost/auth",
	})
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "auth: url to redirect to on success is required", err.Error())
}

func (suite *initTestSuite) TestFuncTemporaryKeyGetIsRequired() {
	// expected := `<title>Home | Rem.land</title>`
	// assert.HTTPBodyContainsf(suite.T(), routes.Router().ServeHTTP, "POST", links.Home(), url.Values{}, expected, "%")
	_, err := NewAuth(Config{
		Endpoint:             "http://localhost/auth",
		UrlRedirectOnSuccess: "http://localhost/dashboard",
	})
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "auth: FuncTemporaryKeyGet function is required", err.Error())
}

func (suite *initTestSuite) TestFuncTemporaryKeySetIsRequired() {
	// expected := `<title>Home | Rem.land</title>`
	// assert.HTTPBodyContainsf(suite.T(), routes.Router().ServeHTTP, "POST", links.Home(), url.Values{}, expected, "%")
	_, err := NewAuth(Config{
		Endpoint:             "http://localhost/auth",
		UrlRedirectOnSuccess: "http://localhost/dashboard",
		FuncTemporaryKeyGet:  func(key string) (value string, err error) { return "", nil },
	})
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "auth: FuncTemporaryKeySet function is required", err.Error())
}

func (suite *initTestSuite) TestFuncUserFindByTokenIsRequired() {
	_, err := NewAuth(Config{
		Endpoint:             "http://localhost/auth",
		UrlRedirectOnSuccess: "http://localhost/dashboard",
		FuncTemporaryKeyGet:  func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:  func(key string, value string, expiresSeconds int) (err error) { return nil },
	})
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "auth: FuncUserFindByToken function is required", err.Error())
}

func (suite *initTestSuite) TestFuncUserFindByUsernameIsRequired() {
	_, err := NewAuth(Config{
		Endpoint:             "http://localhost/auth",
		UrlRedirectOnSuccess: "http://localhost/dashboard",
		FuncTemporaryKeyGet:  func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:  func(key string, value string, expiresSeconds int) (err error) { return nil },
		FuncUserFindByToken:  func(token string) (userID string, err error) { return "", nil },
	})
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "auth: FuncUserFindByUsername function is required", err.Error())
}

func (suite *initTestSuite) TestInitializationSuccess() {
	auth, err := NewAuth(Config{
		Endpoint:               "http://localhost/auth",
		UrlRedirectOnSuccess:   "http://localhost/dashboard",
		FuncTemporaryKeyGet:    func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:    func(key string, value string, expiresSeconds int) (err error) { return nil },
		FuncUserFindByToken:    func(token string) (userID string, err error) { return "", nil },
		FuncUserFindByUsername: func(username string, firstName string, lastName string) (userID string, err error) { return "", nil },
		FuncUserLogin:          func(username string, password string) (userID string, err error) { return "", nil },
		FuncUserLogout:         func(userID string) (err error) { return nil },
		FuncUserStoreToken:     func(sessionID string, userID string) (err error) { return nil },
		FuncEmailSend:          func(userID string, emailSubject string, emailBody string) (err error) { return nil },
	})
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), auth)
}

func (suite *initTestSuite) TestLinks() {
	endpoint := "http://localhost/auth"
	auth, err := NewAuth(Config{
		Endpoint:               endpoint,
		UrlRedirectOnSuccess:   "http://localhost/dashboard",
		FuncTemporaryKeyGet:    func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:    func(key string, value string, expiresSeconds int) (err error) { return nil },
		FuncUserFindByToken:    func(token string) (userID string, err error) { return "", nil },
		FuncUserFindByUsername: func(username string, firstName string, lastName string) (userID string, err error) { return "", nil },
		FuncUserLogin:          func(username string, password string) (userID string, err error) { return "", nil },
		FuncUserLogout:         func(userID string) (err error) { return nil },
		FuncUserStoreToken:     func(sessionID string, userID string) (err error) { return nil },
		FuncEmailSend:          func(userID string, emailSubject string, emailBody string) (err error) { return nil },
	})
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), auth)
	assert.Equal(suite.T(), endpoint+"/"+pathApiLogin, auth.LinkApiLogin())
	assert.Equal(suite.T(), endpoint+"/"+pathApiLogout, auth.LinkApiLogout())
	assert.Equal(suite.T(), endpoint+"/"+pathApiRegister, auth.LinkApiRegister())
	assert.Equal(suite.T(), endpoint+"/"+pathApiResetPassword, auth.LinkApiPasswordReset())
	assert.Equal(suite.T(), endpoint+"/"+pathApiRestorePassword, auth.LinkApiPasswordRestore())

	assert.Equal(suite.T(), endpoint+"/"+pathLogin, auth.LinkLogin())
	assert.Equal(suite.T(), endpoint+"/"+pathLogout, auth.LinkLogout())
	assert.Equal(suite.T(), endpoint+"/"+pathPasswordReset+"?t=mytoken", auth.LinkPasswordReset("mytoken"))
	assert.Equal(suite.T(), endpoint+"/"+pathPasswordRestore, auth.LinkPasswordRestore())
	assert.Equal(suite.T(), endpoint+"/"+pathRegister, auth.LinkRegister())
}
