package auth

import (
	"context"
	"net/http"

	"github.com/gouniverse/api"
)

// ApiAuthOrErrorMiddleware checks that an authentication token
// exists, and then finds the userID based on it. On success appends
// the user ID to the context. On failure it will return an
// unauthenticated JSON response.
func (a Auth) ApiAuthOrErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authToken := authTokenRetrieve(r, a.useCookies)

		if authToken == "" {
			api.Respond(w, r, api.Unauthenticated("auth token is required"))
			return
		}

		userID, err := a.funcUserFindByAuthToken(authToken)

		if err != nil {
			api.Respond(w, r, api.Unauthenticated("auth token is required"))
			return
		}

		if userID == "" {
			api.Respond(w, r, api.Unauthenticated("user id is required"))
			return
		}

		ctx := context.WithValue(r.Context(), AuthenticatedUserID{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
