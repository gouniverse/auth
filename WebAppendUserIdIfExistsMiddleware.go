package auth

import (
	"context"
	"log"
	"net/http"
)

func (a Auth) WebAppendUserIdIfExistsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authToken := authTokenRetrieve(r, a.useCookies)

		if authToken != "" {
			userID, err := a.funcUserFindByAuthToken(authToken)

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			log.Println("USERID", userID)

			ctx := context.WithValue(r.Context(), AuthenticatedUserID{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}
