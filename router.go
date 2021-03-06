package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gouniverse/utils"
)

func (a Auth) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Router(w, r)
	})
}

// Router routes the requests
func (a Auth) Router(w http.ResponseWriter, r *http.Request) {
	path := utils.Req(r, "path", "home")
	uri := r.RequestURI

	if strings.HasSuffix(uri, pathApiLogin) {
		path = pathApiLogin
	} else if strings.HasSuffix(uri, pathApiLogout) {
		path = pathApiLogout
	} else if strings.HasSuffix(uri, pathApiResetPassword) {
		path = pathApiResetPassword
	} else if strings.HasSuffix(uri, pathApiRestorePassword) {
		path = pathApiRestorePassword
	} else if strings.HasSuffix(uri, pathApiRegister) {
		path = pathApiRegister
	} else if strings.HasSuffix(uri, pathLogin) {
		path = pathLogin
	} else if strings.HasSuffix(uri, pathLogout) {
		path = pathLogout
	} else if strings.HasSuffix(uri, pathRegister) {
		path = pathRegister
	} else if strings.HasSuffix(uri, pathPasswordRestore) {
		path = pathPasswordRestore
	}

	log.Println("Path: " + path)

	ctx := context.WithValue(r.Context(), keyEndpoint, r.URL.Path)

	routeFunc := a.getRoute(path)
	routeFunc(w, r.WithContext(ctx))
}

// getRoute finds a route
func (a Auth) getRoute(route string) func(w http.ResponseWriter, r *http.Request) {
	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		pathApiLogin:           a.apiLogin,
		pathApiLogout:          a.apiLogout,
		pathApiRegister:        a.apiRegister,
		pathApiResetPassword:   a.apiPaswordReset,
		pathApiRestorePassword: a.apiPaswordRestore,
		pathLogin:              a.pageLogin,
		pathLogout:             a.pageLogout,
		pathPasswordReset:      a.pagePasswordReset,
		pathPasswordRestore:    a.pagePasswordRestore,
	}

	if a.enableRegistration {
		routes[pathRegister] = a.pageRegister
	}

	if val, ok := routes[route]; ok {
		return val
	}

	return a.notFoundHandler
}

func (a Auth) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.LinkLogin(), http.StatusTemporaryRedirect)
}
