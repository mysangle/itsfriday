package store

import (
	"context"
)

type Book struct {
	ID         int32
    CreatedTs  int64

	UserID     int32
	Title      string
	Author     string
	Translator string
	Pages      int32
	PubYear    int32
	Genre      string
}

type UpdateBook struct {
    ID         int32

	UserID     *int32
	Title      *string
	Author     *string
	Translator *string
	Pages      *int32
	PubYear    *int32
	Genre      *string
}

type FindBook struct {
	ID         *int32

	UserID     *int32
	Title      *string
	Author     *string

	// The maximum number of books to return.
	Limit      *int
}

type DeleteBook struct {
    ID         int32
}

type BookReview struct {
	ID         int32
	CreatedTs  int64

    UserID     int32
	BookID     int32
    DateRead   string
	Rating     float32
	Review     string
}

type UpdateBookReview struct {
	ID         int32

    UserID     *int32
	BookID     *int32
    DateRead   *string
	Rating     *float32
	Review     *string
}

type FindBookReview struct {
    ID         *int32

	UserID     *int32
	BookID     *int32
    DateRead   *string
	Rating     *float32

	// The maximum number of books to return.
	Limit      *int
}

type DeleteBookReview struct {
    ID         int32
}

type BookRead struct {
	BookID     int32
	ReviewID   int32

	Title      string
	Author     string
	Translator string
	Pages      int32
	PubYear    int32
	Genre      string

	DateRead   string
	Rating     float32
	Review     string
}

type ReportBook struct {
	Year       int32
	Count      int32
}

func (s *Store) CreateBook(ctx context.Context, create *Book) (*Book, error) {
    book, err := s.driver.CreateBook(ctx, create)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *Store) UpdateBook(ctx context.Context, update *UpdateBook) (*Book, error) {
    book, err := s.driver.UpdateBook(ctx, update)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *Store) ListBooks(ctx context.Context, find *FindBook) ([]*Book, error) {
    list, err := s.driver.ListBooks(ctx, find)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Store) GetBook(ctx context.Context, find *FindBook) (*Book, error) {
	list, err := s.ListBooks(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}

	book := list[0]
	return book, nil
}

func (s *Store) DeleteBook(ctx context.Context, delete *DeleteBook) error {
    err := s.driver.DeleteBook(ctx, delete)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateBookReview(ctx context.Context, create *BookReview) (*BookReview, error) {
    review, err := s.driver.CreateBookReview(ctx, create)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (s *Store) UpdateBookReview(ctx context.Context, update *UpdateBookReview) (*BookReview, error) {
    review, err := s.driver.UpdateBookReview(ctx, update)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (s *Store) ListBookReviews(ctx context.Context, find *FindBookReview) ([]*BookReview, error) {
    list, err := s.driver.ListBookReviews(ctx, find)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Store) GetBookReview(ctx context.Context, find *FindBookReview) (*BookReview, error) {
	list, err := s.ListBookReviews(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}

	book := list[0]
	return book, nil
}

func (s *Store) DeleteBookReview(ctx context.Context, delete *DeleteBookReview) error {
    err := s.driver.DeleteBookReview(ctx, delete)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) ListBooksReadInYear(ctx context.Context, userID int32, year int32) ([]*BookRead, error) {
	list, err := s.driver.ListBooksReadInYear(ctx, userID, year)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Store) ReportBook(ctx context.Context, userID int32) ([]*ReportBook, error) {
	list, err := s.driver.ReportBook(ctx, userID)
	if err != nil {
		return nil, err
	}

	return list, nil
}
