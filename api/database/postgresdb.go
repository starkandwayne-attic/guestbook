package database

import (
    "database/sql"
    "errors"
    "reflect"
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

func (db *PostgresDB) SelectPostById(post_id int64) (DBResult, error) {
    docCollection := make([]DBResult,0)

    sqlQuery := "SELECT id, url, title, phrase FROM posts WHERE id = $1"

    sqlParams := make([]interface{},1)
    sqlParams[0] = post_id

    err := db.DoSelect(sqlQuery, sqlParams, &docCollection)
    if err != nil {
        return DBResult{}, err
    }
    if len(docCollection) > 0 {
        return docCollection[0], nil
    }
    return DBResult{}, errors.New("No posts found with that id!")
}

func (db *PostgresDB) IsEmailAlreadySubmitted(email string, url string) (bool, error) {
    docCollection := make([]DBResult,0)

    sqlQuery := "SELECT id FROM entries WHERE email = $1 AND url = $2"

    sqlParams := make([]interface{},2)
    sqlParams[0] = email
    sqlParams[1] = url

    err := db.DoSelect(sqlQuery, sqlParams, &docCollection)
    if err != nil {
        return true, err
    }
    if len(docCollection) > 0 {
        return true, nil
    }
    return false, nil
}

func (db *PostgresDB) SelectRandomPost() (DBResult, error) {
    docCollection := make([]DBResult,0)

    sqlQuery := "SELECT id, url, title FROM posts p, "
    sqlQuery += "(SELECT cast(trunc(random() * count(*) + 1) as bigint) as post_id FROM posts) rand_id "
    sqlQuery += "WHERE p.id = rand_id.post_id;"

    sqlParams := make([]interface{},0)

    err := db.DoSelect(sqlQuery, sqlParams, &docCollection)
    if err != nil {
        return DBResult{}, err
    }
    if len(docCollection) > 0 {
        return docCollection[0], nil
    }
    return DBResult{}, errors.New("No posts found!")
}

func (db *PostgresDB) InsertEntry(name string, email string, comment string, url string) (DBResult, error) {
    docCollection := make([]DBResult,0)

    sqlQuery := "INSERT INTO entries (name, email, comment, url) VALUES($1,$2,$3,$4) RETURNING *"

    sqlParams := make([]interface{},4)
    sqlParams[0] = name
    sqlParams[1] = email
    sqlParams[2] = comment
    sqlParams[3] = url

    err := db.DoSelect(sqlQuery, sqlParams, &docCollection)
    if err != nil {
        return DBResult{}, err
    }
    if len(docCollection) > 0 {
        return docCollection[0], nil
    }
    return DBResult{}, errors.New("Entry not inserted!")
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
    //retval := make([]DBResult,0)
    //for _, res := range parsedResults {
    //    var dbr DBResult
    //    err = json.Unmarshal([]byte(res["content"].(string)), &dbr)
    //    if err != nil {
    //        return err
    //    }
    //    retval = append(retval, dbr)
    //}
    *selectTo = parsedResults
    return nil
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
