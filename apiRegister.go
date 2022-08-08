package auth

import (
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
	validator "github.com/gouniverse/validator"
)

func (a Auth) apiRegister(w http.ResponseWriter, r *http.Request) {
	email := strings.Trim(utils.Req(r, "email", ""), " ")
	password := strings.Trim(utils.Req(r, "password", ""), " ")
	first_name := strings.Trim(utils.Req(r, "first_name", ""), " ")
	last_name := strings.Trim(utils.Req(r, "last_name", ""), " ")

	if first_name == "" {
		api.Respond(w, r, api.Error("First name is required field"))
		return
	}

	if last_name == "" {
		api.Respond(w, r, api.Error("Last name is required field"))
		return
	}

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

	err := a.funcUserRegister(email, password, first_name, last_name)

	if err != nil {
		api.Respond(w, r, api.Error("registration failed. "+err.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("registration success", map[string]interface{}{}))
}
