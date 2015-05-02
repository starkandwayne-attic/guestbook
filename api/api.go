package api

import (
    "gitlab.swisscloud.io/Integration/qAPI/api/database"
    "gitlab.swisscloud.io/Integration/qAPI/api/config"
	"fmt"
    "io/ioutil"
    "encoding/json"
	"net/http"
    "strings"
)

var ENTITIES []string = make([]string,0)


var CONFIG config.Config = config.ParseConfig("config/config.json")
var DB database.PostgresDB = database.UsePostgresDB(CONFIG.DatabaseUri, CONFIG.DatabaseName)

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
    fmt.Fprint(w, "{\n\t\"error\": \"" + err.Error() + "\"\n}")
}

func GetHandler(getFrom string, w http.ResponseWriter, r *http.Request, params martini.Params, queryMap map[string]interface{}) {
    var doc database.DBResult

    err := DB.SelectFirst(getFrom, queryMap, &doc)
    if err != nil {
        ReturnError(err, w)
        return
    }

    retval, err := json.MarshalIndent(doc, "", "    ")
    if err != nil {
        ReturnError(err, w)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, string(retval))
}

// added by Steve
func GetLastHandler(getFrom string, w http.ResponseWriter, r *http.Request) {

    var doc database.DBResult

    err := DB.SelectFirst(getFrom, make(map[string]interface{}), &doc)
    if err != nil {
        ReturnError(err, w)
        return
    }

    retval, err := json.MarshalIndent(doc, "", "    ")
    if err != nil {
        ReturnError(err, w)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, string(retval))
}


func GetAllHandler(getFrom string, w http.ResponseWriter, r *http.Request, params martini.Params, queryMap map[string]interface{}) {
    docCollection := make([]database.DBResult,0)

    err := DB.SelectLatest(getFrom, queryMap, &docCollection)
    if err != nil {
        ReturnError(err, w)
        return
    }

    retval, _ := json.MarshalIndent(docCollection, "", "    ")
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, string(retval))
}

func PutHandler(putTo string, w http.ResponseWriter, r *http.Request, params martini.Params) {
    var doc database.DBResult
    var res database.DBResult

    err := UnmarshalBody(r, &doc)
    if err != nil {
        ReturnError(err, w)
        return
    }

    doc.SetID(params["id"])

    err = DB.Update(putTo, doc, &res)
    if err != nil {
        ReturnError(err, w)
        return
    }
    retval, err := json.MarshalIndent(res, "", "    ")
    if err != nil {
        ReturnError(err, w)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, string(retval))
}


func PostSubmitHandler(w http.ResponseWriter, r *http.Request, db *database.PostgresDB) {
    var doc database.DBResult
    var res database.DBResult

    err := UnmarshalBody(r, &doc)

    submitRequest := doc['submit']
    code := submitRequest['code']
    name := submitRequest['name']
    email := submitRequest['email']
    comment := submitRequest['comment']
    post_id := submitRequest['post_id']

    post := db.SelectPostById(post_id)
    alreadyEntered := db.IsEmailAlreadySubmitted(email)

    w.Header().Set("Content-Type", "application/json")

    if alreadyEntered {
        fmt.Fprint(w, "{\n\t\"error\": \"What treachery is this?!  You have already entered the drawing.  Don't make us send our sidekick after you!\"\n}")
        return
    }

    if strings.ToUpper(post['phrase'].(string)) != strings.ToUpper(code) {
        fmt.Fprint(w, "{\n\t\"error\": \"Oops! Looks like some supervillain slipped you the wrong code. Make sure you have the right blog post and try again - before it's too late!\"\n}")
        return
    }

    response, err := db.InsertEntry(name, email, comment)

    if err != nil {
        ReturnError(err, w)
        return
    }

    fmt.Fprint(w, "{\n\t\"success\": \"Nicely done, " + response['name']  + "!  You now have a new entry in the drawing.\"\n}")
    return
}






// added by steve
func DeleteHandler(deleteFrom string, w http.ResponseWriter, r *http.Request, params martini.Params) {
	var res database.DBResult
	
	id := params["id"]
	
	err := DB.Delete(deleteFrom, id, &res)
    if err != nil {
        ReturnError(err, w)
        return
    }
    retval, err := json.MarshalIndent(res, "", "    ")
    if err != nil {
        ReturnError(err, w)
        return
    }
    w.Header().Set("Access-Control-Allow-Origin", "127.0.0.1:3000")
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, string(retval))
}

func BuildMapFromParams(params martini.Params) (map[string]interface{}) {
    retval := make(map[string]interface{})

    for key, param := range params {
        retval[key] = param
    }
    return retval
}

