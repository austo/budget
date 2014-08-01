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

var debug = flag.Bool("d", false, "enable debugging")
var env = flag.String("e", "home", "runtime environment")
var cfgFile = flag.String("f", "config.json", "configuration file")

const (
	DATE_FMT      = "2006-01-02"
	BAD_DATA      = "bad data"
	EBUDGET_ITEMS = "failed to retrieve budget items from database"
)

func main() {
	flag.Parse()

	cfg, err := readConfig(*cfgFile, *env)
	if err != nil {
		log.Fatal(err)
	}

	if *debug {
		fmt.Println(cfg)
	}

	connString := makeConnStr(cfg)

	if *debug {
		fmt.Printf("connString:%s\n", connString)
	}

	db := database.NewDb()
	err = db.Open(connString)
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

	hostAddr := fmt.Sprintf("127.0.0.1:%s", cfg.port)
	log.Printf("Listening at %s...\n", hostAddr)
	http.ListenAndServe(hostAddr, nil)
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

func makeConnStr(cfg config) string {
	return fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%s;database=%s",
		cfg.db.server, cfg.db.user, cfg.db.password, cfg.db.port, cfg.db.database)
}
