package auth

import (
	"context"
	"net/http"
)

func (a Auth) WebJustAppendUserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authToken := authTokenRetrieve(r, a.useCookies)

		if authToken != "" {
			userID, err := a.funcUserFindByAuthToken(authToken)

			if err != nil {
				http.Redirect(w, r, a.LinkLogin(), http.StatusTemporaryRedirect)
			}

			ctx := context.WithValue(r.Context(), AuthenticatedUserID{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}
