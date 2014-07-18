package budget

import (
	"time"
)

type Account struct {
	AccountId       int
	AccountName     string
	Income          bool
	ProjectedIncome float64
	StartingBalance float64
}

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
