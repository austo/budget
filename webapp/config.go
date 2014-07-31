package main

import (
	"encoding/json"
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
