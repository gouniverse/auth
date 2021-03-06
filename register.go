package auth

import (
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/hb"
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

func (a Auth) pageRegister(w http.ResponseWriter, r *http.Request) {
	// Elements for the form
	alertSuccess := hb.NewDiv().Attr("class", "alert alert-success").Attr("style", "display:none")
	alertDanger := hb.NewDiv().Attr("class", "alert alert-danger").Attr("style", "display:none")
	alertGroup := hb.NewDiv().Attr("class", "alert-group").AddChild(alertSuccess).AddChild(alertDanger)

	header := hb.NewHeading3().HTML("Please sign up").Attr("style", "margin:0px;")
	firstNameLabel := hb.NewLabel().HTML("First Name")
	firstNameInput := hb.NewInput().Attr("class", "form-control").Attr("name", "first_name").Attr("placeholder", "Enter first name")
	firstNameFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(firstNameLabel).AddChild(firstNameInput)
	lastNameLabel := hb.NewLabel().HTML("Last Name")
	lastNameInput := hb.NewInput().Attr("class", "form-control").Attr("name", "last_name").Attr("placeholder", "Enter last name")
	lastNameFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(lastNameLabel).AddChild(lastNameInput)
	emailLabel := hb.NewLabel().HTML("E-mail Address")
	emailInput := hb.NewInput().Attr("class", "form-control").Attr("name", "email").Attr("placeholder", "Enter e-mail address")
	emailFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(emailLabel).AddChild(emailInput)
	passwordLabel := hb.NewLabel().AddChild(hb.NewHTML("Password"))
	passwordInput := hb.NewInput().Attr("class", "form-control").Attr("name", "password").Attr("type", "password").Attr("placeholder", "Enter password")
	passwordFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(passwordLabel).AddChild(passwordInput)
	buttonRegister := hb.NewButton().Attr("class", "btn btn-lg btn-success btn-block w-100").HTML("Register").Attr("onclick", "registerFormValidate()")
	buttonRegisterFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(buttonRegister)
	buttonLogin := hb.NewHyperlink().Attr("class", "btn btn-lg btn-info float-start").HTML("Login").Attr("href", a.LinkLogin())
	buttonForgotPassword := hb.NewHyperlink().Attr("class", "btn btn-lg btn-warning float-end").HTML("Forgot password?").Attr("href", a.LinkPasswordRestore())
	//form := hb.NewForm().Attr("method", "POST")

	// Add elements in a card
	cardHeader := hb.NewDiv().Attr("class", "card-header").AddChild(header)
	cardBody := hb.NewDiv().Attr("class", "card-body").AddChildren([]*hb.Tag{
		alertGroup,
		firstNameFormGroup,
		lastNameFormGroup,
		emailFormGroup,
		passwordFormGroup,
		buttonRegisterFormGroup,
	})
	cardFooter := hb.NewDiv().Attr("class", "card-footer").AddChildren([]*hb.Tag{
		buttonLogin,
		buttonForgotPassword,
	})
	card := hb.NewDiv().Attr("class", "card card-default").Attr("style", "margin:0 auto;max-width: 360px;")
	card.AddChild(cardHeader).AddChild(cardBody).AddChild(cardFooter)

	container := hb.NewDiv().Attr("class", "container")
	heading := hb.NewHeading1().Attr("class", "text-center").HTML("Register")

	container.AddChild(heading)
	container.AddChild(card)

	h := container.ToHTML()

	webpage := webpage("Register", h, a.pageRegisterScripts())
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func (a Auth) pageRegisterScripts() string {
	urlApiRegister := a.LinkApiRegister()
	urlSuccess := a.LinkLogin()

	return `
	var urlApiRegister = "` + urlApiRegister + `";
	console.log(urlApiRegister);
	var urlOnSuccess = "` + urlSuccess + `";
    /**
     * Raises an error message
     * @param  {String} error
     * @returns  {Boolean}
     */
    function registerFormRaiseError(error) {
        $('div.alert-success').html('').hide();
        $('div.alert-danger').html(error).show();
        setTimeout(function () {
            $('div.alert-danger').html('').hide();
        }, 10000);
        return false;
    }

    function registerFormRaiseSuccess(success) {
        $('div.alert-danger').html('').hide();
        $('div.alert-success').html(success).show();
        setTimeout(function () {
            $('div.alert-success').html('').hide();
        }, 10000);
        return false;
    }

    /**
     * Validate Login Form
     * @returns  {Boolean}
     */
    function registerFormValidate() {
		var first_name = $.trim($('input[name=first_name]').val());
		var last_name = $.trim($('input[name=last_name]').val());
        var email = $.trim($('input[name=email]').val());
        var password = $.trim($('input[name=password]').val());

		if (first_name === '') {
            return registerFormRaiseError('First name is required');
        }

		if (last_name === '') {
            return registerFormRaiseError('Last name is required');
        }

        if (email === '') {
            return registerFormRaiseError('Email is required');
        }

        if (password === '') {
            return registerFormRaiseError('Password is required');
        }

        $('.buttonLogin .imgLoading').show();

        var data = {"first_name": first_name, "last_name": last_name, "email": email, "password": password};

        $.post(urlApiRegister, data).then(function (response) {
            $('.buttonLogin .imgLoading').hide();

            if (response.status !== "success") {
                return registerFormRaiseError(response.message);
            }

            registerFormRaiseSuccess('Success');
            $('div.alert-danger').html('').hide();
            setTimeout(function () {
                window.location.href=urlOnSuccess;
            }, 2000);
            return;
        }).fail(function (error) {
			console.log(error);
            $('.buttonLogin .imgLoading').hide();
            return registerFormRaiseError('There was an error. Try again later!');
        });
    }
    $(function () {
        $("input[name=first_name").focus();
    });
	`
}
