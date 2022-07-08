package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/gouniverse/api"
	"github.com/gouniverse/hb"
)

func (a Auth) apiLogout(w http.ResponseWriter, r *http.Request) {
	authToken := authTokenRetrieve(r, a.useCookies)

	if authToken == "" {
		api.Respond(w, r, api.Success("logout success"))
	}

	userID, errToken := a.funcUserFindByToken(authToken)

	if errToken != nil {
		api.Respond(w, r, api.Error("logout failed"))
		return
	}

	errLogout := a.funcUserLogout(userID)

	if errLogout != nil {
		api.Respond(w, r, api.Error("logout failed. "+errLogout.Error()))
		return
	}

	if a.useCookies {
		expiration := time.Now().Add(-365 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:     "authtoken",
			Value:    "none",
			Expires:  expiration,
			HttpOnly: false,
			Secure:   true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)
	}

	api.Respond(w, r, api.Success("logout success"))
}

func (a Auth) pageLogout(w http.ResponseWriter, r *http.Request) {
	// Elements for the form
	alertSuccess := hb.NewDiv().Attr("class", "alert alert-success").Attr("style", "display:none")
	alertDanger := hb.NewDiv().Attr("class", "alert alert-danger").Attr("style", "display:none")
	alertGroup := hb.NewDiv().Attr("class", "alert-group").AddChild(alertSuccess).AddChild(alertDanger)

	header := hb.NewHeading3().HTML("Please sign out").Attr("style", "margin:0px;")
	buttonContinue := hb.NewButton().Attr("class", "btn btn-lg btn-success btn-block w-100").HTML("Logout").Attr("onclick", "logoutFormValidate()")
	buttonContinueFormGroup := hb.NewDiv().Attr("class", "form-group mt-3").AddChild(buttonContinue)

	// Add elements in a card
	cardHeader := hb.NewDiv().Attr("class", "card-header").AddChild(header)
	cardBody := hb.NewDiv().Attr("class", "card-body").AddChildren([]*hb.Tag{
		alertGroup,
		buttonContinueFormGroup,
	})

	card := hb.NewDiv().Attr("class", "card card-default").Attr("style", "margin:0 auto;max-width: 360px;")
	card.AddChild(cardHeader).AddChild(cardBody)

	container := hb.NewDiv().Attr("class", "container")
	heading := hb.NewHeading1().Attr("class", "text-center").HTML("Logout")

	container.AddChild(heading)
	container.AddChild(card)

	h := container.ToHTML()

	webpage := webpage("Logout", h, a.pageLogoutScripts())
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func (a Auth) pageLogoutScripts() string {
	urlApiLogout := a.LinkApiLogout()
	urlSuccess := a.LinkLogin()
	log.Println(urlApiLogout)

	return `
	var urlApiLogout = "` + urlApiLogout + `";
	var urlOnSuccess = "` + urlSuccess + `";
    /**
     * Raises an error message
     * @param  {String} error
     * @returns  {Boolean}
     */
    function logoutFormRaiseError(error) {
        $('div.alert-success').html('').hide();
        $('div.alert-danger').html(error).show();
        setTimeout(function () {
            $('div.alert-danger').html('').hide();
        }, 10000);
        return false;
    }

    function logoutFormRaiseSuccess(success) {
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
    function logoutFormValidate() {
        $('.buttonLogin .imgLoading').show();

        var data = {};

        $.post(urlApiLogout, data).then(function (response) {
            $('.buttonLogin .imgLoading').hide();

            if (response.status !== "success") {
                return logoutFormRaiseError(response.message);
            }

            $$.setAuthToken(response.data.token);
            $$.setAuthUser(response.data.user);
            logoutFormRaiseSuccess('Success');
            $('div.alert-danger').html('').hide();
            setTimeout(function () {
                $$.to(urlOnSuccess);
            }, 2000);
            return;
        }).fail(function (error) {
			console.log(error);
            $('.buttonLogin .imgLoading').hide();
            return logoutFormRaiseError('There was an error. Try again later!');
        });
    }
    $(function () {
        $("#email").focus();
    });
	`
}
