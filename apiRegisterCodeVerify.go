package auth

import (
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (a Auth) apiRegisterCodeVerify(w http.ResponseWriter, r *http.Request) {
	verificationCode := strings.Trim(utils.Req(r, "verification_code", ""), " ")

	if verificationCode == "" {
		api.Respond(w, r, api.Error("Verification code is required field"))
		return
	}

	if len(verificationCode) != LoginCodeLength {
		api.Respond(w, r, api.Error("Verification code is invalid length"))
		return
	}

	if !utils.StrContainsOnlySpecifiedCharacters(verificationCode, LoginCodeGamma) {
		api.Respond(w, r, api.Error("Verification code contains invalid characters"))
		return
	}

	registerJSON, errCode := a.funcTemporaryKeyGet(verificationCode)

	if errCode != nil {
		api.Respond(w, r, api.Error("Verification code has expired"))
		return
	}

	register, errJSON := utils.FromJSON(registerJSON, nil)

	if errJSON != nil {
		api.Respond(w, r, api.Error("Serialized format is malformed"))
		return
	}

	email := register.(map[string]interface{})["email"].(string)
	firstName := register.(map[string]interface{})["first_name"].(string)
	lastName := register.(map[string]interface{})["last_name"].(string)

	errRegister := a.passwordlessFuncUserRegister(email, firstName, lastName)

	if errRegister != nil {
		api.Respond(w, r, api.Error("registration failed. "+errRegister.Error()))
		return
	}

	a.authenticateViaEmail(w, r, email)
}
