package main

import (
	"context"
	"finance-tracker/database"
	"log"
)

func main() {
	//establish db connection
	db, err := database.GetConn();
	if err != nil {
		log.Fatalf("Cant Initialize a connection to Database: %v\n", err)
	}
	log.Println("Successful connect to the database...")
	defer func () {
		db.Close(context.Background())
		log.Println("Database connection closes.")
	}()
		
}