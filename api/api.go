package api

import (
    "gitlab.swisscloud.io/Integration/qAPI/api/database"
    "gitlab.swisscloud.io/Integration/qAPI/api/config"
	"fmt"
    "io/ioutil"
    "encoding/json"
	"net/http"
    "github.com/codegangsta/martini"
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

func PostHandler(postTo string, w http.ResponseWriter, r *http.Request, params martini.Params, addFields map[string]interface{}) {
    var doc database.DBResult
    var res database.DBResult

    err := UnmarshalBody(r, &doc)
    if err != nil {
        ReturnError(err, w)
        return
    }

    doc.SetID(params["id"])

    for field, _ := range addFields {
        doc[field] = addFields[field]
    }

    err = DB.Insert(postTo, doc, &res)
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

