package main

import (
    "github.com/starkandwayne/guestbook/api"
    "github.com/starkandwayne/guestbook/api/database"
    "github.com/cloudfoundry-community/go-cfenv"
    "github.com/gorilla/mux"
    "flag"
    "log"
    "net/http"
    "strconv"
)

var DB database.PostgresDB
var postgresUri string
var appName string
var useCFEnv bool
var port int

func main() {
    flag.StringVar(&postgresUri, "uri", "postgres://postgres@127.0.0.1:5432/guestbook?sslmode=disable", "Postgres URI")
    flag.BoolVar(&useCFEnv, "use_cfenv", false, "Use CF Env, overrides other settings")
    flag.StringVar(&appName, "app_name", "guestbook", "Application Name")
    flag.IntVar(&port, "port", 8080, "Application Port")

    flag.Parse()

    if useCFEnv {
        appEnv, enverr := cfenv.Current()
        if enverr != nil {
            log.Fatal("CF Env not found")
        }
        log.Printf("%#v\n", appEnv.Services)
        postgresService, err := appEnv.Services.WithName("guestbook-pg")
        if err == nil {
            postgresUri = postgresService.Credentials["uri"]
        } else {
            log.Fatal("Unable to get cf env")
        }
        //TODO:  Get application name from application_name in VCAP_APPLICATION
        //TODO:  Get port from port in VCAP_APPLICATION
    }

    DB = database.UsePostgresDB(postgresUri)
    DB.EnsureStructure()
    r := mux.NewRouter()
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

    http.Handle("/", r)
    http.HandleFunc("/submit", SubmitHandler)
    http.HandleFunc("/posts/random", RandomPostHandler)
    http.ListenAndServe(":" + strconv.Itoa(port), nil)
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        api.PostSubmitHandler(w, r, &DB, appName)
    }
}

func RandomPostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        api.GetRandomPostHandler(w, r, &DB)
    }
}
