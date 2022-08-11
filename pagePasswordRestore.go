package auth

import (
	"net/http"

	"github.com/gouniverse/hb"
)

func (a Auth) pagePasswordRestore(w http.ResponseWriter, r *http.Request) {

	webpage := webpage("Restore Password", a.pagePasswordRestoreContent(), a.pagePasswordRestoreScripts())
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func (a Auth) pagePasswordRestoreContent() string {
	// Elements for the form
	alertSuccess := hb.NewDiv().Attr("class", "alert alert-success").Attr("style", "display:none")
	alertDanger := hb.NewDiv().Attr("class", "alert alert-danger").Attr("style", "display:none")
	alertGroup := hb.NewDiv().Attr("class", "alert-group").AddChild(alertSuccess).AddChild(alertDanger)

	header := hb.NewHeading5().
		Attr("style", "margin:0px;").
		HTML("Restore password")
	firstNameLabel := hb.NewLabel().HTML("First Name")
	firstNameInput := hb.NewInput().Attr("class", "form-control").Attr("name", "first_name").Attr("placeholder", "Enter first name")
	firstNameFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(firstNameLabel).AddChild(firstNameInput)
	lastNameLabel := hb.NewLabel().HTML("Last Name")
	lastNameInput := hb.NewInput().Attr("class", "form-control").Attr("name", "last_name").Attr("placeholder", "Enter last name")
	lastNameFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(lastNameLabel).AddChild(lastNameInput)
	emailLabel := hb.NewLabel().HTML("E-mail Address")
	emailInput := hb.NewInput().Attr("class", "form-control").Attr("name", "email").Attr("placeholder", "Enter e-mail address")
	emailFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(emailLabel).AddChild(emailInput)
	buttonContinue := hb.NewButton().
		Attr("class", "ButtonContinue btn btn-lg btn-success btn-block w-100").
		Attr("onclick", "passwordRestoreFormValidate()").
		HTML("Send Password Reset Link")
	buttonContinueFormGroup := hb.NewDiv().
		Attr("class", "form-group mt-3 mb-3").
		AddChild(buttonContinue)
	buttonLogin := hb.NewHyperlink().
		Attr("class", "btn btn-info float-start").
		Attr("href", a.LinkLogin()).
		HTML("Login")
	buttonRegister := hb.NewHyperlink().
		Attr("class", "btn btn-warning float-end").
		Attr("href", a.LinkRegister()).
		HTML("Register")

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

	container := hb.NewDiv().Attr("class", "container").
		AddChild(card)

	return container.ToHTML()
}

func (a Auth) pagePasswordRestoreScripts() string {
	urlApiPasswordrestore := a.LinkApiPasswordRestore()
	urlSuccess := a.LinkLogin()

	return `
	var urlApiPasswordRestore = "` + urlApiPasswordrestore + `";
	var urlOnSuccess = "` + urlSuccess + `";

    /**
     * Raises an error message
     * @param  {String} error
     * @returns  {Boolean}
     */
    function passwordRestoreFormRaiseError(error) {
        $('div.alert-success').html('').hide();
        $('div.alert-danger').html(error).show();
        setTimeout(function () {
            $('div.alert-danger').html('').hide();
        }, 10000);
        return false;
    }

    function passwordRestoreFormRaiseSuccess(success) {
        $('div.alert-danger').html('').hide();
        $('div.alert-success').html(success).show();
        setTimeout(function () {
            $('div.alert-success').html('').hide();
        }, 10000);
        return false;
    }

    /**
     * Validates the Password Restore Form
	 * and sends the request to the backend
     * @returns  {Boolean}
     */
    function passwordRestoreFormValidate() {
		var first_name = $.trim($('input[name=first_name]').val());
		var last_name = $.trim($('input[name=last_name]').val());
        var email = $.trim($('input[name=email]').val());

        $('.ButtonContinue .imgLoading').show();

        var data = {"first_name": first_name, "last_name": last_name, "email": email};

        $.post(urlApiPasswordRestore, data).then(function (response) {
            $('.ButtonContinue .imgLoading').hide();

            if (response.status !== "success") {
                return passwordRestoreFormRaiseError(response.message);
            }

            passwordRestoreFormRaiseSuccess('Success');
            $('div.alert-danger').html('').hide();

            setTimeout(function () {
                window.location.href=urlOnSuccess;
            }, 2000);
            return;
        }).fail(function (error) {
			console.log(error);
            $('.ButtonContinue .imgLoading').hide();
            return passwordRestoreFormRaiseError('There was an error. Try again later!');
        });
    }
    $(function () {
        $("input[name=first_name").focus();
    });
	`
}
