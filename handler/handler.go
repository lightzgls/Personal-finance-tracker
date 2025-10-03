package handler

import (
	"encoding/json"
	"errors"
	"finance-tracker/model"
	"finance-tracker/repository"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/jackc/pgx/v5"
)

var decoder = schema.NewDecoder()

func AddSourceHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 1. Check method and parse the form (this part was correct)
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// 2. Decode form data into your struct
		var req model.AddSourceRequest
		err := decoder.Decode(&req, r.PostForm)
		if err != nil {
			log.Printf("!!! Failed to decode form data: %v", err)
			http.Redirect(w, r, "/home?error=bad_form_data", http.StatusSeeOther)
			return
		}

		err = repository.AddSource(db, req)

		if errors.Is(err, repository.ErrDuplicateSource) {
			log.Print("Duplicate source error, redirecting")
			http.Redirect(w, r, "/home?error=source_already_exist", http.StatusSeeOther)
			return
		} else if errors.Is(err, repository.ErrInvalidBalance) {
			http.Redirect(w, r, "/home?error=negative_balance", http.StatusSeeOther)
			return

		} else if err != nil {
			log.Printf("An unexpected error occurred: %v", err)
			http.Error(w, "An internal server error occurred", http.StatusInternalServerError)
			return
		}
		// 5. Success: Redirect back to the home page with no error
		log.Println("Source added successfully")
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func AddTransactionHandler(db *pgx.Conn, tmpl *template.Template) http.HandlerFunc {
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

		transactions, err := repository.GetAllTransactions(db)
		if err != nil {
			http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
			return
		}

		
		balance, monthIncome, monthExpense, err := repository.GetSummary(db)
		if err != nil {
			http.Error(w, "Failed to fetch balance", http.StatusInternalServerError)
			return
		}
		limit := 5
		if len(transactions) < limit {
			limit = len(transactions)
		}
		limitedTransactions := transactions[:limit]

		var TransPopup = false
		var sourcePopup = false

		if r.URL.Query().Get("show_all_transactions") == "true" {
			TransPopup = true
		}
		var AllSources []model.Account
		if r.URL.Query().Get("show_all_sources") == "true" {
			sourcePopup = true
			sourcesWithBalance, err := repository.GetAllSources(db)
			if err != nil {
				http.Error(w, "Failed to fetch source balances", http.StatusInternalServerError)
				return
			}
			AllSources = sourcesWithBalance
		}

		sources, err := repository.GetAllSoucesName(db)
		if err != nil {
			http.Error(w, "Failed to fetch sources", http.StatusInternalServerError)
			return
		}




		var req model.AddTransactionRequest

		err = decoder.Decode(&req, r.PostForm)
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
		if errors.Is(err, repository.ErrNotEnoughBalance) {
			log.Println("Insufficient balance, re-rendering page with error...")



			response := model.PageData{
				Balance:          balance,
				MonthIncome:      monthIncome,
				MonthExpense:     monthExpense,
				Transactions:     limitedTransactions,
				FormErrors:       make(map[string]string),
				ShowTransPopup:   TransPopup,
				AllTransactions:  transactions,
				AvailableSources: sources,
				ShowSourcesPopup: sourcePopup,
				AllSources:       AllSources,
			}




			response.FormErrors["not_enough_balance"] = err.Error()
			err = tmpl.ExecuteTemplate(w, "home.html", response)
			if err != nil {
				log.Printf("Failed to render template: %v", err)
			}
			return


		} else if errors.Is(err,repository.ErrNegativeAmount) {
			log.Println("Negative amount, re-rendering page with error...")


			response := model.PageData{
				Balance:          balance,
				MonthIncome:      monthIncome,
				MonthExpense:     monthExpense,
				Transactions:     limitedTransactions,
				FormErrors:       make(map[string]string),
				ShowTransPopup:   TransPopup,
				AllTransactions:  transactions,
				AvailableSources: sources,
				ShowSourcesPopup: sourcePopup,
				AllSources:       AllSources,
			}

			response.FormErrors["negative_amount"] = err.Error()



			err = tmpl.ExecuteTemplate(w, "home.html", response)
			if err != nil {
				log.Printf("Failed to render template: %v", err)
			}

			return
		} else if err != nil {
			http.Error(w, "An internal server error occurred", http.StatusInternalServerError)
			return
		}



		log.Println("Transaction added successfully, redirecting.")
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func GetSummaryHandler(db *pgx.Conn, tmpl *template.Template) http.HandlerFunc {
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

		
		balance, monthIncome, monthExpense, err := repository.GetSummary(db)
		if err != nil {
			http.Error(w, "Failed to fetch balance", http.StatusInternalServerError)
			return
		}
		limit := 5
		if len(transactions) < limit {
			limit = len(transactions)
		}
		limitedTransactions := transactions[:limit]


		var TransPopup = false
		var sourcePopup = false

		if r.URL.Query().Get("show_all_transactions") == "true" {
			TransPopup = true
		}
		var AllSources []model.Account
		if r.URL.Query().Get("show_all_sources") == "true" {
			sourcePopup = true
			sourcesWithBalance, err := repository.GetAllSources(db)
			if err != nil {
				http.Error(w, "Failed to fetch source balances", http.StatusInternalServerError)
				return
			}
			AllSources = sourcesWithBalance
		}

		sources, err := repository.GetAllSoucesName(db)
		if err != nil {
			http.Error(w, "Failed to fetch sources", http.StatusInternalServerError)
			return
		}

		formErrors := make(map[string]string)
		errorKey := r.URL.Query().Get("error")
		switch errorKey {
			case "source_already_exist":
				formErrors["source_name"] = "This source already exists. Please choose another."
			case "negative_balance":
				formErrors["balance"] = "Initial balance cannot be a negative number."
		}

		response := model.PageData{
			Balance:          balance,
			MonthIncome:      monthIncome,
			MonthExpense:     monthExpense,
			Transactions:     limitedTransactions,
			FormErrors:       formErrors,
			ShowTransPopup:   TransPopup,
			AllTransactions:  transactions,
			AvailableSources: sources,
			ShowSourcesPopup: sourcePopup,
			AllSources:       AllSources,
		}

		err = tmpl.ExecuteTemplate(w, "home.html", response)
		if err != nil {
			log.Printf("Failed to render template: %v", err)
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

func DeleteTransactionsHandler(db *pgx.Conn) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        if err := r.ParseForm(); err != nil {
            http.Error(w, "Failed to parse form", http.StatusBadRequest)
            return
        }

        idStrings := r.PostForm["trans_id"]
        if len(idStrings) == 0 {
            http.Redirect(w, r, "/home?show_all_transactions=true", http.StatusSeeOther)
            return
        }

        var idsToDelete []uuid.UUID
        for _, idStr := range idStrings {
            id, err := uuid.Parse(idStr)
            if err != nil {
                http.Error(w, "Invalid transaction ID found", http.StatusBadRequest)
                return
            }
            idsToDelete = append(idsToDelete, id)
        }

        _, err := repository.DeleteTransactionsByIDs(db, idsToDelete)
        if err != nil {
            log.Printf("Failed to delete transactions: %v", err)
            http.Error(w, "Failed to delete transactions", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/home?show_all_transactions=true", http.StatusSeeOther)
    }
}
func InactiveSoucesHandler(db *pgx.Conn) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        if err := r.ParseForm(); err != nil {
            http.Error(w, "Failed to parse form", http.StatusBadRequest)
            return
        }
		log.Println("handler")
        names := r.PostForm["source_name"]
        if len(names) == 0 {
            http.Redirect(w, r, "/home?show_all_sources=true", http.StatusSeeOther)
            return
        }
        _, err := repository.InactiveSources(db, names)
        if err != nil {
            log.Printf("Failed to delete transactions: %v", err)
            http.Error(w, "Failed to delete transactions", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/home?show_all_sources=true", http.StatusSeeOther)
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
