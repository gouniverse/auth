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
	return "exampleUserId", nil
}

func userLogout(username string) (err error) {
    // your code here
	return nil
}

func userStoreToken(token string, userID string) error {
    // your code here
	return nil
}

func userFindByToken(token string) (userID string, err error) {
    // your code here
	return "exampleUserId", nil
}
```

- Setup the auth settings

```golang
auth, err := auth.NewAuth(auth.Config{
	Endpoint:             "/",
	UrlRedirectOnSuccess: "/user/success",
	FuncUserLogin:        userLogin,
	FuncUserLogout:       userLogout,
	FuncUserRegister:     userRegister, // not providing it will disable registration
	FuncUserStoreToken:   userStoreToken,
	FuncUserFindByToken:  userFindByToken,
})
```

- Attach to router

```golang
mux := http.NewServeMux()
mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Index page. Login at: " + auth.LinkLogin()))
})

mux.HandleFunc("/auth/", auth.Router)

// Put your auth routes after the Auth middleware
mux.Handle("/user/dashboard", auth.AuthMiddleware(dashboardHandler("IN AUTHENTICATED DASHBOARD")))
```