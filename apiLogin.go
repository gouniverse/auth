package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
	validator "github.com/gouniverse/validator"
)

func (a Auth) apiLogin(w http.ResponseWriter, r *http.Request) {
	if a.passwordless {
		a.apiLoginPasswordless(w, r)
	} else {
		a.apiLoginUsernameAndPassword(w, r)
	}
}

func (a Auth) apiLoginPasswordless(w http.ResponseWriter, r *http.Request) {
	email := strings.Trim(utils.Req(r, "email", ""), " ")

	if email == "" {
		api.Respond(w, r, api.Error("Email is required field"))
		return
	}

	if !validator.IsEmail(email) {
		api.Respond(w, r, api.Error("This is not a valid email: "+email))
		return
	}

	verificationCode := utils.StrRandomFromGamma(LoginCodeLength, LoginCodeGamma)

	errTempTokenSave := a.funcTemporaryKeySet(verificationCode, email, 3600)

	if errTempTokenSave != nil {
		api.Respond(w, r, api.Error("token store failed. "+errTempTokenSave.Error()))
		return
	}

	emailContent := a.passwordlessFuncEmailTemplateLoginCode(email, verificationCode, UserAuthOptions{
		UserIp:    utils.IP(r),
		UserAgent: r.UserAgent(),
	})

	errEmailSent := a.passwordlessFuncEmailSend(email, "Login Code", emailContent, UserAuthOptions{
		UserIp:    utils.IP(r),
		UserAgent: r.UserAgent(),
	})

	if errEmailSent != nil {
		log.Println(errEmailSent)
		api.Respond(w, r, api.Error("Login code failed to be send. Please try again later"))
		return
	}

	api.Respond(w, r, api.Success("Login code was sent successfully"))
}

func (a Auth) apiLoginUsernameAndPassword(w http.ResponseWriter, r *http.Request) {
	email := strings.Trim(utils.Req(r, "email", ""), " ")
	password := strings.Trim(utils.Req(r, "password", ""), " ")

	response := a.LoginWithUsernameAndPassword(email, password, utils.IP(r), r.UserAgent())

	if response.ErrorMessage != "" {
		api.Respond(w, r, api.Error(response.ErrorMessage))
		return
	}

	if a.useCookies {
		AuthCookieSet(w, r, response.Token)
	}

	api.Respond(w, r, api.SuccessWithData(response.SuccessMessage, map[string]any{
		"token": response.Token,
	}))
}
