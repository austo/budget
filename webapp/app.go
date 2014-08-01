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
var cfgFile = flag.String("f", "../webapp/config.json", "configuration file")

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

	connStr := cfg.dbConnStr()

	if *debug {
		fmt.Printf("\x1B[1m\x1B[33mconfiguration:\x1B[39m\x1B[22m\n %v\n", cfg)
		fmt.Printf("\x1B[1m\x1B[33mconnString:\x1B[39m\x1B[22m\n %s\n", connStr)
		fmt.Println("\x1B[1m\x1B[33mflag values:\x1B[39m\x1B[22m")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf(" name: %s, value: %v\n", f.Name, f.Value)
		})
		fmt.Println()
	}

	db := database.NewDb()
	err = db.Open(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Dispose()

	handlers := makeHandlers(db)

	r := mux.NewRouter()
	r.HandleFunc("/action", handlers["action"]).Methods("GET")
	r.PathPrefix("/{js|css|img}/").Handler(makeFileServer(http.FileServer(http.Dir("../assets/"))))
	r.PathPrefix("/").Handler(makeFileServer(http.FileServer(http.Dir("../assets/html/"))))
	http.Handle("/", r)

	hostAddr := fmt.Sprintf("127.0.0.1:%s", cfg.port)
	log.Printf("Listening at %s...\n", hostAddr)
	http.ListenAndServe(hostAddr, nil)
}

func makeHandlers(db *database.Db) map[string]http.HandlerFunc {
	handlers := make(map[string]http.HandlerFunc)
	// TODO: proper error handling and return codes
	handlers["action"] = makeHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
	return handlers
}

func makeHandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &loggingResponseWrapper{200, w}
		fn(rw, r)
		log.Printf("%d %s %s\n", rw.status, r.Method, r.RequestURI)
	}
}

type fileServer struct {
	handler http.Handler
}

type loggingResponseWrapper struct {
	status int
	http.ResponseWriter
}

func (w *loggingResponseWrapper) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (fs fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := &loggingResponseWrapper{200, w}
	fs.handler.ServeHTTP(rw, r)
	log.Printf("%d %s %s\n", rw.status, r.Method, r.RequestURI)
}

func makeFileServer(h http.Handler) fileServer {
	fs := fileServer{
		handler: h,
	}
	return fs
}
