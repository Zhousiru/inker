package config

import (
	"encoding/json"
	"os"
)

type configuration struct {
	APIAddr               string
	Key                   string
	NormalTokenLifetime   int64
	RememberTokenLifetime int64
	MongoURI              string
	DBName                string
}

// Conf 包含了用户配置
var Conf *configuration

func init() {
	file, _ := os.Open("config.json")
	defer file.Close()

	Conf = new(configuration)
	err := json.NewDecoder(file).Decode(Conf)
	if err != nil {
		panic(err)
	}

	return
}
