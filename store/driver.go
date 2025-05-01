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

	CreateUser(ctx context.Context, create *User) (*User, error)
	UpdateUser(ctx context.Context, update *UpdateUser) (*User, error)
	ListUsers(ctx context.Context, find *FindUser) ([]*User, error)
	DeleteUser(ctx context.Context, delete *DeleteUser) error

	UpsertUserSetting(ctx context.Context, upsert *UserSetting) (*UserSetting, error)
	ListUserSettings(ctx context.Context, find *FindUserSetting) ([]*UserSetting, error)
	DeleteUserSetting(ctx context.Context, delete *DeleteUserSetting) (error)
}
