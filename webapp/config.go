package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	// "strconv"
)

type dbconfig struct {
	server   string
	port     string
	user     string
	password string
	database string
}

type config struct {
	port string
	db   dbconfig
}

func main() {
	cfg, err := readConfig("home", "config.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg)
}

func readConfig(env, filename string) (cfg config, err error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	var f interface{}
	err = json.Unmarshal(bytes, &f)
	if err != nil {
		return
	}
	m := f.(map[string]interface{})
	cfg, err = populateConfig(env, m)
	return
}

func populateConfig(env string, m map[string]interface{}) (cfg config, err error) {
	f := m[env]
	if f == nil {
		err = fmt.Errorf("no configuration for environment \"%s\"", env)
		return
	}
	c := f.(map[string]interface{})
	for k, v := range c {
		switch strings.ToLower(k) {
		case "port":
			cfg.port = getString(v)
		case "db":
			cfg.db = populateDbConfig(v.(map[string]interface{}))
		}
	}
	return
}

func populateDbConfig(m map[string]interface{}) (db dbconfig) {
	for k, v := range m {
		switch strings.ToLower(k) {
		case "server":
			db.server = v.(string)
		case "port":
			db.port = getString(v)
		case "user":
			db.user = v.(string)
		case "password":
			db.password = v.(string)
		case "database":
			db.database = v.(string)
		}
	}
	return
}

func getString(f interface{}) (str string) {
	switch vv := f.(type) {
	case string:
		str = vv
	case float64:
		str = fmt.Sprintf("%0.f", vv)
	}
	return
}
