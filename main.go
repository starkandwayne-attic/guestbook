package main

import (
    "github.com/starkandwayne/guestbook/api"
    "github.com/starkandwayne/guestbook/api/database"
    "github.com/cloudfoundry-community/go-cfenv"
    "github.com/gorilla/mux"
    "flag"
    "log"
    "net/http"
    _"strings"
)

var DB database.PostgresDB

func main() {
    var postgresUri string
    var useCFEnv bool

    flag.StringVar(&postgresUri, "uri", "postgres://postgres@127.0.0.1:5432/guestbook?sslmode=disable", "Postgres URI")
    flag.BoolVar(&useCFEnv, "use_cfenv", false, "Use CF Env, overrides other settings")

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
    }

    DB = database.UsePostgresDB(postgresUri)
    //DB.EnsureStructure()
    r := mux.NewRouter()
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

    http.Handle("/", r)
    http.HandleFunc("/submit", SubmitHandler)
    http.HandleFunc("/posts/random", RandomPostHandler)
    http.ListenAndServe(":8080", nil)
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        api.PostSubmitHandler(w, r, &DB)
    }
}

func RandomPostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        api.GetRandomPostHandler(w, r, &DB)
    }
}
