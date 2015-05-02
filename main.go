package main

import (
    _ "github.com/starkandwayne/guestbook/api/database"
    "github.com/cloudfoundry-community/go-cfenv"
    "github.com/gorilla/mux"
    "flag"
    "fmt"
    "log"
    "net/http"
    _"strings"
)

func main() {
    var postgresUri string
    var useCFEnv bool

    flag.StringVar(&postgresUri, "uri", "postgres://postgres@127.0.0.1:5432/guestbook", "Postgres URI")
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

    //var DB database.PostgresDB = database.UsePostgresDB(postgresUri)
    //DB.EnsureStructure()
    r := mux.NewRouter()
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

    http.Handle("/", r)
    http.HandleFunc("/fnord", TestHandler)
    http.ListenAndServe(":8080", nil)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "asdf")
}
