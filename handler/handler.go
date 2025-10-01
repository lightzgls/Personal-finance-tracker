package handler

import (
	"encoding/json"
	"finance-tracker/model"
	"finance-tracker/repository"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/schema"
	"github.com/jackc/pgx/v5"
)

var decoder = schema.NewDecoder()

func AddSourceHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		var req model.AddSourceRequest
		err := decoder.Decode(&req, r.PostForm)
		if err != nil {
			http.Error(w, "Failed to decode form data", http.StatusBadRequest)
			log.Printf("!!! Failed to decode form data: %v", err)
			return
		}
		err = repository.AddSource(db, req)
		if err != nil {
			log.Printf("Failed to add new source: %v", err)
			http.Error(w, "Failed to add new source..", http.StatusInternalServerError)
			return
		}
		log.Println("Transaction added successfully")
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func AddTransactionHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		log.Printf("Form data received: %+v", r.PostForm)

		var req model.AddTransactionRequest

		err := decoder.Decode(&req, r.PostForm)
		if err != nil {
			http.Error(w, "Failed to decode form data", http.StatusBadRequest)
			log.Printf("!!! Failed to decode form data: %v", err)
			return
		}
		_, err = time.Parse("2006-01-02", req.TransactionDate)
		if err != nil {
			log.Printf("Invalid date format: %v", err)
			http.Error(w, "Invalid date format. Please use dd/mm/YYYY.", http.StatusBadRequest)
			return
		}
		err = repository.AddTransactions(db, req)
		if err != nil {
			log.Printf("Failed to add transaction: %v", err)
			http.Error(w, "Failed to add transaction", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
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
		response := model.GetSummaryResponse{
			Balance:      balance,
			MonthIncome:  monthIncome,
			MonthExpense: monthExpense,
			Transactions: limitedTransactions,
		}

		tmpl, err := template.ParseFiles(("templates/home.html"))
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, response)
		if err != nil {
			// Log the error instead of calling http.Error
			log.Printf("Failed to render template: %v", err)
			// You can optionally try to close the connection or just return
			return
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

// func GetAllSoucesNameHandler(db *pgx.Conn) http.HandlerFunc{
// 	return func (w http.ResponseWriter,r *http.Request)  {
// 		if r.Method != http.MethodGet {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 		names, err := repository.GetAllSoucesName(db)
// 		if err != nil
// 	}
// }
