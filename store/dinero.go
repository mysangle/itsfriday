package store

import (
	"context"
)

type DineroCategory struct {
	ID          int32
	
	UserID      int32
	Name        string
	Priority    int32
}

type UpdateDineroCategory struct {
	ID          int32

	Name        *string
	Priority    *int32
}

type FindDineroCategory struct {
	ID          *int32

	UserID      *int32
	Name        *string
}

type DeleteDineroCategory struct {
    ID          int32
}

type DineroExpense struct {
	ID          int32

	CreatedTs   int64

	UserID      int32
	CategoryID  int32
	DateUsed    string
	Item        string
	Price       int32
}

type UpdateDineroExpense struct {
	ID          int32

	CategoryID  *int32
	DateUsed    *string
	Item        *string
	Price       *int32
}

type FindDineroExpense struct {
	ID          *int32
	UserID      *int32
    CategoryID  *int32
	Year        *int32
	Month       *int32
}

type DeleteDineroExpense struct {
    ID          int32
}

func (s *Store) CreateDineroCaterory(ctx context.Context, create *DineroCategory) (*DineroCategory, error) {
    category, err := s.driver.CreateDineroCategory(ctx, create)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Store) GetDineroCategory(ctx context.Context, find *FindDineroCategory) (*DineroCategory, error) {
	list, err := s.ListDineroCategories(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}

	book := list[0]
	return book, nil
}

func (s *Store) ListDineroCategories(ctx context.Context, find *FindDineroCategory) ([]*DineroCategory, error) {
    list, err := s.driver.ListDineroCategories(ctx, find)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Store) UpdateDineroCategory(ctx context.Context, update *UpdateDineroCategory) (*DineroCategory, error) {
    category, err := s.driver.UpdateDineroCategory(ctx, update)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Store) DeleteDineroCategory(ctx context.Context, delete *DeleteDineroCategory) error {
    err := s.driver.DeleteDineroCategory(ctx, delete)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateDineroExpense(ctx context.Context, create *DineroExpense) (*DineroExpense, error) {
    category, err := s.driver.CreateDineroExpense(ctx, create)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Store) GetDineroExpense(ctx context.Context, find *FindDineroExpense) (*DineroExpense, error) {
	list, err := s.ListDineroExpenses(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}

	book := list[0]
	return book, nil
}

func (s *Store) ListDineroExpenses(ctx context.Context, find *FindDineroExpense) ([]*DineroExpense, error) {
    list, err := s.driver.ListDineroExpenses(ctx, find)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Store) UpdateDineroExpense(ctx context.Context, update *UpdateDineroExpense) (*DineroExpense, error) {
    category, err := s.driver.UpdateDineroExpense(ctx, update)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Store) DeleteDineroExpense(ctx context.Context, delete *DeleteDineroExpense) error {
    err := s.driver.DeleteDineroExpense(ctx, delete)
	if err != nil {
		return err
	}

	return nil
}
