package store

import (
	"context"
	"database/sql"
)

type Driver interface {
	GetDB() *sql.DB
	Close() error

	FindMigrationHistoryList(ctx context.Context, find *FindMigrationHistory) ([]*MigrationHistory, error)
	UpsertMigrationHistory(ctx context.Context, upsert *UpsertMigrationHistory) (*MigrationHistory, error)

	// user service
	CreateUser(ctx context.Context, create *User) (*User, error)
	UpdateUser(ctx context.Context, update *UpdateUser) (*User, error)
	ListUsers(ctx context.Context, find *FindUser) ([]*User, error)
	DeleteUser(ctx context.Context, delete *DeleteUser) error

	UpsertUserSetting(ctx context.Context, upsert *UserSetting) (*UserSetting, error)
	ListUserSettings(ctx context.Context, find *FindUserSetting) ([]*UserSetting, error)
	DeleteUserSetting(ctx context.Context, delete *DeleteUserSetting) error

	// libro service
	CreateBook(ctx context.Context, create *Book) (*Book, error)
	UpdateBook(ctx context.Context, update *UpdateBook) (*Book, error)
	ListBooks(ctx context.Context, find *FindBook) ([]*Book, error)
	DeleteBook(ctx context.Context, delete *DeleteBook) error
	
	CreateBookReview(ctx context.Context, create *BookReview) (*BookReview, error)
	UpdateBookReview(ctx context.Context, update *UpdateBookReview) (*BookReview, error)
	ListBookReviews(ctx context.Context, find *FindBookReview) ([]*BookReview, error)
	DeleteBookReview(ctx context.Context, delete *DeleteBookReview) error

	ListBooksReadInYear(ctx context.Context, userID int32, year string) ([]*BookRead, error)
	ReportBook(ctx context.Context, userID int32) ([]*ReportBook, error)

	CreateDineroCategory(ctx context.Context, create *DineroCategory) (*DineroCategory, error)
	UpdateDineroCategory(ctx context.Context, update *UpdateDineroCategory) (*DineroCategory, error)
	ListDineroCategories(ctx context.Context, find *FindDineroCategory) ([]*DineroCategory, error)
	DeleteDineroCategory(ctx context.Context, delete *DeleteDineroCategory) error
}
