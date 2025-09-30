package repository

import (
	"context"
	"finance-tracker/model"
	"log"
	"time"
	"strings"
	"github.com/jackc/pgx/v5"
)

func GetAllTransactions(db *pgx.Conn) ([]model.TransactionInfo, error) {
	rows, err := db.Query(context.Background(), `SELECT
													T.AMOUNT,
													T.CATEGORY_TYPE,
													T.CATEGORY_NAME,
													T.DESCRIPTION,
													T.TRANSACTION_DATE,
													A.SOURCE_NAME,
													A.SOURCE_TYPE
												FROM TRANSACTION T
													JOIN ACCOUNT A ON T.SOURCE_TYPE = A.SOURCE_TYPE AND T.SOURCE_NAME = A.SOURCE_NAME;
												`)
	if err != nil {
		log.Printf("ERROR querying : %v\n", err)
	}
	defer rows.Close()

	var AllTransactions []model.TransactionInfo
	for rows.Next() {
		var t model.TransactionInfo
		err := rows.Scan(&t.Amount, &t.CategoryType, &t.CategoryName, &t.Description , &t.TransactionDate, &t.SourceName, &t.SourceType)
		if err != nil {
			log.Printf("ERROR scanning row: %v\n", err)
			return nil, err
		}
		AllTransactions = append(AllTransactions, t)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return AllTransactions, nil
}

func AddTransactions(db *pgx.Conn, amount float64,
	category_type, category_name, description, sourceName, SourceType string, transactionDate time.Time) error {
	UpdateBalanceAdd := `UPDATE ACCOUNT 
						SET balance = balance + $1
						WHERE category_type = $2 AND category_name = $3;`
	UpdateBalanceSub := `UPDATE ACCOUNT 
						SET balance = balance - $1
						WHERE category_type = $2 AND category_name = $3;`
	insertTransaction := `INSERT INTO TRANSACTION 
						(category_type, category_name, amount, description, transaction_date, source_name, source_type)
						VALUES ($1, $2, $3, $4, $5, $6, $7);`
	
	if strings.ToLower(category_type) == "expense" {
		_, err := db.Exec(context.Background(), UpdateBalanceSub, amount, category_type, category_name)
		if err != nil {
			log.Printf("ERROR updating: %v", err)
			return err
		}
	} else {
		if strings.ToLower(category_type) == "income" {
		_, err := db.Exec(context.Background(), UpdateBalanceAdd, amount, category_type, category_name)
			if err != nil {
				log.Printf("ERROR updating: %v", err)
				return err
			}
		}	
	}
	_, err := db.Exec(context.Background(), insertTransaction, category_type, category_name, amount, description, transactionDate, sourceName, SourceType)
	if err != nil {
		log.Printf("ERROR inserting: %v", err)
		return err
	}
	return nil
}

func GetSummary(db *pgx.Conn) (balance, monthIncome, monthExpense float64) {
	balanceQuery := `SELECT COALESCE(SUM(balance), 0) FROM account;`
	
	err := db.QueryRow(context.Background(), balanceQuery).Scan(&balance)
	if err != nil {
		log.Printf("ERROR querying total balance: %v\n", err)
	}

	monthlyQuery := `
		SELECT 
			COALESCE(SUM(CASE WHEN LOWER(category_type) = 'income' THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN LOWER(category_type) = 'expense' THEN amount ELSE 0 END), 0)
		FROM 
			transaction
		WHERE 
			DATE_TRUNC('month', transaction_date) = DATE_TRUNC('month', CURRENT_DATE);
	`
	err = db.QueryRow(context.Background(), monthlyQuery).Scan(&monthIncome, &monthExpense)
	if err != nil {
		log.Printf("ERROR querying monthly summary: %v\n", err)
	}
	return
}
func GetAllSources(db *pgx.Conn) ([]model.Account, error){
	query := `SELECT * FROM ACCOUNT`
	rows,err := db.Query(context.Background(),query)
	if err != nil {
		log.Printf("ERROR querying: %v",err)
	}
	var AllSource []model.Account
	for rows.Next() {
		var a model.Account
		err := rows.Scan(&a.SourceName,&a.SourceType,&a.Balance, &a.CreatedAt)
		if err != nil {
			log.Printf("ERROR scanning row: %v\n", err)
			return nil, err
		}
		AllSource = append(AllSource, a)
	}
	return AllSource,nil
}