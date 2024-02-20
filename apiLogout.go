package auth

import (
	"net/http"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (a Auth) apiLogout(w http.ResponseWriter, r *http.Request) {
	authToken := AuthTokenRetrieve(r, a.useCookies)

	if authToken == "" {
		api.Respond(w, r, api.Success("logout success"))
	}

	userID, errToken := a.funcUserFindByAuthToken(authToken, utils.IP(r), r.UserAgent())

	if errToken != nil {
		api.Respond(w, r, api.Error("logout failed"))
		return
	}

	if userID != "" {
		errLogout := a.funcUserLogout(userID)

		if errLogout != nil {
			api.Respond(w, r, api.Error("logout failed. "+errLogout.Error()))
			return
		}
	}

	if a.useCookies {
		AuthCookieRemove(w, r)
	}

	api.Respond(w, r, api.Success("logout success"))
}
