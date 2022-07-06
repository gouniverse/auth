package main

import (
	"auth"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gouniverse/cms"
	"github.com/gouniverse/utils"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func userRegister(username string, password string, first_name string, last_name string) error {
	return nil
}

func userLogin(username string, password string) (userID string, err error) {
	return "exampleUserId", nil
}

func userStoreToken(token string, userID string) error {
	return nil
}

func userFindByToken(token string) (userID string, err error) {
	return "exampleUserId", nil
}

func main() {
	log.Println("1. Initializing environment variables...")
	utils.EnvInitialize()

	log.Println("2. Initializing database...")
	db, err := mainDb(utils.Env("DB_DRIVER"), utils.Env("DB_HOST"), utils.Env("DB_PORT"), utils.Env("DB_DATABASE"), utils.Env("DB_USERNAME"), utils.Env("DB_PASSWORD"))

	if err != nil {
		log.Panic("Database is NIL: " + err.Error())
		return
	}

	if db == nil {
		log.Panic("Database is NIL")
		return
	}

	auth, err := auth.NewAuth(auth.Config{
		Endpoint:             "/auth",
		UrlRedirectOnSuccess: "/user/dashboard",
		FuncUserLogin:        userLogin,
		FuncUserRegister:     userRegister,
		FuncUserStoreToken:   userStoreToken,
		FuncUserFindByToken:  userFindByToken,
	})

	if err != nil {
		log.Panicln(err.Error())
	}

	log.Println("3. Initializing CMS...")
	myCms, err := cms.NewCms(cms.Config{
		DbInstance: db,
		// BlocksEnable:        true,
		// CacheAutomigrate:    true,
		// CacheEnable:         true,
		// EntitiesAutomigrate: true,
		// LogsAutomigrate:     true,
		// LogsEnable:          true,
		// MenusEnable:         true,
		// PagesEnable:         true,
		// SettingsAutomigate:  true,
		// SettingsEnable:      true,
		// SessionAutomigrate:  true,
		// SessionEnable:       true,
		// TemplatesEnable:     true,
		Prefix:           "cms",
		CustomEntityList: entityList(),
	})

	if err != nil {
		log.Panicln(err.Error())
	}

	log.Println("4. Starting server on http://" + utils.Env("SERVER_HOST") + ":" + utils.Env("SERVER_PORT") + " ...")
	log.Println("URL: http://" + utils.Env("APP_URL") + " ...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", auth.Router)
	mux.HandleFunc("/cms", myCms.Router)

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

func mainDb(driverName string, dbHost string, dbPort string, dbName string, dbUser string, dbPass string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	if driverName == "sqlite" {
		dsn := dbName
		db, err = sql.Open("sqlite", dsn)
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
