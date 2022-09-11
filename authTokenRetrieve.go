package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gouniverse/utils"
)

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
