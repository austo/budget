package database

var storedProcedures map[string]string = map[string]string{
	"getBudgetItemsByFiscalYear": "exec GetBudgetItemsByFiscalYear @accountId = ?, @fiscalYearId = ?",
	"getAccountsByFiscalYear":    "exec GetAccountsByFiscalYear @fiscalYearId = ?",
	"getActivityReport":          "exec ActivityReportActiveFiscalYear @startDate = ?, @endDate = ?",
}
