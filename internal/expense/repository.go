package expense

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, exp *Expense) error
	GetAll(ctx context.Context) ([]Expense, error)
}

// PostgresRepository now holds the real database connection pool
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository injects the DB pool into the repo
func NewPostgresRepository(dbPool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: dbPool,
	}
}

// Create executes an actual SQL query against your live database
func (r *PostgresRepository) Create(ctx context.Context, exp *Expense) error {
	query := `
		INSERT INTO expenses (title, amount, created_at) 
		VALUES ($1, $2, NOW()) 
		RETURNING id, created_at;
	`

	// QueryRow executes the query and scans the database-generated ID and timestamp
	err := r.db.QueryRow(ctx, query, exp.Title, exp.Amount).Scan(&exp.ID, &exp.CreatedAt)
	if err != nil {
		return err // Pass the database error back up to the service layer
	}

	return nil
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]Expense, error) {
	query := `SELECT id, title, amount, created_at FROM expenses ORDER BY created_at DESC;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var exp Expense
		err := rows.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.CreatedAt)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, exp)
	}

	return expenses, nil
}
