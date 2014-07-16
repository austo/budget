package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"os"
	"time"
)

type BudgetItem struct {
	ItemId           int
	AccountId        int
	FiscalYearId     int
	AccountName      string
	ItemDate         time.Time
	CounterParty     string
	ItemDescription  string
	Amount           float64
	RemainingBalance float64
}

var storedProcedures map[string]string = map[string]string{
	"getBudgetItemsByFiscalYear": "exec GetBudgetItemsByFiscalYear @accountId = ?, @fiscalYearId = ?",
	"getAccountsByFiscalYear":    "exec GetAccountsByFiscalYear @fiscalYearId = ?",
}

var statements map[string]*sql.Stmt

var debug = flag.Bool("debug", false, "enable debugging")
var password = flag.String("password", "", "the database password")
var port *int = flag.Int("port", 1433, "the database port")
var server = flag.String("server", "", "the database server")
var user = flag.String("user", "", "the database user")
var dbname = flag.String("dbname", "", "budget database")

func getBudgetItem(rows *sql.Rows, acctId int, fiscalYearId int) (item BudgetItem, err error) {
	var itemId int
	var itemDate time.Time
	var counterparty string
	var itemDescription string
	var amount float64
	var remainingBal float64
	err = rows.Scan(&itemId, &itemDate, &counterparty, &itemDescription, &amount, &remainingBal)
	item = BudgetItem{
		ItemId:           itemId,
		AccountId:        acctId,
		FiscalYearId:     fiscalYearId,
		ItemDate:         itemDate,
		CounterParty:     counterparty,
		ItemDescription:  itemDescription,
		Amount:           amount,
		RemainingBalance: remainingBal,
	}
	return
}

func initializeStatementMap(conn *sql.DB) error {
	if statements != nil {
		return nil
	}
	statements = make(map[string]*sql.Stmt)
	for key, value := range storedProcedures {
		stmt, err := conn.Prepare(value)
		if err != nil {
			return err
		}
		statements[key] = stmt
	}
	return nil
}

func disposeStatementMap() {
	if statements == nil {
		return
	}
	for _, value := range statements {
		value.Close()
	}
}

func main() {
	flag.Parse() // parse the command line args

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
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()

	err = initializeStatementMap(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer disposeStatementMap()

	rows, err := statements["getBudgetItemsByFiscalYear"].Query(30, 1)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		log.Fatal("error fetching rows from GetBudgetItemsByFiscalYear")
	}

	items := make([]BudgetItem, 0)

	for rows.Next() {
		// columnNames, _ := rows.Columns()
		// fmt.Println(columnNames)
		item, _ := getBudgetItem(rows, 30, 1)
		items = append(items, item)
	}

	enc := json.NewEncoder(os.Stdout)
	if err = enc.Encode(&items); err != nil {
		log.Fatal(err)
	}
}
