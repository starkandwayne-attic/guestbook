package main

import (
    "github.com/starkandwayne/guestbook/api"
    "github.com/starkandwayne/guestbook/api/database"
    "github.com/cloudfoundry-community/go-cfenv"
    "github.com/gorilla/mux"
    "flag"
    "log"
    "net/http"
    "os"
    "strconv"
)

var DB database.PostgresDB
var postgresUri string
var appName string
var port int

func main() {
    flag.StringVar(&postgresUri, "uri", "postgres://postgres@127.0.0.1:5432/guestbook?sslmode=disable", "Postgres URI")
    flag.StringVar(&appName, "app_name", "guestbook", "Application Name")
    flag.IntVar(&port, "port", 8080, "Application Port")

    flag.Parse()

    println("Starting Guestbook Application...")

    if os.Getenv("VCAP_APPLICATION") != "" {
        println("Parsing Cloud Foundry environment variables...")
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
        appName = appEnv.Name
        port = appEnv.Port
    } else {
        println("Using Default or Parameter settings.")
    }

    DB = database.UsePostgresDB(postgresUri)
    err := DB.EnsureStructure()
    if err != nil {
        log.Fatal(err)
    }
    r := mux.NewRouter()
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

    http.Handle("/", r)
    http.HandleFunc("/submit/entry", SubmitEntryHandler)
    http.HandleFunc("/submit/email", SubmitEmailHandler)

    println("Listening on Port " + strconv.Itoa(port))
    http.ListenAndServe(":" + strconv.Itoa(port), nil)
}

func SubmitEntryHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        api.PostSubmitEntryHandler(w, r, &DB)
    }
}

func SubmitEmailHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        api.PostSubmitEmailHandler(w, r, &DB)
    }
}
