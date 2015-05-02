package database

type Result interface {
    GetID() string
    SetID(value string)
}

type DBResult map[string]interface{}

func (dbr DBResult) GetID() string {
    return dbr["id"].(string)
}

func (dbr DBResult) SetID(value string) {
    dbr["id"] = value
}
