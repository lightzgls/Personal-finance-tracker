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
	SourceType string    `db:"source_type"`
	Balance    float64   `db:"balance"`
	CreatedAt  time.Time `db:"created_at"`
}

type Transaction struct {
	TransactionID   uuid.UUID `db:"transaction_id"`
	CategoryType    uuid.UUID `db:"category_type"`
	CategoryName    string    `db:"category_name"`
	Amount          float64   `db:"amount"`
	Description     string    `db:"description"`
	TransactionDate time.Time `db:"transaction_date"`
	CreatedAt       time.Time `db:"created_at"`
	SourceName      string    `db:"source_name"`
	SourceType      string    `db:"source_type"`
}
type TransactionInfo struct {
	Amount          float64   `db:"amount"`
	CategoryType    string    `db:"category_type"`
	CategoryName    string    `db:"category_name"`
	Description     string    `db:"description"`
	TransactionDate time.Time `db:"transaction_date"`
	SourceName      string    `db:"source_name"`
	SourceType      string    `db:"source_type"`
}

// APIResponse defines the structure for the JSON response.
type APIResponse struct {
	Balance      float64           `json:"balance"`
	MonthIncome  float64           `json:"month_income"`
	MonthExpense float64           `json:"month_expense"`
	Transactions []TransactionInfo `json:"transactions"`
}
