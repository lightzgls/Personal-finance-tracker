package handler

import (
	"encoding/json"
	"finance-tracker/repository"
	"log"
	"net/http"
	"finance-tracker/model"
	"github.com/jackc/pgx/v5"
)


func AddSourceHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req model.AddSourceRequest 
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}
		err := repository.AddSource(db,req)
		if err != nil {
			log.Printf("Failed to add new source: %v", err)
			http.Error(w,"Failed to add new source..",http.StatusInternalServerError)
			return 
		}
		w.Header().Set("Content-Type","Applcation/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "Source added successfully"}`))
		log.Println("Transaction added successfully")
	}
}


func AddTransactionHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req model.AddTransactionRequest

		// 2. Decode the JSON body into the struct
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		// 4. Call the repository function with data from the struct
		err := repository.AddTransactions(db, req)
		if err != nil {
			log.Printf("Failed to add transaction: %v", err) // Log the actual error
			http.Error(w, "Failed to add transaction", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "Transaction added successfully"}`))
		log.Println("Transaction added successfully")
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

// func GetAllSoucesNameHandler(db *pgx.Conn) http.HandlerFunc{
// 	return func (w http.ResponseWriter,r *http.Request)  {
// 		if r.Method != http.MethodGet {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 		names, err := repository.GetAllSoucesName(db)
// 		if err != nil 
// 	}
// }