package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)
func GetConn() (*pgx.Conn,error){
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	return conn,err
}
