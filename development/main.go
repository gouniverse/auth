package main

import (
	"auth"
	"auth/development/scribble"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gouniverse/cms"
	"github.com/gouniverse/utils"
	// _ "github.com/go-sql-driver/mysql"
	// _ "github.com/lib/pq"
	// _ "github.com/mattn/go-sqlite3"
	// _ "modernc.org/sqlite"
)

var myCms *cms.Cms
var jsonStore *scribble.Driver

func userRegister(username string, password string, first_name string, last_name string) error {
	slug := utils.Slugify(username, rune('_'))
	err := jsonStore.Write("users", slug, map[string]string{
		"username":   username,
		"password":   password,
		"first_name": first_name,
		"last_name":  last_name,
	})
	if err != nil {
		return err
	}
	// isOk, err := myCms.CacheStore.SetJSON(username, map[string]string{
	// 	"username":   username,
	// 	"password":   password,
	// 	"first_name": first_name,
	// 	"last_name":  last_name,
	// }, 600)
	// if err != nil {
	// 	return err
	// }
	// if !isOk {
	// 	return errors.New("unable to register")
	// }
	return nil
}

func userLogin(username string, password string) (userID string, err error) {
	slug := utils.Slugify(username, rune('_'))
	var user map[string]string
	err = jsonStore.Read("users", slug, &user)
	if err != nil {
		return "not found err", err
	}
	log.Println(user)
	// user, err := myCms.CacheStore.GetJSON(username, nil)
	// if err != nil {
	// 	return "not found err", err
	// }
	if user == nil {
		return "not found", errors.New("unable to find user")
	}

	if user["password"] == password {
		return username, nil
	}

	return "password mismatch", errors.New("password mismatch")
}

func userStoreToken(token string, userID string) error {
	slug := utils.Slugify(token, rune('_'))
	err := jsonStore.Write("tokens", slug, userID)
	if err != nil {
		return err
	}
	// isOk, err := myCms.SessionStore.Set(token, userID, 600)
	// if err != nil {
	// 	return err
	// }
	// if !isOk {
	// 	return errors.New("unable to store token")
	// }
	return nil
}

func userFindByToken(token string) (userID string, err error) {
	slug := utils.Slugify(token, rune('_'))
	err = jsonStore.Read("tokens", slug, &userID)
	if err != nil {
		return "not found err", err
	}
	// userID = myCms.SessionStore.Get(token, "")
	return userID, nil
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

	// db, err := mainDb(utils.Env("DB_DRIVER"), utils.Env("DB_HOST"), utils.Env("DB_PORT"), utils.Env("DB_DATABASE"), utils.Env("DB_USERNAME"), utils.Env("DB_PASSWORD"))
	// defer db.Close()

	// if err != nil {
	// 	log.Panic("Database is NIL: " + err.Error())
	// 	return
	// }

	// if db == nil {
	// 	log.Panic("Database is NIL")
	// 	return
	// }

	auth, err := auth.NewAuth(auth.Config{
		Endpoint:             "/auth",
		UrlRedirectOnSuccess: "/user/dashboard",
		FuncUserLogin:        userLogin,
		FuncUserRegister:     userRegister,
		FuncUserStoreToken:   userStoreToken,
		FuncUserFindByToken:  userFindByToken,
		UseCookies:           true,
	})

	if err != nil {
		log.Panicln(err.Error())
	}

	// log.Println("3. Initializing CMS...")
	// myCms, err = cms.NewCms(cms.Config{
	// 	DbInstance: db,
	// 	// BlocksEnable:        true,
	// 	CacheAutomigrate: true,
	// 	CacheEnable:      true,
	// 	// EntitiesAutomigrate: true,
	// 	// LogsAutomigrate:     true,
	// 	// LogsEnable:          true,
	// 	// MenusEnable:         true,
	// 	// PagesEnable:         true,
	// 	// SettingsAutomigate:  true,
	// 	// SettingsEnable:      true,
	// 	SessionAutomigrate: true,
	// 	SessionEnable:      true,
	// 	// TemplatesEnable:     true,
	// 	Prefix:           "cms",
	// 	CustomEntityList: entityList(),
	// })

	// if err != nil {
	// 	log.Panicln(err.Error())
	// }

	log.Println("4. Starting server on http://" + utils.Env("SERVER_HOST") + ":" + utils.Env("SERVER_PORT") + " ...")
	log.Println("URL: http://" + utils.Env("APP_URL") + " ...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Index page. Login at: " + auth.LinkLogin()))
	})
	mux.HandleFunc("/auth/", auth.Router)
	// mux.HandleFunc("/cms", myCms.Router)
	// mux.HandleFunc("/user/dashboard", func(w http.ResponseWriter, r *http.Request) {
	// 	auth.AuthMiddleware(Middleware)
	// 	w.Write([]byte("User dashboard"))
	// })
	mux.Handle("/user/dashboard", auth.AuthMiddleware(messageHandler("INSIDE MIDDLEWARE")))

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

func mainDb(driverName string, dbHost string, dbPort string, dbName string, dbUser string, dbPass string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	if driverName == "sqlite" {
		dsn := dbName + "?parseTime=true"
		db, err = sql.Open("sqlite3", dsn)
		// dsn := dbName
		// db, err = sql.Open("sqlite", dsn)
	}
	if driverName == "mysql" {
		dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = sql.Open("mysql", dsn)
	}
	if driverName == "postgres" {
		dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Europe/London"
		db, err = sql.Open("postgres", dsn)
	}
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("database for driver " + driverName + " could not be intialized")
	}

	return db, nil
}

func entityList() []cms.CustomEntityStructure {
	list := []cms.CustomEntityStructure{}
	return list
}
