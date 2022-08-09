package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func init() {
	// mailServer := smtpmock.New(smtpmock.ConfigurationAttr{
	// 	LogToStdout:       false, // enable if you have errors sending emails
	// 	LogServerActivity: true,
	// 	PortNumber:        32435,
	// 	HostAddress:       "127.0.0.1",
	// })

	// if err := mailServer.Start(); err != nil {
	// 	fmt.Println(err)
	// }
}

func TestHomeAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

type AuthTestSuite struct {
	suite.Suite
	// funnelID string
}

func (suite *AuthTestSuite) SetupTest() {
	// config.SetupTests()
}

func (suite *AuthTestSuite) TestEndpointIsRequired() {
	// expected := `<title>Home | Rem.land</title>`
	// assert.HTTPBodyContainsf(suite.T(), routes.Router().ServeHTTP, "POST", links.Home(), url.Values{}, expected, "%")
	_, err := NewAuth(Config{})
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "auth: endpoint is required", err.Error())
}

func (suite *AuthTestSuite) TesUrlRedirectOnSuccessIsRequired() {
	// expected := `<title>Home | Rem.land</title>`
	// assert.HTTPBodyContainsf(suite.T(), routes.Router().ServeHTTP, "POST", links.Home(), url.Values{}, expected, "%")
	_, err := NewAuth(Config{
		Endpoint: "http://localhost/auth",
	})
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "auth: url to redirect to on success is required", err.Error())
}

func (suite *AuthTestSuite) TestFuncTemporaryKeyGetIsRequired() {
	// expected := `<title>Home | Rem.land</title>`
	// assert.HTTPBodyContainsf(suite.T(), routes.Router().ServeHTTP, "POST", links.Home(), url.Values{}, expected, "%")
	_, err := NewAuth(Config{
		Endpoint:             "http://localhost/auth",
		UrlRedirectOnSuccess: "http://localhost/dashboard",
	})
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "auth: FuncTemporaryKeyGet function is required", err.Error())
}

func (suite *AuthTestSuite) TestFuncTemporaryKeySetIsRequired() {
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
