package repository

import (
	"context"
	"finance-tracker/model"
	"log"
	"github.com/jackc/pgx/v5"
)


func GetAllTransactions(db *pgx.Conn) ([]model.Transaction, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM TRANSACTION")
	if err != nil {
		log.Printf("ERROR querying from database: %v\n", err)
	}
	defer rows.Close()
	
	var AllTransactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.TransactionID, &t.CategoryID, &t.Amount, &t.Description, &t.TransactionDate,&t.CreatedAt)
		if err != nil {
			log.Printf("ERROR scanning row: %v\n",err)
			return nil,err
		}
		AllTransactions = append(AllTransactions, t)
	}
	if rows.Err() != nil {
		return  nil,rows.Err()
	}
	return AllTransactions, nil
}
