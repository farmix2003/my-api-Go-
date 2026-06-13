package expense

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateExpense(ctx context.Context, exp *Expense) error {
	exp.Title = strings.TrimSpace(exp.Title)
	if exp.Title == "" {
		return errors.New("title is required")
	}

	if exp.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	return s.repo.Create(ctx, exp)
}

func (s *Service) GetAllExpenses(ctx context.Context) ([]Expense, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateExpense(ctx context.Context, id int, exp *Expense) error {
	return s.repo.Update(ctx, id, exp)
}
