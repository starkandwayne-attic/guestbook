package database

import (
    _ "github.com/lib/pq"
)

func (db *PostgresDB) EnsureStructure() (error) {
    err := db.CreateEntriesTable()
    if err != nil {
        return err
    }
    err = db.CreatePostsTable()
    if err != nil {
        return err
    }
    return nil
}

func (db *PostgresDB) CreateEntriesTable() (error) {
    sqlCreateTable := "CREATE TABLE IF NOT EXISTS entries\n"
    sqlCreateTable += "(\n"
    sqlCreateTable += "     id bigserial NOT NULL,\n"
    sqlCreateTable += "     name text,\n"
    sqlCreateTable += "     email text,\n"
    sqlCreateTable += "     contact boolean,\n"
    sqlCreateTable += "     comment text,\n"
    sqlCreateTable += "     url text,\n"
    sqlCreateTable += "     CONSTRAINT entries_pkey PRIMARY KEY (id)\n"
    sqlCreateTable += ")\n"
    sqlCreateTable += "WITH (\n"
    sqlCreateTable += "     OIDS=FALSE\n"
    sqlCreateTable += ");\n"
    sqlCreateTable += "ALTER TABLE entries\n"
    sqlCreateTable += "OWNER TO postgres;"

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    _, err = session.Exec(sqlCreateTable)
    if err != nil {
        return err
    }

    return nil
}

func (db *PostgresDB) CreatePostsTable() (error) {
    sqlCreateTable := "CREATE TABLE IF NOT EXISTS posts\n"
    sqlCreateTable += "(\n"
    sqlCreateTable += "     id bigserial NOT NULL,\n"
    sqlCreateTable += "     url text,\n"
    sqlCreateTable += "     title text,\n"
    sqlCreateTable += "     phrase text,\n"
    sqlCreateTable += "     CONSTRAINT posts_pkey PRIMARY KEY (id)\n"
    sqlCreateTable += ")\n"
    sqlCreateTable += "WITH (\n"
    sqlCreateTable += "     OIDS=FALSE\n"
    sqlCreateTable += ");\n"
    sqlCreateTable += "ALTER TABLE posts\n"
    sqlCreateTable += "OWNER TO postgres;"

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    _, err = session.Exec(sqlCreateTable)
    if err != nil {
        return err
    }

    return nil
}

func (db *PostgresDB) TableExists(name string) (bool) {
    sqlCheckForTable := "SELECT 1\n"
    sqlCheckForTable += "FROM pg_catalog.pg_class c\n"
    sqlCheckForTable += "JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace\n"
    sqlCheckForTable += "WHERE n.nspname = 'public'\n"
    sqlCheckForTable += "AND c.relname = '" + name + "'"

    session, err := db.connect()
    if err != nil {
        return false
    }
    defer session.Close()

    results, err := session.Query(sqlCheckForTable)
    if err != nil {
        return false
    }

    tableList := db.parseResults(results)
    for _ = range tableList {
        return true
    }
    return false
}
