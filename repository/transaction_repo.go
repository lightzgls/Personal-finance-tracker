package repository

import (
	"context"
	"errors"
	"finance-tracker/model"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)
func CheckSourceActive(db *pgx.Conn, name string) (string, error) {
    var isActive bool
    sql := "SELECT is_active FROM account WHERE source_name = $1"

    err := db.QueryRow(context.Background(), sql, name).Scan(&isActive)
    if err != nil {
        if err == pgx.ErrNoRows {
            // No row was found, so the source does not exist.
            return "not_found", nil
        }
        // Any other error is a real database problem.
        return "", err
    }

    if isActive {
        return "active", nil
    }
    
    return "inactive", nil
}

func GetAllTransactions(db *pgx.Conn) ([]model.TransactionInfo, error) {

	rows, err := db.Query(context.Background(), `SELECT
													T.TRANSACTION_ID,
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
		err := rows.Scan(&t.TransactionID,&t.Amount, &t.CategoryType, &t.CategoryName, &t.TransactionDate, &t.SourceName)
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
	} else if amount < 0 {
		return ErrNegativeAmount
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
	defer tx.Rollback(context.Background())

	var updateQuery string
	if categoryType == "expense" {
		var currentBalance float64
		checkBalanceQuery := `SELECT balance FROM ACCOUNT WHERE source_name = $1;`

		err := db.QueryRow(context.Background(),checkBalanceQuery, req.SourceName).Scan(&currentBalance)
		if err != nil {
			return fmt.Errorf("error checking balance for source '%s': %w", req.SourceName, err)
		}
		if currentBalance < amount {
			return ErrNotEnoughBalance
		}
		updateQuery = `UPDATE ACCOUNT SET balance = balance - $1 WHERE source_name = $2;` 
	} else {
		updateQuery = `UPDATE ACCOUNT SET balance = balance + $1 WHERE source_name = $2;`
	}

	cmdTag, err := tx.Exec(context.Background(), updateQuery, amount, req.SourceName)
	if err != nil {
		log.Printf("ERROR updating balance: %v", err)
		return err
	}
    
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("account with source_name '%s' not found", req.SourceName)
	}


	insertQuery := `INSERT INTO TRANSACTION 
					  (category_type, category_name, amount, transaction_date, source_name)
					  VALUES ($1, $2, $3, $4, $5);`

	_, err = tx.Exec(context.Background(), insertQuery, req.CategoryType, req.CategoryName, amount, req.TransactionDate, req.SourceName)
	if err != nil {
		log.Printf("ERROR inserting transaction: %v", err)
		return err
	}

	log.Println("Success adding new transaction")
	return tx.Commit(context.Background())
}

var ErrDuplicateSource = errors.New("repository: source with that name already exists")
var ErrInvalidBalance = errors.New("repository: initial balance cannot be negative")
var ErrNotEnoughBalance = errors.New("repository: the choosen source doesnt have enough in balance")
var ErrNegativeAmount = errors.New("repository: The transaction amount cant be negative")

func AddSource(db *pgx.Conn,a model.AddSourceRequest) error {
	var SourceQuery string
	status,err := CheckSourceActive(db,a.SourceName);
	if status == "active" {
		return ErrDuplicateSource
	} else if status == "inactive" {
		SourceQuery = "UPDATE account set is_active = true, balance = balance + $1 WHERE source_name = $2"
	} else if  status == "not_found" {
		SourceQuery = "INSERT INTO ACCOUNT (source_name,balance) VALUES ($1,$2);"
	} else {
		return err
	}
	var balance float64 

	if a.Balance == "" {
		balance = 0.0
	} else {
		var err error
		balance, err = strconv.ParseFloat(a.Balance, 64)
		if err != nil {
			log.Printf("An unexpected error occurred: %v", err)
			return err
		}
	}

	if balance < 0 {
        return ErrInvalidBalance
    }
	_,err = db.Exec(context.Background(),SourceQuery,a.SourceName,balance)
	if err != nil {
		log.Printf("Error adding new source: %v\n",err)
		return err
	}
	return nil
}



func GetSummary(db *pgx.Conn) (balance, monthIncome, monthExpense float64, Error error) {
	BalanceQuery := `SELECT COALESCE(SUM(balance), 0) FROM account WHERE is_active = TRUE;` 

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
	query := `SELECT * FROM ACCOUNT WHERE is_active = TRUE`
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		log.Printf("ERROR querying: %v", err)
	}
	var AllSource []model.Account
	for rows.Next() {
		var a model.Account
		err := rows.Scan(&a.SourceName, &a.Balance, &a.CreatedAt,&a.IsActive)
		if err != nil {
			log.Printf("ERROR scanning row: %v\n", err)
			return nil, err
		}
		AllSource = append(AllSource, a)
	}
	return AllSource, nil
}


func GetAllSoucesName(db *pgx.Conn) ([]string, error) {
	query := `SELECT source_name FROM account WHERE is_active = TRUE`
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
func DeleteTransactionsByIDs(db *pgx.Conn, ids []uuid.UUID) (int64, error) {
    sql := "DELETE FROM transaction WHERE transaction_id = ANY($1)"
    
    result, err := db.Exec(context.Background(), sql, ids)
    if err != nil {
        return 0, err
    }

    return result.RowsAffected(), nil
}

func InactiveSources(db *pgx.Conn, names []string) (int64, error) {
    sql := "UPDATE account SET is_active = FALSE WHERE source_name = ANY($1)"
    log.Println("repo")
    result, err := db.Exec(context.Background(), sql, names)
    if err != nil {
        return 0, err
    }

    return result.RowsAffected(), nil
}