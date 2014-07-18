package budget

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"time"
)

// TODO: read from text file
var storedProcedures map[string]string = map[string]string{
	"getBudgetItemsByFiscalYear": "exec GetBudgetItemsByFiscalYear @accountId = ?, @fiscalYearId = ?",
	"getAccountsByFiscalYear":    "exec GetAccountsByFiscalYear @fiscalYearId = ?",
}

type Db struct {
	conn       *sql.DB
	statements map[string]*sql.Stmt
}

func NewDb() *Db {
	return new(Db)
}

func (db *Db) Open(connString string) (err error) {
	db.conn, err = sql.Open("mssql", connString)
	if err != nil {
		return
	}
	return db.initializeStatementMap()
}

func (db *Db) Dispose() {
	db.disposeStatementMap()
	db.conn.Close()
}

func (db *Db) initializeStatementMap() (err error) {
	if db.statements != nil {
		return
	}
	db.statements = make(map[string]*sql.Stmt)
	for key, value := range storedProcedures {
		var stmt *sql.Stmt
		stmt, err = db.conn.Prepare(value)
		if err != nil {
			return
		}
		db.statements[key] = stmt
	}
	return
}

func (db *Db) disposeStatementMap() {
	if db.statements == nil {
		return
	}
	for _, value := range db.statements {
		value.Close()
	}
}

func (db *Db) GetBudgetItems(accountId int, fiscalYearId int) (items []BudgetItem, err error) {
	rows, err := db.statements["getBudgetItemsByFiscalYear"].Query(accountId, fiscalYearId)
	if err != nil {
		return
	}
	items = make([]BudgetItem, 0)

	for rows.Next() {
		item, _ := getBudgetItem(rows, accountId, fiscalYearId)
		items = append(items, item)
	}
	return
}

func (db *Db) GetAccounts(fiscalYearId int) (accounts []Account, err error) {
	rows, err := db.statements["getAccountsByFiscalYear"].Query(fiscalYearId)
	if err != nil {
		return
	}
	accounts = make([]Account, 0)

	for rows.Next() {
		account, _ := getAccount(rows)
		accounts = append(accounts, account)
	}
	return
}

func getBudgetItem(rows *sql.Rows, accountId int, fiscalYearId int) (item BudgetItem, err error) {
	var itemId int
	var itemDate time.Time
	var counterparty string
	var itemDescription string
	var amount float64
	var remainingBal float64
	err = rows.Scan(&itemId, &itemDate, &counterparty, &itemDescription, &amount, &remainingBal)
	if err != nil {
		return
	}
	item = BudgetItem{
		ItemId:           itemId,
		AccountId:        accountId,
		FiscalYearId:     fiscalYearId,
		ItemDate:         itemDate,
		CounterParty:     counterparty,
		ItemDescription:  itemDescription,
		Amount:           amount,
		RemainingBalance: remainingBal,
	}
	return
}

func getAccount(rows *sql.Rows) (account Account, err error) {
	var accountId int
	var accountName string
	var income bool
	var projectedIncome float64
	var startingBalance float64
	err = rows.Scan(&accountId, &accountName, &income, &projectedIncome, &startingBalance)
	if err != nil {
		return
	}
	account = Account{
		AccountId:       accountId,
		AccountName:     accountName,
		Income:          income,
		ProjectedIncome: projectedIncome,
		StartingBalance: startingBalance,
	}
	return
}
