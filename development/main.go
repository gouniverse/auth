package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/darkoatanasovski/htmltags"
	"github.com/gouniverse/auth/development/scribble"
	"github.com/jordan-wright/email"

	"github.com/gouniverse/auth"

	"github.com/gouniverse/utils"
)

var jsonStore *scribble.Driver

func emailSend(userID string, subject string, body string) error {
	emailSendTo("info@sinevia.com", []string{"info@sinevia.com"}, subject, body)
	return nil
}

func userLogin(username string, password string) (userID string, err error) {
	slug := utils.StrSlugify(username, rune('_'))
	var user map[string]string
	err = jsonStore.Read("users", slug, &user)
	if err != nil {
		return "not found err", err
	}
	log.Println(user)
	if user == nil {
		return "not found", errors.New("unable to find user")
	}

	if user["password"] == password {
		return username, nil
	}

	return "password mismatch", errors.New("password mismatch")
}

func userLogout(username string) error {
	return nil
}

func userFindByToken(token string) (userID string, err error) {
	slug := utils.StrSlugify(token, rune('_'))
	err = jsonStore.Read("tokens", slug, &userID)
	if err != nil {
		return "not found err", err
	}
	return userID, nil
}

func userFindByUsername(username string, firstName string, lastName string) (userID string, err error) {
	slug := utils.StrSlugify(username, rune('_'))
	var user map[string]string
	err = jsonStore.Read("users", slug, &user)
	if err != nil {
		return "not found err", err
	}

	if user == nil {
		return "not found", errors.New("unable to find user")
	}

	return user["id"], nil
}

func userPasswordChange(userID string, password string) error {
	user, err := userFindByID(userID)
	if err != nil {
		return err
	}

	user["password"] = password

	slug := utils.StrSlugify(user["username"], rune('_'))
	errSave := jsonStore.Write("users", slug, user)
	if errSave != nil {
		return errSave
	}

	jsonStore.Delete("users", slug)

	return nil
}

func userFindByID(userID string) (user map[string]string, err error) {
	users, errReadAll := jsonStore.ReadAll("users")
	if errReadAll != nil {
		return nil, errReadAll
	}

	for _, userJson := range users {
		json.Unmarshal(userJson, &user)
		log.Println(userID)
		log.Println(user)
		if user["id"] == userID {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

func userStoreToken(token string, userID string) error {
	slug := utils.StrSlugify(token, rune('_'))
	err := jsonStore.Write("tokens", slug, userID)
	if err != nil {
		return err
	}
	return nil
}

func userRegister(username string, password string, first_name string, last_name string) error {
	slug := utils.StrSlugify(username, rune('_'))
	err := jsonStore.Write("users", slug, map[string]string{
		"id":         utils.StrRandomFromGamma(16, "abcdef0123456789"),
		"username":   username,
		"password":   password,
		"first_name": first_name,
		"last_name":  last_name,
	})
	if err != nil {
		return err
	}
	return nil
}

func temporaryKeyGet(key string) (value string, err error) {
	slug := utils.StrSlugify(key, rune('_'))
	var record map[string]string
	err = jsonStore.Read("temp", slug, &record)
	if err != nil {
		return "", err
	}
	return record["value"], nil
}

func temporaryKeySet(key string, value string, expiresSeconds int) (err error) {
	slug := utils.StrSlugify(key, rune('_'))
	expiresAt := time.Now().Add(time.Duration(expiresSeconds))
	err = jsonStore.Write("temp", slug, map[string]string{
		"id":           utils.StrRandomFromGamma(16, "abcdef0123456789"),
		"value":        value,
		"expires":      utils.ToString(expiresSeconds),
		"expires_time": utils.ToString(expiresAt),
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	os.Remove(utils.Env("DB_DATABASE")) // remove database
	log.Println("1. Initializing environment variables...")
	utils.EnvInitialize()

	log.Println("2. Initializing database...")
	var err error
	jsonStore, err = scribble.New("temp", nil)
	if err != nil {
		log.Panic("Database is NIL: " + err.Error())
		return
	}

	auth, err := auth.NewAuth(auth.Config{
		Endpoint:               utils.Env("APP_URL") + "/auth",
		UrlRedirectOnSuccess:   "/user/dashboard",
		FuncEmailSend:          emailSend,
		FuncUserFindByToken:    userFindByToken,
		FuncUserFindByUsername: userFindByUsername,
		FuncUserLogin:          userLogin,
		FuncUserLogout:         userLogout,
		FuncUserPasswordChange: userPasswordChange,
		FuncUserRegister:       userRegister,
		FuncUserStoreToken:     userStoreToken,
		FuncTemporaryKeyGet:    temporaryKeyGet,
		FuncTemporaryKeySet:    temporaryKeySet,

		UseCookies: true,
	})

	if err != nil {
		log.Panicln(err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html>Index page. Login at: <a href='" + auth.LinkLogin() + "'>" + auth.LinkLogin() + "</a>"))
	})
	mux.HandleFunc("/auth/", auth.Router)
	mux.Handle("/user/dashboard", auth.AuthMiddleware(messageHandler("<html>User page. Logout at: <a href='"+auth.LinkLogout()+"'>"+auth.LinkLogout()+"</a>")))

	log.Println("4. Starting server on http://" + utils.Env("SERVER_HOST") + ":" + utils.Env("SERVER_PORT") + " ...")
	log.Println("URL: http://" + utils.Env("APP_URL") + " ...")

	srv := &http.Server{
		Handler: mux,
		Addr:    utils.Env("SERVER_HOST") + ":" + utils.Env("SERVER_PORT"),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Here"))
		next.ServeHTTP(w, r)
	})
}

func messageHandler(message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(message))
	})
}

// func mainDb(driverName string, dbHost string, dbPort string, dbName string, dbUser string, dbPass string) (*sql.DB, error) {
// 	var db *sql.DB
// 	var err error
// 	if driverName == "sqlite" {
// 		dsn := dbName + "?parseTime=true"
// 		db, err = sql.Open("sqlite3", dsn)
// 		// dsn := dbName
// 		// db, err = sql.Open("sqlite", dsn)
// 	}
// 	if driverName == "mysql" {
// 		dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
// 		db, err = sql.Open("mysql", dsn)
// 	}
// 	if driverName == "postgres" {
// 		dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Europe/London"
// 		db, err = sql.Open("postgres", dsn)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	if db == nil {
// 		return nil, errors.New("database for driver " + driverName + " could not be intialized")
// 	}

// 	return db, nil
// }

// EmailSend sends an email
func emailSendTo(from string, to []string, subject string, htmlMessage string) (bool, error) {
	//drvr := os.Getenv("MAIL_DRIVER")
	host := utils.Env("MAIL_HOST")
	port := utils.Env("MAIL_PORT")
	user := utils.Env("MAIL_USERNAME")
	pass := utils.Env("MAIL_PASSWORD")
	addr := host + ":" + port

	nodes, errStripped := htmltags.Strip(htmlMessage, []string{}, true)

	textMessage := ""

	if errStripped == nil {
		//nodes.Elements   //HTML nodes structure of type *html.Node
		textMessage = nodes.ToString() //returns stripped HTML string
	}

	e := email.NewEmail()
	e.From = from
	e.To = to
	e.Subject = subject
	e.Text = []byte(textMessage)
	e.HTML = []byte(htmlMessage)
	err := e.Send(addr, smtp.PlainAuth("", user, pass, host))

	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}
