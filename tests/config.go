package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
)

type config struct {
	Port string
	Db   struct {
		Server   string
		Port     string
		User     string
		Password string
		Database string
	}
}

func main() {
	cfg, err := readConfig("../webapp/config.json", "home")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
	examine(cfg)
}

func readConfig(filename, env string) (config, error) {
	var cfg config
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, err
	}
	var m = map[string]config{}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return cfg, err
	}
	if cfg, ok := m[env]; ok {
		return cfg, nil
	}
	return cfg, fmt.Errorf("no configuration for environment \"%s\"", env)
}

func examine(cfg config) {
	s := reflect.ValueOf(cfg)
	cfgType := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s: %v\n", cfgType.Field(i).Name, f)
	}
}
