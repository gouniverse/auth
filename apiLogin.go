package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
	validator "github.com/gouniverse/validator"
)

func (a Auth) apiLogin(w http.ResponseWriter, r *http.Request) {
	if a.Passwordless {
		a.apiLoginPasswordless(w, r)
	} else {
		a.apiLoginEmailAndPassword(w, r)
	}
}

func (a Auth) apiLoginEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	email := strings.Trim(utils.Req(r, "email", ""), " ")
	password := strings.Trim(utils.Req(r, "password", ""), " ")

	if email == "" {
		api.Respond(w, r, api.Error("Email is required field"))
		return
	}

	if password == "" {
		api.Respond(w, r, api.Error("Password is required field"))
		return
	}

	if !validator.IsEmail(email) {
		api.Respond(w, r, api.Error("This is not a valid email: "+email))
		return
	}

	userID, err := a.funcUserLogin(email, password)

	if err != nil {
		api.Respond(w, r, api.Error("authentication failed. "+err.Error()))
		return
	}

	token := utils.StrRandom(32)

	errSession := a.funcUserStoreAuthToken(token, userID)

	if errSession != nil {
		api.Respond(w, r, api.Error("token store failed. "+errSession.Error()))
		return
	}

	if a.useCookies {
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:     "authtoken",
			Value:    token,
			Expires:  expiration,
			HttpOnly: false,
			Secure:   true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)
	}

	api.Respond(w, r, api.SuccessWithData("login success", map[string]interface{}{
		"token": token,
	}))
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

	mailSent := true

	if mailSent {
		api.Respond(w, r, api.Success("Your login code has been sent"))
		return
	}

	api.Respond(w, r, api.Error("Sorry login code failed to be emailed"))

	// userID, err := a.funcUserFindByUsername(email)

	// if err != nil {
	// 	api.Respond(w, r, api.Error("authentication failed. "+err.Error()))
	// 	return
	// }

	// token := utils.StrRandom(32)

	// errSession := a.funcUserStoreAuthToken(token, userID)

	// if errSession != nil {
	// 	api.Respond(w, r, api.Error("token store failed. "+errSession.Error()))
	// 	return
	// }

	// if a.useCookies {
	// 	expiration := time.Now().Add(365 * 24 * time.Hour)
	// 	cookie := http.Cookie{
	// 		Name:     "authtoken",
	// 		Value:    token,
	// 		Expires:  expiration,
	// 		HttpOnly: false,
	// 		Secure:   true,
	// 		Path:     "/",
	// 	}
	// 	http.SetCookie(w, &cookie)
	// }

	// api.Respond(w, r, api.SuccessWithData("login success", map[string]interface{}{
	// 	"token": token,
	// }))
}
