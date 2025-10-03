package main

import (
	"context"
	"finance-tracker/database"
	"finance-tracker/handler"
	"fmt"
	"log"
	"net/http"
	"html/template"
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
	log.Println("Server is starting on http://localhost:8080/home")
	fmt.Println("Homepage: http://localhost:8080/home")
	
	templates := template.Must(template.ParseFiles("templates/home.html"));



	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
	})
	http.HandleFunc(("/home"), handler.GetSummaryHandler(db,templates))
	http.HandleFunc(("/Balances"), handler.GetAllSourcesHandler(db))
	http.HandleFunc(("/AddTransaction"), handler.AddTransactionHandler(db,templates))
	http.HandleFunc(("/AddSource"), handler.AddSourceHandler(db))
	http.HandleFunc(("/delete-transactions"), handler.DeleteTransactionsHandler(db))
	http.HandleFunc(("/delete-sources"),handler.InactiveSoucesHandler(db))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
