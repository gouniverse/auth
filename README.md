# Auth
Authentication library for Golang

## Usage

- Implement your functions

```golang
func userRegister(username string, password string, first_name string, last_name string) error {
    // your code here
	return nil
}

func userLogin(username string, password string) (userID string, err error) {
    // your code here
	return "yourUserId", nil
}

func userLogout(username string) (err error) {
    // your code here (remove token from session or cache store)
	return nil
}

func userStoreToken(token string, userID string) error {
    // your code here (store in session or cache store with desired timeout)
	return nil
}

func userFindByToken(token string) (userID string, err error) {
    // your code here
	return "yourUserId", nil
}
```

- Setup the auth settings

```golang
auth, err := auth.NewAuth(auth.Config{
	EnableRegistration:   true,
	Endpoint:                "/",
	UrlRedirectOnSuccess:   "http://localhost/user/dashboard",
	FuncUserFindByToken:     userFindByToken,
	FuncUserFindByUsername:  userFindByUsername,
	FuncUserLogin:           userLogin,
	FuncUserLogout:          userLogout,
	FuncUserRegister:        userRegister, // optional, required only if registration is enabled
	FuncUserStoreToken:      userStoreToken,
	FuncEmailSend:           emailSend,
	FuncTemporaryKeyGet:     tempKeyGet,
	FuncTemporaryKeySet:     tempKeySet,
})
```

- Attach to router

```golang
mux := http.NewServeMux()
mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Index page. Login at: " + auth.LinkLogin()))
})

mux.HandleFunc("/auth/", auth.Router)
```

- Used the AuthMiddleware to protect the authenticated routes

```golang
// Put your auth routes after the Auth middleware
mux.Handle("/user/dashboard", auth.AuthMiddleware(dashboardHandler("IN AUTHENTICATED DASHBOARD")))
```


## Other Noteable Auth Projects

- https://github.com/authorizerdev/authorizer
- https://github.com/markbates/goth
