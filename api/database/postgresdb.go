package database

import (
    "database/sql"
    "encoding/json"
    "errors"
    "reflect"
    "strings"
    "strconv"
    "time"
    _ "github.com/coopernurse/gorp"
    _ "github.com/lib/pq"
)

type PostgresDB struct {
    DatabaseUri string
}

func UsePostgresDB(databaseUri string) (PostgresDB) {
    pgdb := PostgresDB{databaseUri}
    return pgdb
}

func (db *PostgresDB) connect() (*sql.DB, error) {
    conn, err := sql.Open("postgres", db.DatabaseUri)
    if err != nil {
        return &sql.DB{}, err
    }
    return conn, nil
}

func (db *PostgresDB) SelectLatest(selectFrom string, query map[string]interface{}, selectTo *[]DBResult) (error) {
    strParams, err := json.Marshal(query)
    if err != nil {
        return err
    }
    sqlQuery := "SELECT * FROM select_latest($1, $2)"
    sqlParams := make([]interface{},0)
    sqlParams = append(sqlParams, selectFrom, strParams)

    return db.doselect(sqlQuery, sqlParams, selectTo)
}

func (db *PostgresDB) SelectFirst(selectFrom string, query map[string]interface{}, selectTo *DBResult) (error) {
    selectResults := make([]DBResult,0)
    err := db.SelectLatest(selectFrom, query, &selectResults)
    if err != nil {
        return err
    }
    if len(selectResults) > 0 {
        *selectTo = selectResults[0]
    }
    return nil
}


func (db *PostgresDB) DoSelect(query string, queryParams []interface{}, selectTo *[]DBResult) (error) {
	return db.doselect(query, queryParams, selectTo)
}

func (db *PostgresDB) doselect(query string, queryParams []interface{}, selectTo *[]DBResult) (error) {
    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    results, err := session.Query(query, queryParams...)

    if err != nil {
        return err
    }

    parsedResults := db.parseResults(results)
    retval := make([]DBResult,0)
    for _, res := range parsedResults {
        var dbr DBResult
        err = json.Unmarshal([]byte(res["content"].(string)), &dbr)
        if err != nil {
            return err
        }
        retval = append(retval, dbr)
    }
    *selectTo = retval
    return nil
}

// added by steve
func (db *PostgresDB) Delete(deleteFrom string, id string, selectTo *DBResult) (error) {
    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

	deleteQuery := "DELETE FROM " + deleteFrom + " WHERE id = '" + id + "'"

    results, err := session.Query(deleteQuery)
    if err != nil {
        return err
    }

	parsedResults := db.parseResults(results)
    retval := make([]DBResult,0)
    for _, res := range parsedResults {
        var dbr DBResult
        err = json.Unmarshal([]byte(res["content"].(string)), &dbr)
        if err != nil {
            return err
        }
        retval = append(retval, dbr)
    }
	return nil
}

func (db *PostgresDB) Insert(insertInto string, insertObject DBResult, selectTo *DBResult) (error) {
    now := time.Now().UTC()
    insertObject.SetCreated(now)
    insertObject.SetUpdated(now)

    if insertObject.GetID() != "" {
        queryMap := make(map[string]interface{})
        queryMap["id"] = insertObject.GetID()

        var existingDoc DBResult
        err := db.SelectFirst(insertInto, queryMap, &existingDoc)
        if existingDoc != nil && err == nil {
            insertObject.SetCreated(existingDoc.GetCreated())
        }
    }

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    jsonToInsert, err := json.Marshal(insertObject)
    if err != nil {
        return err
    }

    insertParams := make([]interface{},4)
    insertParams[0] = insertObject.GetID()
    insertParams[1] = insertObject.GetCreated().Format(time.RFC3339)
    insertParams[2] = insertObject.GetUpdated().Format(time.RFC3339)
    insertParams[3] = string(jsonToInsert)

    insertQuery := "INSERT INTO " + insertInto + " (id,created,updated,content)\n VALUES("
    insertQuery += "$1,$2,$3,$4"
    insertQuery += ") RETURNING *"

    selectResults := make([]DBResult,0)
    err = db.doselect(insertQuery, insertParams, &selectResults)
    if err != nil {
        return err
    }
    if len(selectResults) > 0 {
        *selectTo = selectResults[0]
    }
    return nil
}

func (db *PostgresDB) Update(updateInto string, updateObject DBResult, selectTo *DBResult) (error) {
    if updateObject.GetID() == "" {
        err := errors.New("Cannot use Update without an ID.")
        return err
    }
    now := time.Now().UTC()
    updateObject.SetUpdated(now)

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    jsonToUpdate, err := json.Marshal(updateObject)
    if err != nil {
        return err
    }

    updateParams := make([]interface{},3)
    updateParams[0] = jsonToUpdate
    updateParams[1] = updateObject.GetUpdated().Format(time.RFC3339)
    updateParams[2] = string(updateObject.GetID())


    updateQuery := "UPDATE " + updateInto  + " SET content = jsonb_merge(content, $1),\n"
    updateQuery += "updated = $2\n"
    updateQuery += "WHERE id = $3\n"
    updateQuery += "RETURNING *"

    selectResults := make([]DBResult,0)
    err = db.doselect(updateQuery, updateParams, &selectResults)
    if err != nil {
        return err
    }
    if len(selectResults) > 0 {
        *selectTo = selectResults[0]
    }
    return nil
}


func (db *PostgresDB) buildWhereClauseFromMap(query map[string]interface{}, startParamsFrom int) (string, []interface{}) {
    sqlWhere := ""
    sqlParams := make([]interface{},0)
    currentParam := startParamsFrom

    for key, param := range query {
        if sqlWhere == "" {
            sqlWhere += "WHERE "
        }
        if sqlWhere != "WHERE " {
            sqlWhere += " AND "
        }
        operator := "="
        var value interface{} = param.(string)
        if param == nil {
            operator = "IS"
            value = nil
        }
        if strings.Contains(key, ".") {
            col := strings.Split(key, ".")[1]
            sqlWhere += "cast(content->>$" + strconv.Itoa(currentParam) + " as text) " + operator + " $" + strconv.Itoa(currentParam+1)
            sqlParams = append(sqlParams, col)
            sqlParams = append(sqlParams, value)
            currentParam+=2
        } else {
            sqlWhere += "\"$" + strconv.Itoa(currentParam) + "\" " + operator + " $" + strconv.Itoa(currentParam+1)
            sqlParams = append(sqlParams, key)
            sqlParams = append(sqlParams, value)
            currentParam+=2
        }
    }
    if sqlWhere != "" {
        sqlWhere += "\n"
    }

    return sqlWhere, sqlParams
}

func (db *PostgresDB) parseResults(r *sql.Rows) []DBResult {
    cols, _ := r.Columns()

    var newMapSlice = make([]DBResult,0)

    var counter int = 0
    for r.Next() {
        counter = counter + 1
        var newRow = make(DBResult)
        var scanargs = make([]interface{}, len(cols))
        var scanvals = make([]interface{}, len(cols))

        for i := range scanargs {
            scanargs[i] = &scanvals[i]
        }
        r.Scan(scanargs...)

        for i, columnname := range cols {
            if scanvals[i] != nil {
                if reflect.TypeOf(scanvals[i]).String() == "[]uint8" {
                    newRow[columnname] = string(scanvals[i].([]byte))
                } else {
                    newRow[columnname] = scanvals[i]
                }
            }
        }
        newMapSlice = append(newMapSlice, newRow)
    }
    return newMapSlice
}
