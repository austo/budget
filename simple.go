package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

// type BudgetItem struct {
// 	ItemId           int
// 	AccountId        int
// 	FiscalYearId     int
// 	AccountName      int
// 	ItemDate         int
// 	Counterparty     string
// 	ItemDescription  string
// 	Amount           float64
// 	RemainingBalance float64
// }

var debug = flag.Bool("debug", false, "enable debugging")
var password = flag.String("password", "", "the database password")
var port *int = flag.Int("port", 1433, "the database port")
var server = flag.String("server", "", "the database server")
var user = flag.String("user", "", "the database user")
var dbname = flag.String("dbname", "", "budget database")

func main() {
	flag.Parse() // parse the command line args

	if *debug {
		fmt.Printf(" password:%s\n", *password)
		fmt.Printf(" port:%d\n", *port)
		fmt.Printf(" server:%s\n", *server)
		fmt.Printf(" user:%s\n", *user)
		fmt.Printf(" dbname:%s\n", *dbname)
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", *server, *user, *password, *port, *dbname)
	if *debug {
		fmt.Printf(" connString:%s\n", connString)
	}
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()

	// stmt, err := conn.Prepare("select 1, 'abc'")
	// if err != nil {
	// 	log.Fatal("Prepare failed:", err.Error())
	// }

	stmt, err := conn.Prepare("CALL GetBudgetItemsByFiscalYear(?,?)")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		log.Fatal("error calling GetBudgetItemsByFiscalYear")
	}
	fmt.Print("statement: %+v\n", stmt)
	defer stmt.Close()

	rows, err := stmt.Query(30, 1)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		log.Fatal("error fetching rows from GetBudgetItemsByFiscalYear")
	}

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		fmt.Printf("id: %d\n", id)
	}

	// var somenumber int64
	// var somechars string
	// err = row.Scan(&somenumber, &somechars)
	// if err != nil {
	// 	log.Fatal("Scan failed:", err.Error())
	// }
	// fmt.Printf("somenumber:%d\n", somenumber)
	// fmt.Printf("somechars:%s\n", somechars)

	fmt.Printf("bye\n")

}
