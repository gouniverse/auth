# Auth

![tests](https://github.com/gouniverse/auth/workflows/tests/badge.svg)

Authentication library for Golang

## Usage

- Implement your functions

```golang
// userRegister registers the user
//
// save the user to your databases
// note that the username can be an email (if you wish)
//
func userRegister(username string, password string, first_name string, last_name string) error {
    // your code here
	return nil
}

// userLogin logs the user in
//
// find the user by the specified username and password, 
// note that the username can be an email (if you wish)
//
func userLogin(username string, password string) (userID string, err error) {
    // your code here
	return "yourUserId", nil
}

// userLogout logs the user out
//
// remove the auth token from wherever you have stored it (i.e. session store or the cache store)
//
func userLogout(username string) (err error) {
    // your code here (remove token from session or cache store)
	return nil
}

// userStoreAuthToken stores the auth token with the provided user ID
//
// save the auth token to your selected store it (i.e. session store or the cache store)
// make sure you set an expiration time (i.e. 2 hours)
//
func userStoreAuthToken(token string, userID string) error {
    // your code here (store in session or cache store with desired timeout)
	return nil
}

// userFindByAuthToken find the user by the provided token, and returns the user ID
//
// retrieve the userID from your selected store  (i.e. session store or the cache store)
//
func userFindByAuthToken(token string) (userID string, err error) {
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
	FuncUserFindByAuthToken: userFindByAuthToken,
	FuncUserFindByUsername:  userFindByUsername,
	FuncUserLogin:           userLogin,
	FuncUserLogout:          userLogout,
	FuncUserRegister:        userRegister, // optional, required only if registration is enabled
	FuncUserStoreAuthToken:  userStoreAuthToken,
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
