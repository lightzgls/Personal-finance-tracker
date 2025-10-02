package model

import (
	"time"

	"github.com/google/uuid"
)

type PageData struct {
	Title   string
	Content string
}

type Account struct {
	SourceName string    `db:"source_name"`
	Balance    float64   `db:"balance"`
	CreatedAt  time.Time `db:"created_at"`
}

type Transaction struct {
	TransactionID   uuid.UUID `db:"transaction_id"`
	CategoryType    uuid.UUID `db:"category_type"`
	CategoryName    string    `db:"category_name"`
	Amount          float64   `db:"amount"`
	TransactionDate time.Time `db:"transaction_date"`
	CreatedAt       time.Time `db:"created_at"`
	SourceName      string    `db:"source_name"`
}
type TransactionInfo struct {
	Amount          float64   `db:"amount"`
	CategoryType    string    `db:"category_type"`
	CategoryName    string    `db:"category_name"`
	TransactionDate time.Time `db:"transaction_date"`
	SourceName      string    `db:"source_name"`
}

type GetSummaryResponse struct {
	Balance          float64
	MonthIncome      float64
	MonthExpense     float64
	Transactions     []TransactionInfo
	FormErrors       map[string]string
	ShowTransPopup   bool
	AllTransactions  []TransactionInfo
	AvailableSources []string
	ShowSourcesPopup bool
	AllSources       []Account
}

type AddTransactionRequest struct {
	Amount          string `schema:"amount"`
	CategoryType    string `schema:"transaction_type"`
	CategoryName    string `schema:"category_name"`
	SourceName      string `schema:"source_name"`
	TransactionDate string `schema:"transaction_date"`
}
type AddSourceRequest struct {
	SourceName string `schema:"source_name"`
	Balance    string `schema:"balance"`
	FormErrors map[string]string
}
