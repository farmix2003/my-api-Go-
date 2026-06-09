package main

import (
	"context"
	"fmt"
	"log"
	"my-api/internal/expense"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	connString := fmt.Sprintf("postgres://postgres:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	ctx := context.Background()

	log.Println("Connecting to PostgreSQL...")

	dbPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	if err := dbPool.Ping(pingCtx); err != nil {
		log.Fatalf("Database is unreachable: %v\n", err)
	}

	log.Println("Successfully connected to the database!")

	repo := expense.NewPostgresRepository(dbPool)

	service := expense.NewService(repo)

	handler := expense.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/expenses", handler.CreateExpenseHandler)
	mux.HandleFunc("/allexpenses", handler.GetAllExpensesHandler)

	log.Println("Server starting on port :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}
