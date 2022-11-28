package auth

import (
	"net/http"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/icons"
)

func (a Auth) pageRegister(w http.ResponseWriter, r *http.Request) {
	content := ""
	scripts := ""
	if a.passwordless {
		content = a.pageRegisterPasswordlessContent()
		scripts = a.pageRegisterPasswordlessScripts()
	} else {
		content = a.pageRegisterUsernameAndPasswordContent()
		scripts = a.pageRegisterUsernameAndPasswordScripts()
	}

	webpage := webpage("Register", a.funcLayout(content), scripts)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func (a Auth) pageRegisterPasswordlessContent() string {
	// Elements for the form
	alertSuccess := hb.NewDiv().Class("alert alert-success").Style("display:none")
	alertDanger := hb.NewDiv().Class("alert alert-danger").Style("display:none")
	alertGroup := hb.NewDiv().Class("alert-group").Child(alertSuccess).Child(alertDanger)

	header := hb.NewHeading5().HTML("Register").Style("margin:0px;")

	firstNameLabel := hb.NewLabel().HTML("First Name")
	firstNameInput := hb.NewInput().Class("form-control").Attr("name", "first_name").Attr("placeholder", "Enter first name")
	firstNameFormGroup := hb.NewDiv().Class("form-group mt-3").Child(firstNameLabel).Child(firstNameInput)

	lastNameLabel := hb.NewLabel().HTML("Last Name")
	lastNameInput := hb.NewInput().Class("form-control").Attr("name", "last_name").Attr("placeholder", "Enter last name")
	lastNameFormGroup := hb.NewDiv().Class("form-group mt-3").Child(lastNameLabel).Child(lastNameInput)

	emailLabel := hb.NewLabel().HTML("E-mail Address")
	emailInput := hb.NewInput().Class("form-control").Attr("name", "email").Attr("placeholder", "Enter e-mail address")
	emailFormGroup := hb.NewDiv().Class("form-group mt-3").AddChild(emailLabel).AddChild(emailInput)

	buttonRegister := hb.NewButton().Class("btn btn-lg btn-success btn-block w-100").Children([]*hb.Tag{
		icons.Icon("bi-person-circle", 24, 24, "white").Style("margin-right:8px;margin-top:-2px;"),
		hb.NewSpan().HTML("Register"),
	}).OnClick("registerFormValidate()")

	buttonRegisterFormGroup := hb.NewDiv().Class("form-group mt-3 mb-3").Child(buttonRegister)

	buttonLogin := hb.NewHyperlink().Class("btn btn-info text-white float-start").Children([]*hb.Tag{
		icons.Icon("bi-send", 16, 16, "white").Style("margin-right:8px;margin-top:-2px;"),
		hb.NewSpan().HTML("Login"),
	}).Attr("href", a.LinkLogin())

	// Add elements in a card
	cardHeader := hb.NewDiv().Class("card-header").AddChild(header)
	cardBody := hb.NewDiv().Class("card-body").AddChildren([]*hb.Tag{
		alertGroup,
		firstNameFormGroup,
		lastNameFormGroup,
		emailFormGroup,
		buttonRegisterFormGroup,
	})
	cardFooter := hb.NewDiv().Class("card-footer").Children([]*hb.Tag{
		buttonLogin,
	})
	card := hb.NewDiv().Class("card card-default").Style("margin:0 auto;max-width: 360px;")
	card.AddChild(cardHeader).Child(cardBody).Child(cardFooter)

	container := hb.NewDiv().Class("container").Child(card)

	return container.ToHTML()
}

func (a Auth) pageRegisterPasswordlessScripts() string {
	urlApiRegister := a.LinkApiRegister()
	urlSuccess := a.LinkRegisterCodeVerify()

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

        $('.buttonLogin .imgLoading').show();

        var data = {"first_name": first_name, "last_name": last_name, "email": email};

        $.post(urlApiRegister, data).then(function (response) {
            $('.buttonLogin .imgLoading').hide();

            if (response.status !== "success") {
                return registerFormRaiseError(response.message);
            }

            registerFormRaiseSuccess('Success');
            $('div.alert-danger').html('').hide();
            setTimeout(function () {
                window.location.href=urlOnSuccess;
            }, 100);
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

func (a Auth) pageRegisterUsernameAndPasswordContent() string {
	// Elements for the form
	alertSuccess := hb.NewDiv().Attr("class", "alert alert-success").Attr("style", "display:none")
	alertDanger := hb.NewDiv().Attr("class", "alert alert-danger").Attr("style", "display:none")
	alertGroup := hb.NewDiv().Attr("class", "alert-group").AddChild(alertSuccess).AddChild(alertDanger)

	header := hb.NewHeading5().HTML("Register").Attr("style", "margin:0px;")
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
	buttonRegisterFormGroup := hb.NewDiv().Attr("class", "form-group mt-3 mb-3").AddChild(buttonRegister)
	buttonLogin := hb.NewHyperlink().Attr("class", "btn btn-info float-start").HTML("Login").Attr("href", a.LinkLogin())
	buttonForgotPassword := hb.NewHyperlink().Attr("class", "btn btn-warning float-end").HTML("Forgot password?").Attr("href", a.LinkPasswordRestore())

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
	// heading := hb.NewHeading1().Attr("class", "text-center").HTML("Register")

	// container.AddChild(heading)
	container.AddChild(card)

	return container.ToHTML()
}

func (a Auth) pageRegisterUsernameAndPasswordScripts() string {
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
