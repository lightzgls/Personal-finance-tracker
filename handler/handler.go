package handler

import (
	"encoding/json"
	"finance-tracker/repository"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"finance-tracker/model"
	"github.com/jackc/pgx/v5"
)

func AddTransactionHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		// Get values from the form

		amountStr := r.FormValue("amount")
		category_type := r.FormValue("category_type")
		category_name := r.FormValue("category_name")
		description := r.FormValue("description")
		sourceName := r.FormValue("source_name")
		sourceType := r.FormValue("source_type")
		transactionDateStr := r.FormValue("transaction_date")


		transactionDate, err := time.Parse(time.DateOnly,transactionDateStr)
		if err != nil {
			log.Printf("ERROR convert string to Date: %v\n",err)
			return
		}

		amount,err := strconv.ParseFloat(amountStr,64)
		if err != nil {
			log.Printf("ERROR convert string to float: %v\n",err)
			return
		}
		if strings.ToLower(category_type) != "income" && strings.ToLower(category_type) != "expense" {
			http.Error(w, "Invalid category_type", http.StatusBadRequest)
			return
		}

		// Call the repository function
		err = repository.AddTransactions(db, amount,category_type,category_name,description,sourceName,sourceType,transactionDate)
		if err != nil {
			http.Error(w, "Failed to add transaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Transaction added successfully"))
	}
}



func GetSummaryHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Ensure the method is GET
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 2. Fetch all transactions from the repository
		transactions, err := repository.GetAllTransactions(db)
		if err != nil {
			http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
			return
		}

		// 3. Fetch summary data from the repository
		balance, monthIncome, monthExpense := repository.GetSummary(db)
		limit := 5
		if len(transactions) < limit {
			limit = len(transactions)
		}
		limitedTransactions := transactions[:limit]

		// 4. Construct the response data structure
		response := model.APIResponse{
			Balance:      balance,
			MonthIncome:  monthIncome,
			MonthExpense: monthExpense,
			Transactions: limitedTransactions,
		}

		// 5. Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// 6. Encode the response struct to JSON and write it to the response writer
		if err := json.NewEncoder(w).Encode(response); err != nil {
			// This error happens if the response struct can't be converted to JSON.
			http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		}
	}
}


func GetAllTransactionsHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		transactions, err := repository.GetAllTransactions(db)
		if err != nil {
			http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
			return 
		}
		response := transactions

		w.Header().Set("Content-Type", "application/json")

		// 6. Encode the response struct to JSON and write it to the response writer
		if err := json.NewEncoder(w).Encode(response); err != nil {
			// This error happens if the response struct can't be converted to JSON.
			http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		}

	}
}


func GetAllSourcesHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		sources, err := repository.GetAllSources(db)
		if err != nil {
			http.Error(w, "Failed to fetch sources", http.StatusInternalServerError)
			return 
		}
		response := sources

		w.Header().Set("Content-Type", "application/json")

		// 6. Encode the response struct to JSON and write it to the response writer
		if err := json.NewEncoder(w).Encode(response); err != nil {
			// This error happens if the response struct can't be converted to JSON.
			http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		}

	}
}