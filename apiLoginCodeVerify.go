package auth

import (
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (a Auth) apiLoginCodeVerify(w http.ResponseWriter, r *http.Request) {
	verificationCode := strings.Trim(utils.Req(r, "verification_code", ""), " ")

	if verificationCode == "" {
		api.Respond(w, r, api.Error("Verification code is required field"))
		return
	}

	api.Respond(w, r, api.Error("Verification code not implemented yet"))
}
