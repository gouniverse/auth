package auth

import (
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (a Auth) apiPaswordReset(w http.ResponseWriter, r *http.Request) {
	token := strings.Trim(utils.Req(r, "token", ""), " ")
	password := strings.Trim(utils.Req(r, "password", ""), " ")
	passwordConfirm := strings.Trim(utils.Req(r, "password_confirm", ""), " ")

	if token == "" {
		api.Respond(w, r, api.Error("Token is required field"))
		return
	}

	if password == "" {
		api.Respond(w, r, api.Error("Password is required field"))
		return
	}

	if password != passwordConfirm {
		api.Respond(w, r, api.Error("Passwords do not match"))
		return
	}

	userID, errToken := a.funcTemporaryKeyGet(token)

	if errToken != nil {
		api.Respond(w, r, api.Error("Link not valid of expired"))
		return
	}

	errPasswordChange := a.funcUserPasswordChange(userID, password)

	if errPasswordChange != nil {
		api.Respond(w, r, api.Error("authentication failed. "+errPasswordChange.Error()))
		return
	}

	// token := utils.RandStr(32)

	// errSession := a.funcUserStoreToken(token, userID)

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

	api.Respond(w, r, api.SuccessWithData("login success", map[string]interface{}{
		"token": token,
	}))
}
