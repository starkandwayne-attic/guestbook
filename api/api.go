package api

import (
    "github.com/starkandwayne/guestbook/api/database"
	"fmt"
    "io/ioutil"
    "encoding/json"
	"net/http"
    "strings"
    "errors"
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

func ReturnError(err error, w http.ResponseWriter, err_no int) {
    println("{\n\t\"error\": \"" + err.Error() + "\"\n}")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(err_no)
    if err_no == 500 {
        fmt.Fprint(w, "{\n\t\"error\": \"Egads!  Some supervillain blew up our application - please let our daring superheroes know so they can make things right.\"\n}")
    } else {
        fmt.Fprint(w, "{\n\t\"error\": \"" + err.Error()  + "\"\n}")
    }
}

func PostSubmitEmailHandler(w http.ResponseWriter, r *http.Request, db *database.PostgresDB) {
    var doc database.DBResult

    err := UnmarshalBody(r, &doc)
    if err != nil {
        ReturnError(err, w, 500)
        return
    }

    submitRequest := doc["submit"].(map[string]interface{})
    name := submitRequest["name"].(string)
    email := submitRequest["email"].(string)
    comment := ""
    _, hasComment := submitRequest["comment"]
    if hasComment {
        comment = submitRequest["comment"].(string)
    }

    post, err := db.SelectRandomPost(email)
    if err != nil {
        if err.Error() == "No posts found!" {
            noposts_err := errors.New("Wow - you read our entire blog?  Congratulations - you truly are a superhero!")
            ReturnError(noposts_err, w, 403)
            return
        }
        ReturnError(err, w, 500)
        return
    }
    post_id := post["id"].(int64)

    entry, err := db.InsertEntry(name, email, comment, post_id)

    if err != nil {
        ReturnError(err, w, 500)
        return
    }

    post["entry_id"] = entry["id"]
    post["success"] = "The first step is complete - now, click on the link to the blog post to obtain the code!"
    delete(post, "phrase")
    delete(post, "row_num")
    delete(post, "orig_id")


    w.Header().Set("Content-Type", "application/json")
    retval, err := json.MarshalIndent(post, "", "    ")
    if err != nil {
        ReturnError(err, w, 500)
        return
    }

    fmt.Fprint(w, string(retval))
}

func PostSubmitEntryHandler(w http.ResponseWriter, r *http.Request, db *database.PostgresDB) {
    var doc database.DBResult

    err := UnmarshalBody(r, &doc)
    if err != nil {
        ReturnError(err, w, 500)
        return
    }

    submitRequest := doc["submit"].(map[string]interface{})
    code := submitRequest["code"].(string)
    entry_id := int64(submitRequest["entry_id"].(float64))
    post_id := int64(submitRequest["post_id"].(float64))
    email := submitRequest["email"].(string)

    alreadyEntered, err := db.IsEmailAlreadySubmitted(email, post_id)
    if err != nil {
        ReturnError(err, w, 500)
        return
    }

    post, err := db.SelectPostById(post_id)
    if err != nil {
        ReturnError(err, w, 500)
        return
    }

    w.Header().Set("Content-Type", "application/json")

    if alreadyEntered {
        err = errors.New("What treachery is this?!  You have already entered the drawing.  Don't make us send our sidekick after you!")
        ReturnError(err, w, 403)
        return
    }

    if strings.ToUpper(post["phrase"].(string)) != strings.ToUpper(code) {
        err = errors.New("Oops! Looks like some supervillain slipped you the wrong code. Make sure you have the right blog post and try again - before it's too late!")
        ReturnError(err, w, 403)
        return
    }

    response, err := db.ValidateEntry(entry_id, post_id)

    if err != nil {
        ReturnError(err, w, 500)
        return
    }

    fmt.Fprint(w, "{\n\t\"success\": \"Nicely done, " + response["name"].(string)  + "!  You now have a new entry in the drawing.  If you want to keep entering, another link has been generated for you.  Keep reading our posts and collecting those codes so you too can be a superhero!\"\n}")
}
