package store

import (
	"context"
	"database/sql"
)

type Users interface {
	Create(ctx context.Context, name, email, password string, role Role) error
}

type Storage struct {
	Users
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &userRepo{db},
	}
}
