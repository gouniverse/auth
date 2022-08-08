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

	token := utils.RandStr(32)

	errSession := a.funcUserStoreToken(token, userID)

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
