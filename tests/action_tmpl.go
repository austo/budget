package main

import (
	"flag"
	"fmt"
	"github.com/austo/budget/database"
	"github.com/austo/budget/dto"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"time"
)

var debug = flag.Bool("debug", false, "enable debugging")
var password = flag.String("password", "", "the database password")
var port *int = flag.Int("port", 1433, "the database port")
var server = flag.String("server", "", "the database server")
var user = flag.String("user", "", "the database user")
var dbname = flag.String("dbname", "GardenClubAccounting", "budget database")

const (
	DATE_FMT      = "2006-Jan-02"
	BAD_DATA      = "bad data"
	EBUDGET_ITEMS = "failed to retrieve budget items from database"
)

var templates = template.Must(template.ParseFiles("../tmpl/activityReport.html"))

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

	handlers := makeHandlers(db)

	rtr := mux.NewRouter()
	rtr.HandleFunc("/action", handlers["action"]).Methods("GET")
	http.Handle("/", rtr)

	fmt.Println(templates)
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
		activityReport := dto.ActivityReport{
			Start:           start,
			End:             end,
			StartingBalance: actions[0].PreviousBalance,
			EndingBalance:   actions[len(actions)-1].RemainingBalance,
			Items:           actions,
		}
		tmplErr := templates.ExecuteTemplate(w, "activityReport.html", activityReport)
		if tmplErr != nil {
			http.Error(w, tmplErr.Error(), http.StatusInternalServerError)
		}
	}
	return handlers
}
