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

func (a Auth) apiPaswordRestore(w http.ResponseWriter, r *http.Request) {
	email := strings.Trim(utils.Req(r, "email", ""), " ")
	firstName := strings.Trim(utils.Req(r, "first_name", ""), " ")
	lastName := strings.Trim(utils.Req(r, "last_name", ""), " ")

	if email == "" {
		api.Respond(w, r, api.Error("Email is required field"))
		return
	}
	
	if firstName == "" {
		api.Respond(w, r, api.Error("First name is required field"))
		return
	}

	if lastName == "" {
		api.Respond(w, r, api.Error("Last name is required field"))
		return
	}

	user, err := a.funcUserFindByEmail(email)

	if err != nil {
		log.Println(err.Error())
		api.Respond(w, r, api.Error("Internal server error"))
		return
	}

	if user == nil {
		api.Respond(w, r, api.Error("E-mail not registered"))
		return
	}

	if strings.ToLower(user.FirstName) != strings.ToLower(firstName) {
		api.Respond(w, r, api.Error("First or last name not matching"))
		return
	}
	
	
	if strings.ToLower(user.LastName) != strings.ToLower(lastName) {
		api.Respond(w, r, api.Error("First or last name not matching"))
		return
	}
	
	token := utils.RandStr(32)

	errTempTokenSave := a.funcStoreTempToken(token, user.Email, 3600)

	if errTempTokenSave != nil {
		api.Respond(w, r, api.Error("token store failed. "+errSession.Error()))
		return
	}
	
	emailContent := helpers.EmailPasswordChangeTemplate(user.Email, helpers.LinkAuthPasswordChange(emailPasswordChangeToken))

	isSent, err := a.funcEmailSend(a.EmailFromAddress, user.Email, "Password Restore", emailContent)

	log.Println(err)

	if isSent {
		api.Respond(w, r, api.Success("Password reset link was sent to your e-mail"))
		return
	}

	api.Respond(w, r, api.Error("Password reset link failed to be sent. Please try again later"))
}

func (a Auth) pagePasswordRestore(w http.ResponseWriter, r *http.Request) {
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
	buttonContinue := hb.NewButton().Attr("class", "btn btn-lg btn-success btn-block w-100").HTML("Email Password Reset Link").Attr("onclick", "registerFormValidate()")
	buttonContinueFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(buttonContinue)
	buttonLogin := hb.NewHyperlink().Attr("class", "btn btn-lg btn-info float-start").HTML("Login").Attr("href", a.LinkLogin())
	buttonRegister := hb.NewHyperlink().Attr("class", "btn btn-lg btn-warning float-end").HTML("Register").Attr("href", a.LinkRegister())
	//form := hb.NewForm().Attr("method", "POST")

	// Add elements in a card
	cardHeader := hb.NewDiv().Attr("class", "card-header").AddChild(header)
	cardBody := hb.NewDiv().Attr("class", "card-body").AddChildren([]*hb.Tag{
		alertGroup,
		firstNameFormGroup,
		lastNameFormGroup,
		emailFormGroup,
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

	webpage := webpage("Fogotten Password", h, a.pagePasswordRestoreScripts())
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func (a Auth) pagePasswordRestoreScripts() string {
	urlApiRegister := a.LinkApiRegister()
	urlSuccess := a.LinkLogin()

	return `
	var urlApiRegister = "` + urlApiRegister + `";
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
