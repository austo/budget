package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/austo/budget/database"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var debug = flag.Bool("debug", false, "enable debugging")
var password = flag.String("p", "", "the database password")
var port *int = flag.Int("port", 1433, "the database port")
var server = flag.String("s", "", "the database server")
var user = flag.String("u", "", "the database user")
var dbname = flag.String("d", "GardenClubAccounting", "budget database")

const (
	DATE_FMT      = "2006-01-02"
	BAD_DATA      = "bad data"
	EBUDGET_ITEMS = "failed to retrieve budget items from database"
)

func main() {
	flag.Parse()

	if *debug {
		fmt.Printf("password: %s\n", *password)
		fmt.Printf("port: %d\n", *port)
		fmt.Printf("server: %s\n", *server)
		fmt.Printf("user: %s\n", *user)
		fmt.Printf("dbname: %s\n", *dbname)
	}

	connString := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%d;database=%s",
		*server, *user, *password, *port, *dbname)

	if *debug {
		fmt.Printf("connString:%s\n", connString)
	}

	db := database.NewDb()
	err := db.Open(connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Dispose()

	handlers := makeHandlers(db)

	r := mux.NewRouter()
	r.HandleFunc("/action", handlers["action"]).Methods("GET")
	r.PathPrefix("/{js|css|img}/").Handler(http.FileServer(http.Dir("../assets/")))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../assets/html/")))
	http.Handle("/", r)

	log.Println("Listening...")
	http.ListenAndServe("127.0.0.1:3000", nil)
}

func makeHandlers(db *database.Db) map[string]http.HandlerFunc {
	handlers := make(map[string]http.HandlerFunc)
	// TODO: proper error handling and return codes
	handlers["action"] = func(w http.ResponseWriter, r *http.Request) {
		rawStart, rawEnd := r.URL.Query().Get("start"), r.URL.Query().Get("end")
		// TODO: if start/end not present, default to fiscal year start/end
		if rawStart == "" || rawEnd == "" {
			http.Error(w, BAD_DATA, http.StatusBadRequest)
			return
		}
		start, tErr := time.Parse(DATE_FMT, rawStart)
		if tErr != nil {
			http.Error(w, BAD_DATA, http.StatusBadRequest)
			return
		}
		end, tErr := time.Parse(DATE_FMT, rawEnd)
		if tErr != nil {
			http.Error(w, BAD_DATA, http.StatusBadRequest)
			return
		}
		actions, dbErr := db.GetActivityReportItems(start, end)
		if dbErr != nil {
			http.Error(w, EBUDGET_ITEMS, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(&actions)
	}
	return handlers
}
