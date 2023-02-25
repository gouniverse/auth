package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBlockRobotsMiddlewareShouldPassForHomeWithoutSlash(t *testing.T) {
	auth, err := NewPasswordlessAuth(ConfigPasswordless{
		Endpoint:                "/auth",
		UrlRedirectOnSuccess:    "/user",
		FuncTemporaryKeyGet:     func(key string) (value string, err error) { return "", nil },
		FuncTemporaryKeySet:     func(key, value string, expiresSeconds int) (err error) { return nil },
		FuncUserFindByAuthToken: func(sessionID string) (userID string, err error) { return "", nil },
		FuncUserFindByEmail:     func(email string) (userID string, err error) { return "", nil },
		FuncUserLogout:          func(userID string) (err error) { return nil },
		FuncUserStoreAuthToken:  func(sessionID, userID string) error { return nil },
		FuncEmailSend:           func(email, emailSubject, emailBody string) (err error) { return nil },
		UseCookies:              true,
		UseLocalStorage:         false,
	})

	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Cookie", "name=authtoken; count=x")
	req.AddCookie(&http.Cookie{Name: "authtoken", Value: "123456"})

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// DEBUG: t.Log("Passes as expected")
		value := r.Context().Value(AuthenticatedUserID{})

		if value == "" {
			t.Fatal("Response SHOULD BE empty but found ", value)
		}
	})

	rw := httptest.NewRecorder()
	handler := auth.WebAppendUserIdIfExistsMiddleware(testHandler)
	handler.ServeHTTP(rw, req)

	if rw.Body.String() != "" {
		t.Fatal("Response SHOULD BE empty but found ", rw.Body.String())
	}
}
