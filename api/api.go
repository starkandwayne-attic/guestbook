package api

import (
    "github.com/starkandwayne/guestbook/api/database"
	"fmt"
    "io/ioutil"
    "encoding/json"
	"net/http"
    "strings"
)

func UnmarshalBody(r *http.Request, unmarshalTo database.Result) (error) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return err
    }
    err = json.Unmarshal(body, unmarshalTo)
    if err != nil {
        return err
    }
    return nil
}

func ReturnError(err error, w http.ResponseWriter) {
    println("{\n\t\"error\": \"" + err.Error() + "\"\n}")
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, "{\n\t\"error\": \"Egads!  Some supervillain blew up our application - please let our daring superheroes know so they can make things right.\"\n}")
}

func GetRandomPostHandler(w http.ResponseWriter, r *http.Request, db *database.PostgresDB) {
    var post database.DBResult

    post, err := db.SelectRandomPost()
    if err != nil {
        ReturnError(err, w)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    retval, err := json.MarshalIndent(post, "", "    ")
    if err != nil {
        ReturnError(err, w)
        return
    }
    fmt.Fprint(w, string(retval))
}

func PostSubmitHandler(w http.ResponseWriter, r *http.Request, db *database.PostgresDB, appName string) {
    var doc database.DBResult

    err := UnmarshalBody(r, &doc)
    if err != nil {
        ReturnError(err, w)
        return
    }

    submitRequest := doc["submit"].(map[string]interface{})
    code := submitRequest["code"].(string)
    name := submitRequest["name"].(string)
    email := submitRequest["email"].(string)
    comment := submitRequest["comment"].(string)
    post_id := int64(submitRequest["post_id"].(float64))

    post, err := db.SelectPostById(post_id)
    url := strings.ToLower(appName)
    alreadyEntered, err := db.IsEmailAlreadySubmitted(email, url)
    if err != nil {
        ReturnError(err, w)
        return
    }

    w.Header().Set("Content-Type", "application/json")

    if alreadyEntered {
        fmt.Fprint(w, "{\n\t\"error\": \"What treachery is this?!  You have already entered the drawing.  Don't make us send our sidekick after you!\"\n}")
        return
    }

    if strings.ToUpper(post["phrase"].(string)) != strings.ToUpper(code) {
        fmt.Fprint(w, "{\n\t\"error\": \"Oops! Looks like some supervillain slipped you the wrong code. Make sure you have the right blog post and try again - before it's too late!\"\n}")
        return
    }

    response, err := db.InsertEntry(name, email, comment, url)

    if err != nil {
        ReturnError(err, w)
        return
    }

    fmt.Fprint(w, "{\n\t\"success\": \"Nicely done, " + response["name"].(string)  + "!  You now have a new entry in the drawing.\"\n}")
    return
}
