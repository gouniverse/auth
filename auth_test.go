package auth

import (
	"testing"

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

func TestInitializationTestSuite(t *testing.T) {
	suite.Run(t, new(initTestSuite))
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(apiTestSuite))
}

func TestUiTestSuite(t *testing.T) {
	suite.Run(t, new(uiTestSuite))
}
