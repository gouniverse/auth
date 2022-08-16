package auth

import (
	"net/http"

	"github.com/gouniverse/hb"
)

func (a Auth) pageLogin(w http.ResponseWriter, r *http.Request) {
	webpage := webpage("Login", a.funcLayout(a.pageLoginContent()), a.pageLoginScripts())
	
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func (a Auth) pageLoginContent() string {
	// Elements for the form
	alertSuccess := hb.NewDiv().Attr("class", "alert alert-success").Attr("style", "display:none")
	alertDanger := hb.NewDiv().Attr("class", "alert alert-danger").Attr("style", "display:none")
	alertGroup := hb.NewDiv().Attr("class", "alert-group").AddChild(alertSuccess).AddChild(alertDanger)

	header := hb.NewHeading5().HTML("Login").Attr("style", "margin:0px;")
	emailLabel := hb.NewLabel().HTML("E-mail Address")
	emailInput := hb.NewInput().Attr("class", "form-control").Attr("name", "email").Attr("placeholder", "Enter e-mail address")
	emailFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(emailLabel).AddChild(emailInput)
	passwordLabel := hb.NewLabel().AddChild(hb.NewHTML("Password"))
	passwordInput := hb.NewInput().Attr("class", "form-control").Attr("name", "password").Attr("type", "password").Attr("placeholder", "Enter password")
	passwordFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(passwordLabel).AddChild(passwordInput)
	buttonLogin := hb.NewButton().Attr("class", "btn btn-lg btn-success btn-block w-100").HTML("Login").Attr("onclick", "loginFormValidate()")
	buttonLoginFormGroup := hb.NewDiv().Attr("class", "form-group mt-3 mb-3").AddChild(buttonLogin)
	buttonRegister := hb.NewHyperlink().Attr("class", "btn btn-info float-start").HTML("Register").Attr("href", a.LinkRegister())
	buttonForgotPassword := hb.NewHyperlink().Attr("class", "btn btn-warning float-end").HTML("Forgot password?").Attr("href", a.LinkPasswordRestore())

	// Add elements in a card
	cardHeader := hb.NewDiv().Attr("class", "card-header").AddChild(header)
	cardBody := hb.NewDiv().Attr("class", "card-body").
		// Attr("style", "margin-bottom:20px;").
		AddChildren([]*hb.Tag{
			alertGroup,
			emailFormGroup,
			passwordFormGroup,
			buttonLoginFormGroup,
		})
	cardFooter := hb.NewDiv().Attr("class", "card-footer").AddChildren([]*hb.Tag{
		buttonForgotPassword,
	})

	if a.enableRegistration {
		cardFooter.AddChild(buttonRegister)
	}

	card := hb.NewDiv().Attr("class", "card card-default").
		Attr("style", "margin:0 auto;max-width: 360px;").
		AddChild(cardHeader).
		AddChild(cardBody).
		AddChild(cardFooter)

	container := hb.NewDiv().Attr("class", "container").
		AddChild(card)

	return container.ToHTML()
}

func (a Auth) pageLoginScripts() string {
	urlApiLogin := a.LinkApiLogin()
	urlSuccess := a.LinkRedirectOnSuccess()

	return `
	var urlApiLogin = "` + urlApiLogin + `";
	var urlOnSuccess = "` + urlSuccess + `";
    /**
     * Raises an error message
     * @param  {String} error
     * @returns  {Boolean}
     */
    function loginFormRaiseError(error) {
        $('div.alert-success').html('').hide();
        $('div.alert-danger').html(error).show();
        setTimeout(function () {
            $('div.alert-danger').html('').hide();
        }, 10000);
        return false;
    }

    function loginFormRaiseSuccess(success) {
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
    function loginFormValidate() {
        var email = $.trim($('input[name=email]').val());
        var password = $.trim($('input[name=password]').val());

        if (email === '') {
            return loginFormRaiseError('Email is required');
        }

        if (password === '') {
            return loginFormRaiseError('Password is required');
        }

        $('.buttonLogin .imgLoading').show();

        var data = {"email": email, "password": password};

        $.post(urlApiLogin, data).then(function (response) {
            $('.buttonLogin .imgLoading').hide();

            if (response.status !== "success") {
                return loginFormRaiseError(response.message);
            }

            $$.setAuthToken(response.data.token);
            $$.setAuthUser(response.data.user);
            loginFormRaiseSuccess('Success');
            $('div.alert-danger').html('').hide();
            setTimeout(function () {
                $$.to(urlOnSuccess);
            }, 2000);
            return;
        }).fail(function (error) {
			console.log(error);
            $('.buttonLogin .imgLoading').hide();
            return loginFormRaiseError('There was an error. Try again later!');
        });
    }
    $(function () {
        $("#email").focus();
    });
	`
}
