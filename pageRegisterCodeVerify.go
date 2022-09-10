package auth

import (
	"net/http"

	"github.com/gouniverse/hb"
)

func (a Auth) pageRegisterCodeVerify(w http.ResponseWriter, r *http.Request) {
	webpage := webpage("Verify Registration Code", a.funcLayout(a.pageRegisterCodeVerifyContent()), a.pageRegisterCodeVerifyScripts())

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func (a Auth) pageRegisterCodeVerifyContent() string {
	// Elements for the form
	alertSuccess := hb.NewDiv().Class("alert alert-success").Style("display:none")
	alertDanger := hb.NewDiv().Class("alert alert-danger").Style("display:none")
	alertGroup := hb.NewDiv().Class("alert-group").AddChild(alertSuccess).AddChild(alertDanger)

	header := hb.NewHeading5().HTML("Registration Code Verification").Attr("style", "margin:0px;")
	verificationCodeLabel := hb.NewLabel().HTML("Verification code")
	verificationCodeInput := hb.NewInput().Attr("class", "form-control").Attr("name", "verification_code").Attr("placeholder", "Enter verification code")
	verificationCodeFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").Child(verificationCodeLabel).Child(verificationCodeInput)
	buttonLogin := hb.NewButton().Attr("class", "btn btn-lg btn-success btn-block w-100").HTML("Login").Attr("onclick", "registerCodeFormValidate()")
	buttonLoginFormGroup := hb.NewDiv().Attr("class", "form-group mt-3 mb-3").AddChild(buttonLogin)
	buttonBack := hb.NewHyperlink().Attr("class", "btn btn-info float-start").HTML("Resend code").Attr("href", a.LinkLogin())

	// Add elements in a card
	cardHeader := hb.NewDiv().Class("card-header").Child(header)
	cardBody := hb.NewDiv().Class("card-body").Children([]*hb.Tag{
		alertGroup,
		verificationCodeFormGroup,
		buttonLoginFormGroup,
	})
	cardFooter := hb.NewDiv().Class("card-footer").Children([]*hb.Tag{
		buttonBack,
	})

	card := hb.NewDiv().Class("card card-default").Style("margin:0 auto;max-width: 360px;").Children([]*hb.Tag{
		cardHeader,
		cardBody,
		cardFooter,
	})

	container := hb.NewDiv().Class("container").Child(card)

	return container.ToHTML()
}

func (a Auth) pageRegisterCodeVerifyScripts() string {
	urlApiRegisterCodeVerify := a.LinkApiRegisterCodeVerify()
	urlSuccess := a.LinkRedirectOnSuccess()

	return `
	var urlApiRegisterCodeVerify = "` + urlApiRegisterCodeVerify + `";
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
    function registerCodeFormValidate() {
        var verificationCode = $.trim($('input[name=verification_code]').val());

        if (verificationCode === '') {
            return loginFormRaiseError('Code is required');
        }

        $('.buttonLogin .imgLoading').show();

        var data = {"verification_code": verificationCode};

        $.post(urlApiRegisterCodeVerify, data).then(function (response) {
            $('.buttonLogin .imgLoading').hide();

            if (response.status !== "success") {
                return loginFormRaiseError(response.message);
            }

            $$.setAuthToken(response.data.token);
            $$.setAuthUser(response.data.user);
            loginFormRaiseSuccess('Verification successful');
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