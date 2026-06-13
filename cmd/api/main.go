package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"my-api/internal/expense"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepositoryAdapter struct {
	*expense.PostgresRepository
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	connString := fmt.Sprintf("postgres://postgres:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

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
	mux.HandleFunc("POST /expenses", handler.CreateExpenseHandler)
	mux.HandleFunc("GET /expenses", handler.GetAllExpensesHandler)
	mux.HandleFunc("PUT /expenses/{id}", handler.UpdateExpense)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Println("Server starting on port :8080...")
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server crashed: %v", err)
		}
	case <-ctx.Done():
		log.Println("Shutdown signal received...")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped gracefully")
}
