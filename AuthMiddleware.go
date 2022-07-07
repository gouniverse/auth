package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

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

		userID, err := a.funcUserFindByToken(authToken)

		if err != nil {
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

// authTokenRetrieve retrieves the auth token from the request
// Several attempts are made:
//  1. Authorization header (aka Bearer token)
//  2. Request param "api_key"
//  3. Request param "token"
func authTokenRetrieve(r *http.Request, useCookies bool) string {
	if useCookies {
		authTokenFromCookie, err := r.Cookie("authtoken")
		if err != nil {
			log.Println(err.Error())
			return ""
		}
		return authTokenFromCookie.Value
	}

	authorizationHeader := r.Header.Get("Authorization")
	authTokenFromBearerToken := bearerToken(authorizationHeader)

	if authTokenFromBearerToken != "" {
		return authTokenFromBearerToken
	}

	apiKeyFromRequest := utils.Req(r, "api_key", "")

	if apiKeyFromRequest != "" {
		return apiKeyFromRequest
	}

	tokenFromRequest := utils.Req(r, "token", "")

	if tokenFromRequest != "" {
		return tokenFromRequest
	}

	return ""
}

// BearerAuthHeader validates incoming `r.Header.Get("Authorization")` header
// and returns token otherwise an empty string.
func bearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, "Bearer")
	if len(parts) != 2 {
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if len(token) < 1 {
		return ""
	}

	return token
}
