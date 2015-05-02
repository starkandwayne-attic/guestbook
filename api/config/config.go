package config

import (
    "io/ioutil"
    "encoding/json"
) 

type Config struct {
    Port int `json:"port"`
    DatabaseUri string `json:"databaseUri"`
    DatabaseName string `json:"databaseName"`
}

func ParseConfig(filePath string) (Config) {
    var retval Config
    file, err := ioutil.ReadFile(filePath)
    if err != nil {
        return Config{}
    }
    json.Unmarshal(file, &retval)
    return retval
}
