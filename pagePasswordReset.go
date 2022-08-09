package auth

import (
	"net/http"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/utils"
)

func (a Auth) pagePasswordReset(w http.ResponseWriter, r *http.Request) {
	token := utils.Req(r, "t", "")
	errorMessage := ""

	if token == "" {
		errorMessage = "Link is invalid"
	} else {
		tokenValue, errToken := a.funcTemporaryKeyGet(token)
		if errToken != nil {
			errorMessage = "Link has expired"
		} else if tokenValue == "" {
			errorMessage = "Link is invalid or expired"
		}
	}

	h := a.pagePasswordResetContent(token, errorMessage)
	webpage := webpage("Reset Password", h, a.pagePasswordResetScripts())
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func (a Auth) pagePasswordResetContent(token string, errorMessage string) string {
	urlPasswordRestore := a.LinkPasswordRestore()
	urlLogin := a.LinkLogin()
	urlRegister := a.LinkRegister()
	// Elements for the form
	alertSuccess := hb.NewDiv().Attr("class", "alert alert-success").Attr("style", "display:none")
	alertDanger := hb.NewDiv().Attr("class", "alert alert-danger")
	if errorMessage != "" {
		alertDanger.HTML(errorMessage)
	} else {
		alertDanger.Attr("style", "display:none")
	}
	alertGroup := hb.NewDiv().Attr("class", "alert-group").AddChild(alertSuccess).AddChild(alertDanger)

	header := hb.NewHeading3().HTML("Please fill the form bellow").Attr("style", "margin:0px;")
	tokenInput := hb.NewInput().Attr("name", "token").Attr("value", token)
	passwordLabel := hb.NewLabel().HTML("New Password")
	passwordInput := hb.NewInput().Attr("class", "form-control").Attr("name", "password").Attr("placeholder", "Enter new password")
	passwordFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(passwordLabel).AddChild(passwordInput)
	passwordConfirmLabel := hb.NewLabel().HTML("Confirm New Password")
	passwordConfirmInput := hb.NewInput().Attr("class", "form-control").Attr("name", "password_confirm").Attr("placeholder", "Enter confirnation of new password")
	passwordConfirmFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(passwordConfirmLabel).AddChild(passwordConfirmInput)
	buttonContinue := hb.NewButton().Attr("class", "ButtonContinue btn btn-lg btn-success btn-block w-100").HTML("Reset Password").Attr("onclick", "resetFormValidate()")
	buttonContinueFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(buttonContinue)
	buttonLogin := hb.NewHyperlink().Attr("class", "btn btn-lg btn-info float-start").HTML("Login").Attr("href", a.LinkLogin())
	buttonRegister := hb.NewHyperlink().Attr("class", "btn btn-lg btn-warning float-end").HTML("Register").Attr("href", a.LinkRegister())
	//form := hb.NewForm().Attr("method", "POST")

	// Add elements in a card
	cardHeader := hb.NewDiv().Attr("class", "card-header").AddChild(header)
	cardBody := hb.NewDiv().Attr("class", "card-body").AddChildren([]*hb.Tag{
		alertGroup,
	})

	if errorMessage == "" {
		cardBody.AddChild(tokenInput)
		cardBody.AddChild(passwordFormGroup)
		cardBody.AddChild(passwordConfirmFormGroup)
		cardBody.AddChild(buttonContinueFormGroup)
	} else {
		cardBody.AddChild(hb.NewParagraph().HTML("Sorry, there was an error processing your request. Please select one of the following options:"))
		cardBody.AddChild(hb.NewParagraph().AddChild(hb.NewHyperlink().Attr("href", urlPasswordRestore).HTML("request a reset of your password")))
		cardBody.AddChild(hb.NewParagraph().AddChild(hb.NewHyperlink().Attr("href", urlLogin).HTML("login to the system")))
		cardBody.AddChild(hb.NewParagraph().AddChild(hb.NewHyperlink().Attr("href", urlRegister).HTML("create a new account")))
	}

	cardFooter := hb.NewDiv().Attr("class", "card-footer").AddChildren([]*hb.Tag{
		buttonLogin,
	})

	if a.enableRegistration {
		cardFooter.AddChild(buttonRegister)
	}

	card := hb.NewDiv().
		Attr("class", "card card-default").
		Attr("style", "margin:0 auto;max-width: 360px;")

	card.AddChild(cardHeader).AddChild(cardBody).AddChild(cardFooter)

	container := hb.NewDiv().Attr("class", "container")
	heading := hb.NewHeading1().Attr("class", "text-center").HTML("Change Password")

	container.AddChild(heading)
	container.AddChild(card)

	return container.ToHTML()
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
    function resetFormRaiseError(error) {
        $('div.alert-success').html('').hide();
        $('div.alert-danger').html(error).show();
        setTimeout(function () {
            $('div.alert-danger').html('').hide();
        }, 10000);
        return false;
    }

    function resetFormRaiseSuccess(success) {
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
    function resetFormValidate() {
        var token = $.trim($('input[name=token]').val());
        var password = $.trim($('input[name=password]').val());
        var passwordConfirm = $.trim($('input[name=password_confirm]').val());

        $('.ButtonContinue .imgLoading').show();

        var data = {"password": password, "password_confirm": passwordConfirm, "token": token};

        $.post(urlApiPasswordReset, data).then(function (response) {
            $('.ButtonContinue .imgLoading').hide();

            if (response.status !== "success") {
                return resetFormRaiseError(response.message);
            }

            resetFormRaiseSuccess('Success');
            $('div.alert-danger').html('').hide();
            setTimeout(function () {
                window.location.href=urlOnSuccess;
            }, 2000);
            return;
        }).fail(function (error) {
			console.log(error);
            $('.ButtonContinue .imgLoading').hide();
            return resetFormRaiseError('There was an error. Try again later!');
        });
    }
    $(function () {
        $("input[name=first_name").focus();
    });
	`
}
