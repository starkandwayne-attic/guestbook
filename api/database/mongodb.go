package database

import (
    "time"
    "gopkg.in/mgo.v2"
    _"gopkg.in/mgo.v2/bson"
)

type MongoDB struct {
    DatabaseUri string
    DatabaseName string
}

func UseMongoDB(databaseUri string, databaseName string) (MongoDB) {
    return MongoDB{databaseUri, databaseName}
}

func (db *MongoDB) EnsureStructure(tables []string) (error) {
    return nil
}

func (db *MongoDB) connect() (*mgo.Session, error) {
    mgoSession, err := mgo.Dial(db.DatabaseUri)
    if err != nil {
        return &mgo.Session{}, err
    }
    mgoSession.SetMode(mgo.Monotonic, true)

    return mgoSession, nil
}

func (db *MongoDB) SelectAll(selectFrom string, selectTo *[]DBResult) (error) {
    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    conn := session.DB(db.DatabaseName).C(selectFrom)
    conn.Find(nil).Sort("-updated").All(selectTo)
    if err != nil {
        return err
    }
    return nil
}

func (db *MongoDB) SelectFirst(selectFrom string, query map[string]interface{}, selectTo *DBResult) (error) {
    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    conn := session.DB(db.DatabaseName).C(selectFrom)
    conn.Find(query).Sort("-updated").One(selectTo)
    if err != nil {
        return err
    }
    return nil
}

func (db *MongoDB) Insert(insertInto string, insertObject DBResult) (error) {
    now := time.Now().UTC()
    insertObject.SetCreated(now)
    insertObject.SetUpdated(now)

    queryMap := make(map[string]interface{})
    queryMap["id"] = insertObject.GetID()

    var existingDoc DBResult
    err := db.SelectFirst(insertInto, queryMap, &existingDoc)
    if existingDoc != nil && err == nil {
        insertObject.SetCreated(existingDoc.GetCreated())
    }

    session, err := db.connect()
    if err != nil {
        return err
    }
    defer session.Close()

    conn := session.DB(db.DatabaseName).C(insertInto)

    err = conn.Insert(insertObject)
    if err != nil {
        return err
    }
    return nil
}
