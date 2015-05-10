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
    err = db.CreatePhrasesTable()
    if err != nil {
        return err
    }
    err = db.Create_select_remaining_posts_for_email_Function()
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
    sqlCreateTable += "     comment text,\n"
    sqlCreateTable += "     post_id int,\n"
    sqlCreateTable += "     entered bool DEFAULT FALSE,\n"
    sqlCreateTable += "     CONSTRAINT entries_pkey PRIMARY KEY (id)\n"
    sqlCreateTable += ")\n"
    sqlCreateTable += "WITH (\n"
    sqlCreateTable += "     OIDS=FALSE\n"
    sqlCreateTable += ");\n"

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
    sqlCreateTable += "     orig_id integer NOT NULL,\n"
    sqlCreateTable += "     url text,\n"
    sqlCreateTable += "     title text,\n"
    sqlCreateTable += "     phrase text,\n"
    sqlCreateTable += "     CONSTRAINT posts_pkey PRIMARY KEY (id)\n"
    sqlCreateTable += ")\n"
    sqlCreateTable += "WITH (\n"
    sqlCreateTable += "     OIDS=FALSE\n"
    sqlCreateTable += ");\n"

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

func (db *PostgresDB) CreatePhrasesTable() (error) {
    sqlCreateTable := "CREATE TABLE IF NOT EXISTS phrases\n"
    sqlCreateTable += "(\n"
    sqlCreateTable += "     id bigserial NOT NULL,\n"
    sqlCreateTable += "     phrase text,\n"
    sqlCreateTable += "     CONSTRAINT phrases_pkey PRIMARY KEY (id)\n"
    sqlCreateTable += ")\n"
    sqlCreateTable += "WITH (\n"
    sqlCreateTable += "     OIDS=FALSE\n"
    sqlCreateTable += ");\n"

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

func (db *PostgresDB) Create_select_remaining_posts_for_email_Function() (error) {
    sqlCreateFunc := "CREATE OR REPLACE FUNCTION select_remaining_posts_for_email(IN email_address text)\n"
    sqlCreateFunc += "  RETURNS TABLE(row_num bigint, id bigint, url text, title text, phrase text) AS\n"
    sqlCreateFunc += "$BODY$\n"
    sqlCreateFunc += "BEGIN\n"
    sqlCreateFunc += "    RETURN QUERY SELECT row_number() over (order by p.id) AS row_num, p.id, p.url, p.title, p.phrase FROM posts p\n"
    sqlCreateFunc += "    WHERE NOT EXISTS(SELECT entries.id FROM entries where entries.post_id = p.id AND entries.email = email_address);\n"
    sqlCreateFunc += "END; $BODY$\n"
    sqlCreateFunc += "  LANGUAGE plpgsql VOLATILE\n"
    sqlCreateFunc += "  COST 100\n"
    sqlCreateFunc += "  ROWS 1000;\n"
    sqlCreateFunc += "ALTER FUNCTION select_remaining_posts_for_email(text)\n"
    sqlCreateFunc += "  OWNER TO postgres;"

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    _, err = session.Exec(sqlCreateFunc)
    if err != nil {
        return err
    }

    return nil
}
