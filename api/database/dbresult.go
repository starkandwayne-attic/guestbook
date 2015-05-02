package database

import (
    "time"
)

type Result interface {
    GetID() string
    SetID(value string)
    GetCreated() time.Time
    SetCreated(value time.Time)
    GetUpdated() time.Time
    SetUpdated(value time.Time)
}

type DBResult map[string]interface{}

func (dbr DBResult) GetID() string {
    return dbr["id"].(string)
}

func (dbr DBResult) SetID(value string) {
    dbr["id"] = value
}

func (dbr DBResult) GetCreated() time.Time {
    return EnsurePropertyIsTime(dbr["created"])
}

func (dbr DBResult) SetCreated(value time.Time) {
    dbr["created"] = value
}

func (dbr DBResult) GetUpdated() time.Time {
    return EnsurePropertyIsTime(dbr["updated"])
}

func (dbr DBResult) SetUpdated(value time.Time) {
    dbr["updated"] = value
}

func EnsurePropertyIsTime(prop interface{}) (time.Time) {
    dateValue, isDate := prop.(time.Time)

    if isDate {
        return dateValue
    }

    stringValue, isString := prop.(string)

    if isString {
        t, _ := time.Parse(time.RFC3339, stringValue)
        return t
    }
    return time.Now()
}
