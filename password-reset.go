package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gouniverse/api"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/utils"
	validator "github.com/gouniverse/validator"
)

func (a Auth) apiPaswordReset(w http.ResponseWriter, r *http.Request) {
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

func (a Auth) pagePasswordReset(w http.ResponseWriter, r *http.Request) {
	// Elements for the form
	alertSuccess := hb.NewDiv().Attr("class", "alert alert-success").Attr("style", "display:none")
	alertDanger := hb.NewDiv().Attr("class", "alert alert-danger").Attr("style", "display:none")
	alertGroup := hb.NewDiv().Attr("class", "alert-group").AddChild(alertSuccess).AddChild(alertDanger)

	header := hb.NewHeading3().HTML("Please fill the form bellow").Attr("style", "margin:0px;")
	passwordLabel := hb.NewLabel().HTML("New Password")
	passwordInput := hb.NewInput().Attr("class", "form-control").Attr("name", "password").Attr("placeholder", "Enter new password")
	passwordFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(passwordLabel).AddChild(passwordInput)
	passwordConfirmLabel := hb.NewLabel().HTML("Confirm New Password")
	passwordConfirmInput := hb.NewInput().Attr("class", "form-control").Attr("name", "password_confirm").Attr("placeholder", "Enter confirnation of new password")
	passwordConfirmFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(passwordConfirmLabel).AddChild(passwordConfirmInput)
	buttonContinue := hb.NewButton().Attr("class", "btn btn-lg btn-success btn-block w-100").HTML("Reset Password").Attr("onclick", "registerFormValidate()")
	buttonContinueFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(buttonContinue)
	buttonLogin := hb.NewHyperlink().Attr("class", "btn btn-lg btn-info float-start").HTML("Login").Attr("href", a.LinkLogin())
	buttonRegister := hb.NewHyperlink().Attr("class", "btn btn-lg btn-warning float-end").HTML("Register").Attr("href", a.LinkRegister())
	//form := hb.NewForm().Attr("method", "POST")

	// Add elements in a card
	cardHeader := hb.NewDiv().Attr("class", "card-header").AddChild(header)
	cardBody := hb.NewDiv().Attr("class", "card-body").AddChildren([]*hb.Tag{
		alertGroup,
		passwordFormGroup,
		passwordConfirmFormGroup,
		buttonContinueFormGroup,
	})
	cardFooter := hb.NewDiv().Attr("class", "card-footer").AddChildren([]*hb.Tag{
		buttonLogin,
	})

	if a.enableRegistration {
		cardFooter.AddChild(buttonRegister)
	}

	card := hb.NewDiv().Attr("class", "card card-default").Attr("style", "margin:0 auto;max-width: 360px;")
	card.AddChild(cardHeader).AddChild(cardBody).AddChild(cardFooter)

	container := hb.NewDiv().Attr("class", "container")
	heading := hb.NewHeading1().Attr("class", "text-center").HTML("Forgot Password")

	container.AddChild(heading)
	container.AddChild(card)

	h := container.ToHTML()

	webpage := webpage("Reset Password", h, a.pagePasswordResetScripts())
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func (a Auth) pagePasswordResetScripts() string {
	urlApiPasswordReset := a.LinkApiPasswordReset()
	urlSuccess := a.LinkLogin()

	return `
	var urlApiPasswordReset = "` + urlApiPasswordReset + `";
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

        $.post(urlApiPasswordReset, data).then(function (response) {
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
