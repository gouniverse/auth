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

func userStoreToken(token string, userID string) error {
    // your code here
	return nil
}

func userFindByToken(token string) (userID string, err error) {
    // your code here
	return "exampleUserId", nil
}
```


```golang
auth, err := auth.NewAuth(auth.Config{
	Endpoint:             "/",
	UrlRedirectOnSuccess: "/user/success",
	FuncUserLogin:        userLogin,
	FuncUserRegister:     userRegister, // not providing this will disable registration
	FuncUserStoreToken:   userStoreToken,
	FuncUserFindByToken:  userFindByToken,
})
```