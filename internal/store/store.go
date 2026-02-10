package store

import (
	"context"
	"database/sql"
)

type Users interface {
	Create(ctx context.Context, name, email, password string, role Role) (*User, error)
	GetUser(ctx context.Context, email string) (*User, error)
	Me(ctx context.Context, id string) (*User, error)
}

type Classes interface {
	Create(ctx context.Context, className, teacherId string ) (*Class, error)
}

type Storage struct {
	Users
	Classes
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &userRepo{db},
		Classes: &classRepo{db},
	}
}
