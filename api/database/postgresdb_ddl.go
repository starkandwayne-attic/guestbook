package database

import (
    _ "github.com/coopernurse/gorp"
    _ "github.com/lib/pq"
    "strings"
)

func (db *PostgresDB) EnsureStructure(tables []string) (error) {
    db.Create_jsonb_object_set_key_Function()
    db.Create_jsonb_merge_Function()
    db.Create_uuidgen_Function()
    db.Create_generate_where_from_parameters_Function()
    db.Create_select_latest_Function()
    db.Create_archive_on_delete_Function()

    for _, table := range tables {
        if (db.TableExists(table) == false || db.TableExists(table + "_archive") == false) {
            err := db.CreateEntityTable(table)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func (db *PostgresDB) CreateEntityTable(name string) (error) {
    sqlCreateTable := "CREATE TABLE IF NOT EXISTS " + name + " (\n"
    sqlCreateTable += "record_id bigserial,\n"
    sqlCreateTable += "id text,\n"
    sqlCreateTable += "created timestamptz,\n"
    sqlCreateTable += "updated timestamptz,\n"
    sqlCreateTable += "content jsonb,\n"
    sqlCreateTable += "deprecated boolean DEFAULT False,\n"
    sqlCreateTable += "deleted timestamptz DEFAULT NULL,\n"
    sqlCreateTable += "CONSTRAINT pk_" + name + " PRIMARY KEY(record_id)"
    sqlCreateTable += ");"

    sqlCreateArchiveTable := "CREATE TABLE IF NOT EXISTS " + name + "_archive (\n"
    sqlCreateArchiveTable += "deleted timestamptz DEFAULT NOW(),\n"
    sqlCreateArchiveTable += "CONSTRAINT pk_" + name + "_archive PRIMARY KEY(record_id)"
    sqlCreateArchiveTable += ") INHERITS (" + name + ");"

    sqlCreateIndex_Created := "CREATE INDEX idx_" + name + "_created ON " + name + " (created)"
    sqlCreateIndex_Updated := "CREATE INDEX idx_" + name + "_updated ON " + name + " (updated)"

    sqlCreateArchiveIndex_Created := "CREATE INDEX idx_" + name + "_archive__created ON " + name + "_archive (created)"
    sqlCreateArchiveIndex_Updated := "CREATE INDEX idx_" + name + "_archive_updated ON " + name + "_archive (updated)"

    sqlCreateInsertTrigger := "CREATE TRIGGER uuidgen BEFORE INSERT ON " + name + " FOR EACH ROW EXECUTE PROCEDURE uuidgen();"
    sqlCreateDeleteTrigger := "CREATE TRIGGER archive_on_delete BEFORE DELETE ON " + name + " FOR EACH ROW EXECUTE PROCEDURE archive_on_delete();"

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    _, err = session.Exec(sqlCreateTable)
    if err != nil {
        return err
    }

    _, err = session.Exec(sqlCreateIndex_Created)
    _, err = session.Exec(sqlCreateIndex_Updated)
    _, err = session.Exec(sqlCreateInsertTrigger)
    _, err = session.Exec(sqlCreateDeleteTrigger)
    if (err != nil && strings.Contains(err.Error(), "already exists") == false) {
        return err
    }

    _, err = session.Exec(sqlCreateArchiveTable)
    if err != nil {
        return err
    }

    _, err = session.Exec(sqlCreateArchiveIndex_Created)
    _, err = session.Exec(sqlCreateArchiveIndex_Updated)

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

func (db *PostgresDB) Create_archive_on_delete_Function() (error) {
    sqlCreateFunction := "CREATE OR REPLACE FUNCTION archive_on_delete() RETURNS TRIGGER AS \n"
    sqlCreateFunction += "$BODY$\n"
    sqlCreateFunction += "BEGIN\n"
    sqlCreateFunction += "  EXECUTE 'INSERT INTO ' || quote_ident(TG_TABLE_NAME || '_archive') || "
    sqlCreateFunction += " ' SELECT ($1).record_id, ($1).id, ($1).created, ($1).updated, ($1).content, ($1).deprecated, NOW()'\n"
    sqlCreateFunction += "  USING OLD;\n"
    sqlCreateFunction += "  RETURN OLD;\n"
    sqlCreateFunction += "END\n"
    sqlCreateFunction += "$BODY$\n"
    sqlCreateFunction += "  LANGUAGE plpgsql VOLATILE;\n"

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    _, err = session.Exec(sqlCreateFunction)
    if err != nil {
        return err
    }
    return nil

}

func (db *PostgresDB) Create_uuidgen_Function() (error) {
    sqlCreateFunction := "CREATE OR REPLACE FUNCTION uuidgen() RETURNS trigger AS '\n"
    sqlCreateFunction += "DECLARE\n"
    sqlCreateFunction += "  new_uuid text;\n"
    sqlCreateFunction += "BEGIN\n"
    sqlCreateFunction += "  IF NEW.id = '''' THEN\n"
    sqlCreateFunction += "        new_uuid := uuid_in(md5(now()::text)::cstring);\n"
    sqlCreateFunction += "        NEW.id := new_uuid;\n"
    sqlCreateFunction += "        NEW.content := jsonb_object_set_key(NEW.content, ''id'', new_uuid);\n"
    sqlCreateFunction += "  END IF;\n"
    sqlCreateFunction += "  RETURN NEW;\n"
    sqlCreateFunction += "END;\n"
    sqlCreateFunction += "' LANGUAGE plpgsql;"

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    _, err = session.Exec(sqlCreateFunction)
    if err != nil {
        return err
    }
    return nil
}    

func (db *PostgresDB) Create_jsonb_object_set_key_Function() (error) {
    sqlCreateFunction := "CREATE OR REPLACE FUNCTION \"jsonb_object_set_key\"(\n"
    sqlCreateFunction += "  \"json\"          jsonb,\n"
    sqlCreateFunction += "  \"key_to_set\"    TEXT,\n"
    sqlCreateFunction += "  \"value_to_set\"  anyelement\n"
    sqlCreateFunction += ")\n"
    sqlCreateFunction += "  RETURNS jsonb\n"
    sqlCreateFunction += "  LANGUAGE sql\n"
    sqlCreateFunction += "  IMMUTABLE\n"
    sqlCreateFunction += "  STRICT\n"
    sqlCreateFunction += "AS $function$\n"
    sqlCreateFunction += "SELECT COALESCE(\n"
    sqlCreateFunction += "  (SELECT ('{' || string_agg(to_json(\"key\")::jsonb || ':' || \"value\", ',') || '}')\n"
    sqlCreateFunction += "     FROM (SELECT *\n"
    sqlCreateFunction += "           FROM jsonb_each(\"json\")\n"
    sqlCreateFunction += "           WHERE \"key\" <> \"key_to_set\"\n"
    sqlCreateFunction += "           UNION ALL\n"
    sqlCreateFunction += "           SELECT \"key_to_set\", to_json(\"value_to_set\")::jsonb) AS \"fields\"),\n"
    sqlCreateFunction += "  '{}'\n"
    sqlCreateFunction += ")::jsonb\n"
    sqlCreateFunction += "$function$;\n"

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    _, err = session.Exec(sqlCreateFunction)
    if err != nil {
        return err
    }
    return nil
}

func (db *PostgresDB) Create_jsonb_merge_Function() (error) {
    sqlCreateFunction := "CREATE OR REPLACE FUNCTION public.jsonb_merge(data jsonb, merge_data jsonb)\n"
    sqlCreateFunction += "   RETURNS jsonb\n"
    sqlCreateFunction += "LANGUAGE sql\n"
    sqlCreateFunction += "AS $$\n"
    sqlCreateFunction += "SELECT ('{'||string_agg(to_json(key)||':'||value, ',')||'}')::jsonb\n"
    sqlCreateFunction += "FROM (\n"
    sqlCreateFunction += "WITH to_merge AS (\n"
    sqlCreateFunction += "SELECT * FROM jsonb_each(merge_data)\n"
    sqlCreateFunction += ")\n"
    sqlCreateFunction += "SELECT *\n"
    sqlCreateFunction += "FROM jsonb_each(data)\n"
    sqlCreateFunction += "WHERE key NOT IN (SELECT key FROM to_merge)\n"
    sqlCreateFunction += "UNION ALL\n"
    sqlCreateFunction += "SELECT * FROM to_merge\n"
    sqlCreateFunction += ") t;\n"
    sqlCreateFunction += "$$;\n"

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    _, err = session.Exec(sqlCreateFunction)
    if err != nil {
        return err
    }
    return nil
}

func (db *PostgresDB) Create_generate_where_from_parameters_Function() (error) {
    sqlCreateFunction := "CREATE OR REPLACE FUNCTION generate_where_from_parameters(parameters jsonb)\n"
    sqlCreateFunction += "RETURNS text AS $$\n"
    sqlCreateFunction += "DECLARE\n"
    sqlCreateFunction += "    ary_keys text[];\n"
    sqlCreateFunction += "    field_key text;\n"
    sqlCreateFunction += "    where_clause text := '';\n"
    sqlCreateFunction += "    param_value text := '';\n"
    sqlCreateFunction += "BEGIN\n"
    sqlCreateFunction += "    SELECT array_agg(x.key_name)\n"
    sqlCreateFunction += "    FROM (SELECT jsonb_object_keys(parameters) key_name) x\n"
    sqlCreateFunction += "    INTO ary_keys;\n"
    sqlCreateFunction += "\n"
    sqlCreateFunction += "    IF ary_keys IS NULL THEN\n"
    sqlCreateFunction += "        RETURN '';\n"
    sqlCreateFunction += "    END IF;\n"
    sqlCreateFunction += "\n"
    sqlCreateFunction += "    FOREACH field_key IN ARRAY ary_keys LOOP\n"
    sqlCreateFunction += "        IF where_clause = '' THEN\n"
    sqlCreateFunction += "            where_clause = 'WHERE ';\n"
    sqlCreateFunction += "        END IF;\n"
    sqlCreateFunction += "        IF where_clause <> 'WHERE ' THEN\n"
    sqlCreateFunction += "            where_clause = where_clause || ' AND ';\n"
    sqlCreateFunction += "        END IF;\n"
    sqlCreateFunction += "        param_value = ' = ' || quote_literal(cast(parameters->>field_key as text));\n"
    sqlCreateFunction += "        IF (cast(parameters->>field_key as text) IS NULL) THEN\n"
    sqlCreateFunction += "            param_value = ' IS NULL';\n"
    sqlCreateFunction += "        END IF;\n"
    sqlCreateFunction += "        IF POSITION('.' IN field_key) > 0 THEN\n"
    sqlCreateFunction += "            where_clause = where_clause || 'cast(content->>' || quote_literal(split_part(field_key,'.',2)) || ' as text)' || param_value;\n"
    sqlCreateFunction += "        ELSE\n"
    sqlCreateFunction += "            where_clause = where_clause || quote_ident(field_key) || param_value;\n"
    sqlCreateFunction += "        END IF;\n"
    sqlCreateFunction += "    END LOOP;\n"
    sqlCreateFunction += "    \n"
    sqlCreateFunction += "    RETURN where_clause;\n"
    sqlCreateFunction += "END; $$\n"
    sqlCreateFunction += "LANGUAGE plpgsql;\n"

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    _, err = session.Exec(sqlCreateFunction)
    if err != nil {
        return err
    }
    return nil
}

func (db *PostgresDB) Create_select_latest_Function() (error) {
    sqlCreateFunction := "CREATE OR REPLACE FUNCTION select_latest(select_from text, filter jsonb)\n"
    sqlCreateFunction += " RETURNS TABLE (\n"
    sqlCreateFunction += "   record_id bigint,\n"
    sqlCreateFunction += "   id text,\n"
    sqlCreateFunction += "   created timestamptz,\n"
    sqlCreateFunction += "   updated timestamptz,\n"
    sqlCreateFunction += "   content jsonb,\n"
    sqlCreateFunction += "   deprecated boolean,\n"
    sqlCreateFunction += "   deleted timestamptz\n"
    sqlCreateFunction += "   ) AS $$\n"
    sqlCreateFunction += "DECLARE\n"
    sqlCreateFunction += "    where_clause text;\n"
    sqlCreateFunction += "    sql_query text;\n"
    sqlCreateFunction += "BEGIN\n"
    sqlCreateFunction += "    SELECT * FROM generate_where_from_parameters(filter) INTO where_clause;\n"
    sqlCreateFunction += "\n"
    sqlCreateFunction += "    sql_query = 'SELECT * FROM (SELECT x.* FROM ONLY ' || quote_ident(select_from) || ' x, ' ||\n"
    sqlCreateFunction += "            '(SELECT id, max(record_id) AS max_record_id ' ||\n"
    sqlCreateFunction += "            'FROM ONLY ' || quote_ident(select_from) || ' ';\n"
    sqlCreateFunction += "            IF where_clause <> '' THEN\n"
    sqlCreateFunction += "                sql_query = sql_query || where_clause || ' ';\n"
    sqlCreateFunction += "            END IF;\n"
    sqlCreateFunction += "    sql_query = sql_query ||\n"
    sqlCreateFunction += "            'GROUP BY id ' ||\n"
    sqlCreateFunction += "            'ORDER BY max_record_id DESC) mcd ' ||\n"
    sqlCreateFunction += "            'WHERE x.record_id = mcd.max_record_id) results ';\n"
    sqlCreateFunction += "    sql_query = sql_query || 'ORDER BY record_id DESC';\n"
    sqlCreateFunction += "\n"
    sqlCreateFunction += "    RETURN QUERY EXECUTE sql_query;\n"
    sqlCreateFunction += "END; $$\n"
    sqlCreateFunction += "LANGUAGE plpgsql;\n"

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    _, err = session.Exec(sqlCreateFunction)
    if err != nil {
        return err
    }

    sqlCreateSelectLatestWithArchivedFunction := strings.Replace(sqlCreateFunction, "FROM ONLY ", "FROM ", -1)
    sqlCreateSelectLatestWithArchivedFunction =  strings.Replace(
        sqlCreateSelectLatestWithArchivedFunction,
        "select_latest",
        "select_latest_with_archived", -1)

    _, err = session.Exec(sqlCreateSelectLatestWithArchivedFunction)
    if err != nil {
        return err
    }

    return nil
}
