package auth

import (
	"net/http"

	"github.com/gouniverse/utils"
)

// authTokenRetrieve retrieves the auth token from the request
// Several attempts are made:
//  1. From cookie
//  2. Authorization header (aka Bearer token)
//  3. Request param "api_key"
//  4. Request param "token"
func AuthTokenRetrieve(r *http.Request, useCookies bool) string {
	// 1. Token from cookie
	if useCookies {
		return AuthCookieGet(r)
	}

	// 2. Bearer token
	authTokenFromBearerToken := BearerTokenFromHeader(r.Header.Get("Authorization"))

	if authTokenFromBearerToken != "" {
		return authTokenFromBearerToken
	}

	// 3. API key
	apiKeyFromRequest := utils.Req(r, "api_key", "")

	if apiKeyFromRequest != "" {
		return apiKeyFromRequest
	}

	// 4. Token
	tokenFromRequest := utils.Req(r, "token", "")

	if tokenFromRequest != "" {
		return tokenFromRequest
	}

	return ""
}
