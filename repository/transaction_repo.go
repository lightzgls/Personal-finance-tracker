package repository

import (
	"context"
	"finance-tracker/model"
	"log"
	"strconv"
	"strings"
	"github.com/jackc/pgx/v5"
)

func GetAllTransactions(db *pgx.Conn) ([]model.TransactionInfo, error) {
	rows, err := db.Query(context.Background(), `SELECT
													T.AMOUNT,
													T.CATEGORY_TYPE,
													T.CATEGORY_NAME,
													T.TRANSACTION_DATE,
													A.SOURCE_NAME
												FROM TRANSACTION T
													JOIN ACCOUNT A ON T.SOURCE_NAME = A.SOURCE_NAME
												ORDER BY T.TRANSACTION_DATE DESC;
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
	UpdateBalanceAddQuery := `UPDATE ACCOUNT 
						SET balance = balance + $1
						WHERE  source_name = $2;`
	UpdateBalanceSubQuery := `UPDATE ACCOUNT 
						SET balance = balance - $1
						WHERE source_name = $2;`
	InsertTransactionQuery := `INSERT INTO TRANSACTION 
						(category_type, category_name, amount, transaction_date, source_name)
						VALUES ($1, $2, $3, $4, $5);`


	amount, err := strconv.ParseFloat(req.Amount,64)
	if err != nil {
		log.Printf("Error parsing string to Float64: %v\n",err)
	}

	if strings.ToLower(req.CategoryType) != "income" && strings.ToLower(req.CategoryType) != "expense" {
		return err
	}

	if strings.ToLower(req.CategoryType) == "expense" {
		_, err := db.Exec(context.Background(), UpdateBalanceSubQuery, amount, req.CategoryName)
		if err != nil {
			log.Printf("ERROR updating: %v", err)
			return err
		}
	} else {
		if strings.ToLower(req.CategoryType) == "income" {
			_, err := db.Exec(context.Background(), UpdateBalanceAddQuery, amount, req.CategoryName)
			if err != nil {
				log.Printf("ERROR updating: %v", err)
				return err
			}
		}
	}
	_, err = db.Exec(context.Background(), InsertTransactionQuery, req.CategoryType, req.CategoryName, amount, req.TransactionDate, req.SourceName)
	if err != nil {
		log.Printf("ERROR inserting: %v", err)
		return err
	}
	log.Println("Success adding new transaction")
	return nil

}



func AddSource(db *pgx.Conn,a model.AddSourceRequest) error {
	InsertSourceQuery := "INSERT INTO ACCOUNT (source_name,balance) VALUES ($1,$2);"
	balance, err := strconv.ParseFloat(a.Balance,64)
	if err != nil {
		log.Printf("Error parsing string to float64: %v\n",err)
	}
	_,err = db.Exec(context.Background(),InsertSourceQuery,a.SourceName,balance)
	if err != nil {
		log.Printf("Error adding new source: %v\n",err)
		return err
	}
	return nil
}

func GetSummary(db *pgx.Conn) (balance, monthIncome, monthExpense float64) {
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
	}
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
	return Name, nil
}
