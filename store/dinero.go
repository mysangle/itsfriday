package store

import (
	"context"
)

type DineroCategory struct {
	ID         int32
	
	UserID     int32
	Name       string
	Priority   int32
}

type UpdateDineroCategory struct {
	ID         int32
	UserID     int32

	Name       *string
	Priority   *int32
}

type FindDineroCategory struct {
	ID         int32

	UserID     *int32
	Name       *string
}

type DeleteDineroCategory struct {
    ID         int32
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
