package dto

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

type ActivityReportItem struct {
	ItemDate         time.Time
	AccountName      string
	CounterParty     string
	Amount           float64
	PreviousBalance  float64
	RemainingBalance float64
}

type ActivityReport struct {
	Start           time.Time
	End             time.Time
	StartingBalance float64
	EndingBalance   float64
	Items           []ActivityReportItem
}
