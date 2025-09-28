package model

import (
	"time"

	"github.com/google/uuid"
)


type Category struct {
	CategoryID   uuid.UUID `db:"category_id"`
	CategoryName string    `db:"category_name"`
	CategoryType string    `db:"category_type"`
}

type Account struct {
	SourceName string    `db:"source_name"`
	SourceType string    `db:"source_type"`
	Balance    float64   `db:"balance"`
	CreatedAt  time.Time `db:"created_at"`
}

type Transaction struct {
	TransactionID   uuid.UUID `db:"transaction_id"`
	CategoryID      uuid.UUID `db:"category_id"`
	Amount          float64   `db:"amount"`
	Description     string    `db:"description"`
	TransactionDate time.Time `db:"transaction_date"`
	CreatedAt       time.Time `db:"created_at"`
}
