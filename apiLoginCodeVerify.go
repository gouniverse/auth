package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (a Auth) apiLoginCodeVerify(w http.ResponseWriter, r *http.Request) {
	verificationCode := strings.Trim(utils.Req(r, "verification_code", ""), " ")

	if verificationCode == "" {
		api.Respond(w, r, api.Error("Verification code is required field"))
		return
	}

	if len(verificationCode) != LoginCodeLength {
		api.Respond(w, r, api.Error("Verification code is invalid length"))
		return
	}

	if !strCotainsOnlyCharacters(verificationCode, LoginCodeGamma) {
		api.Respond(w, r, api.Error("Verification code contains invalid characters"))
		return
	}

	email, errCode := a.funcTemporaryKeyGet(verificationCode)

	if errCode != nil {
		api.Respond(w, r, api.Error("Verification code has expired"))
		return
	}

	userID, errUser := a.passwordlessFuncUserFindByEmail(email)

	if errUser != nil {
		api.Respond(w, r, api.Error("authentication failed. "+errUser.Error()))
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

func strCotainsOnlyCharacters(str string, chars string) bool {
	nonCharacter := func(c rune) bool { return !strings.ContainsAny(chars, string(c)) }
	words := strings.FieldsFunc(str, nonCharacter)
	return str == strings.Join(words, "")
}
