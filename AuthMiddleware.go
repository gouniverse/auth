package auth

import (
	"context"
	"net/http"

	"github.com/gouniverse/api"
)

// DEPRECATED use the Web or the API middleware instead
func (a Auth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authToken := authTokenRetrieve(r, a.useCookies)

		if authToken == "" {
			if a.useCookies {
				http.Redirect(w, r, a.LinkLogin(), http.StatusTemporaryRedirect)
				return
			}

			api.Respond(w, r, api.Unauthenticated("auth token is required"))
			return
		}

		userID, err := a.funcUserFindByAuthToken(authToken)

		if err != nil || userID == "" {
			if a.useCookies {
				http.Redirect(w, r, a.LinkLogin(), http.StatusTemporaryRedirect)
				return
			}

			api.Respond(w, r, api.Unauthenticated("auth token is required"))
			return
		}

		ctx := context.WithValue(r.Context(), AuthenticatedUserID{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
