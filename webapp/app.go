package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/austo/budget/database"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var debug = flag.Bool("debug", false, "enable debugging")
var password = flag.String("password", "", "the database password")
var port *int = flag.Int("port", 1433, "the database port")
var server = flag.String("server", "", "the database server")
var user = flag.String("user", "", "the database user")
var dbname = flag.String("dbname", "GardenClubAccounting", "budget database")

func main() {
	flag.Parse()

	if *debug {
		fmt.Printf(" password:%s\n", *password)
		fmt.Printf(" port:%d\n", *port)
		fmt.Printf(" server:%s\n", *server)
		fmt.Printf(" user:%s\n", *user)
		fmt.Printf(" dbname:%s\n", *dbname)
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

	r := mux.NewRouter()
	r.HandleFunc("/accounts/{fiscalYearId:\\d+}", getAccounts(db)).Methods("GET")
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../assets/"))))
	r.PathPrefix("/{js|css}/").Handler(http.FileServer(http.Dir("../assets/")))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../assets/html/")))
	http.Handle("/", r)

	log.Println("Listening...")
	http.ListenAndServe("127.0.0.1:3000", nil)
}

func getAccounts(db *database.Db) http.HandlerFunc {
	// TODO: proper error handling and return codes
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		idStr := params["fiscalYearId"]
		fiscalYearId, _ := strconv.ParseInt(idStr, 10, 32)
		enc := json.NewEncoder(w)
		accounts, _ := db.GetAccounts(int(fiscalYearId))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = enc.Encode(&accounts)
	}
}
