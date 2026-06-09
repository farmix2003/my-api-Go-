package expense

import (
	"context"
	"errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateExpense(ctx context.Context, exp *Expense) error {
	if exp.Amount <= 0 {
		return errors.New("Amount must be greater than 0")
	}
	return s.repo.Create(ctx, exp)
}

func (s *Service) GetAllExpenses(ctx context.Context) ([]Expense, error) {
	return s.repo.GetAll(ctx)
}
