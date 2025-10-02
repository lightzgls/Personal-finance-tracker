package repository

import (
	"context"
	"errors"
	"finance-tracker/model"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

func sourceExist(db *pgx.Conn, name string) (bool, error) {
	var placeholderID string
	ExistCheck := `SELECT source_name FROM ACCOUNT WHERE source_name = $1;`

	err := db.QueryRow(context.Background(), ExistCheck,name).Scan(&placeholderID)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Found no %v in account list\n",name)
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetAllTransactions(db *pgx.Conn) ([]model.TransactionInfo, error) {

	rows, err := db.Query(context.Background(), `SELECT
													T.AMOUNT,
													T.CATEGORY_TYPE,
													T.CATEGORY_NAME,
													T.TRANSACTION_DATE,
													A.SOURCE_NAME
												FROM TRANSACTION T
													JOIN ACCOUNT A ON T.SOURCE_NAME = A.SOURCE_NAME
												ORDER BY T.TRANSACTION_DATE DESC, T.CREATED_AT DESC;
												`)
	if err != nil {
		log.Printf("ERROR querying transactions : %v\n", err)
	}
	defer rows.Close()

	var AllTransactions []model.TransactionInfo
	for rows.Next() {
		var t model.TransactionInfo
		err := rows.Scan(&t.Amount, &t.CategoryType, &t.CategoryName, &t.TransactionDate, &t.SourceName)
		if err != nil {
			log.Printf("ERROR scanning row: %v\n", err)
			return nil, err
		}
		t.CategoryType = strings.ToTitle(t.CategoryType)
		t.SourceName = strings.ToTitle(t.SourceName)
		AllTransactions = append(AllTransactions, t)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return AllTransactions, nil
}

func AddTransactions(db *pgx.Conn, req model.AddTransactionRequest) error {
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		log.Printf("Error parsing string to Float64: %v\n", err)
		return err
	}

	categoryType := strings.ToLower(req.CategoryType)
	if categoryType != "income" && categoryType != "expense" {
		return errors.New("invalid category_type: must be 'income' or 'expense'")
	}

	// 1. Begin a database transaction
	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Printf("ERROR begin a transaction: %v", err)
		return err
	}
	// Defer a rollback in case anything goes wrong. It will be a no-op if we commit.
	defer tx.Rollback(context.Background())

	// 2. Update the account balance within the transaction
	var updateQuery string
	if categoryType == "expense" {
		updateQuery = `UPDATE ACCOUNT SET balance = balance - $1 WHERE source_name = $2;`
	} else {
		updateQuery = `UPDATE ACCOUNT SET balance = balance + $1 WHERE source_name = $2;`
	}

	// Use tx.Exec, not db.Exec
	cmdTag, err := tx.Exec(context.Background(), updateQuery, amount, req.SourceName)
	if err != nil {
		log.Printf("ERROR updating balance: %v", err)
		return err
	}
    
    // Check if any row was actually updated. If not, the account doesn't exist.
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("account with source_name '%s' not found", req.SourceName)
	}


	// 3. Insert the transaction record within the transaction
	insertQuery := `INSERT INTO TRANSACTION 
					  (category_type, category_name, amount, transaction_date, source_name)
					  VALUES ($1, $2, $3, $4, $5);`

	// Use tx.Exec, not db.Exec
	_, err = tx.Exec(context.Background(), insertQuery, req.CategoryType, req.CategoryName, amount, req.TransactionDate, req.SourceName)
	if err != nil {
		log.Printf("ERROR inserting transaction: %v", err)
		return err
	}

	// 4. If all commands succeed, commit the transaction to make the changes permanent
	log.Println("Success adding new transaction")
	return tx.Commit(context.Background())
}

var ErrDuplicateSource = errors.New("repository: source with that name already exists")
var ErrInvalidBalance = errors.New("repository: initial balance cannot be negative")


func AddSource(db *pgx.Conn,a model.AddSourceRequest) error {
	exist, err:= sourceExist(db,a.SourceName)
	if err != nil {
		return err
	}
	if exist {
		return ErrDuplicateSource
	}
	InsertSourceQuery := "INSERT INTO ACCOUNT (source_name,balance) VALUES ($1,$2);"
	balance, err := strconv.ParseFloat(a.Balance,64)
	if err != nil {
		log.Printf("Error parsing string to float64: %v\n",err)
		return err
	}
	if balance < 0 {
        return ErrInvalidBalance
    }
	_,err = db.Exec(context.Background(),InsertSourceQuery,a.SourceName,balance)
	if err != nil {
		log.Printf("Error adding new source: %v\n",err)
		return err
	}
	return nil
}



func GetSummary(db *pgx.Conn) (balance, monthIncome, monthExpense float64, Error error) {
	BalanceQuery := `SELECT COALESCE(SUM(balance), 0) FROM account;` 

	err := db.QueryRow(context.Background(), BalanceQuery).Scan(&balance)
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
func GetAllSources(db *pgx.Conn) ([]model.Account, error) {
	query := `SELECT * FROM ACCOUNT`
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		log.Printf("ERROR querying: %v", err)
	}
	var AllSource []model.Account
	for rows.Next() {
		var a model.Account
		err := rows.Scan(&a.SourceName, &a.Balance, &a.CreatedAt)
		if err != nil {
			log.Printf("ERROR scanning row: %v\n", err)
			return nil, err
		}
		AllSource = append(AllSource, a)
	}
	return AllSource, nil
}


func GetAllSoucesName(db *pgx.Conn) ([]string, error) {
	query := `SELECT source_name FROM ACCOUNT`
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		log.Printf("ERROR querying: %v", err)
		return nil, err
	}
	defer rows.Close()

	var Name []string
	for rows.Next() {
		var a string
		err := rows.Scan(&a)
		if err != nil {
			log.Printf("ERROR scanning row: %v\n", err)
			return nil, err
		}
		Name = append(Name, a)
	}
	if err = rows.Err(); err != nil {
		return nil,err
	}
	return Name, nil
}
