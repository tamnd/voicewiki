package middleware

import (
	"encoding/json"
	"io/ioutil"
)

var Config = &struct {
	App struct {
		Bind        string `json:"bind"`
		WebRoot     string `json:"webroot"`
		Env string `json:"env"`
	} `json:"app"`
}{}

func LoadConfig(file string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, Config)
}
