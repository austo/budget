package database

import (
	"database/sql"
	"github.com/austo/budget/dto"
	_ "github.com/denisenkom/go-mssqldb"
)

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

func (db *Db) GetBudgetItems(accountId int, fiscalYearId int) (items []dto.BudgetItem, err error) {
	rows, err := db.statements["getBudgetItemsByFiscalYear"].Query(accountId, fiscalYearId)
	if err != nil {
		return
	}
	items = make([]dto.BudgetItem, 0)

	for rows.Next() {
		item, _ := getBudgetItem(rows, accountId, fiscalYearId)
		items = append(items, item)
	}
	return
}

func (db *Db) GetAccounts(fiscalYearId int) (accounts []dto.Account, err error) {
	rows, err := db.statements["getAccountsByFiscalYear"].Query(fiscalYearId)
	if err != nil {
		return
	}
	accounts = make([]dto.Account, 0)

	for rows.Next() {
		account, _ := getAccount(rows)
		accounts = append(accounts, account)
	}
	return
}
