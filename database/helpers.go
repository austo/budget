package database

import (
	"database/sql"
	"github.com/austo/budget/dto"
	"time"
)

func getBudgetItem(rows *sql.Rows, accountId int, fiscalYearId int) (item dto.BudgetItem, err error) {
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
	item = dto.BudgetItem{
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

func getAccount(rows *sql.Rows) (account dto.Account, err error) {
	var accountId int
	var accountName string
	var income bool
	var projectedIncome float64
	var startingBalance float64
	err = rows.Scan(&accountId, &accountName, &income, &projectedIncome, &startingBalance)
	if err != nil {
		return
	}
	account = dto.Account{
		AccountId:       accountId,
		AccountName:     accountName,
		Income:          income,
		ProjectedIncome: projectedIncome,
		StartingBalance: startingBalance,
	}
	return
}

func getActivityReportItem(rows *sql.Rows) (item dto.ActivityReportItem, err error) {
	var itemDate time.Time
	var accountName string
	var counterparty string
	var amount float64
	var previousBalance float64
	var remainingBalance float64
	err = rows.Scan(&itemDate, &accountName, &counterparty, &amount, &previousBalance, &remainingBalance)
	if err != nil {
		return
	}
	item = dto.ActivityReportItem{
		ItemDate:         itemDate,
		AccountName:      accountName,
		CounterParty:     counterparty,
		Amount:           amount,
		PreviousBalance:  previousBalance,
		RemainingBalance: remainingBalance,
	}
	return
}
