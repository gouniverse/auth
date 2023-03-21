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

	registerMap := register.(map[string]interface{})

	email := ""
	if val, ok := registerMap["email"]; ok {
		email = val.(string)
	}

	firstName := ""
	if val, ok := registerMap["first_name"]; ok {
		firstName = val.(string)
	}

	lastName := ""
	if val, ok := registerMap["last_name"]; ok {
		lastName = val.(string)
	}

	password := ""
	if val, ok := registerMap["password"]; ok {
		password = val.(string)
	}

	var errRegister error = nil

	if a.passwordless {
		errRegister = a.passwordlessFuncUserRegister(email, firstName, lastName)
	} else {
		errRegister = a.funcUserRegister(email, password, firstName, lastName)
	}

	if errRegister != nil {
		api.Respond(w, r, api.Error("registration failed. "+errRegister.Error()))
		return
	}

	a.authenticateViaUsername(w, r, email, firstName, lastName)
}
