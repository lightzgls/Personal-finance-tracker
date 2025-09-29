package main

import (
	"context"
	"finance-tracker/database"
	"finance-tracker/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	//establish db connection
	db, err := database.GetConn()
	if err != nil {
		log.Fatalf("Cant Initialize a connection to Database: %v\n", err)
	}
	log.Println("Successful connect to the database...")
	defer func() {
		db.Close(context.Background())
		log.Println("Database connection closes.")
	}()

	//start the server
	log.Println("Server is starting on http://localhost:5500/home")
	

	fmt.Println("Homepage: http://localhost:5500/home")
	fmt.Println("Transaction History: http://localhost:5500/Transactions")
	fmt.Println("Balance: http://localhost:5500/Balances")

	http.HandleFunc(("/home"), handler.GetSummaryHandler(db))
	http.HandleFunc(("/Transactions"), handler.GetAllTransactionsHandler(db))
	http.HandleFunc(("/Balances"), handler.GetAllSourcesHandler(db))

	err = http.ListenAndServe(":5500", nil)
	if err != nil {
		log.Fatal("ListenandServe: ", err)
	}
}
