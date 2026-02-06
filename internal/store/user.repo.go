package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Role string

const (
	Teacher Role = "teacher"
	Student Role = "student"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type userRepo struct {
	db *sql.DB
}

func (r *userRepo) Create(ctx context.Context, name, email, password string, role Role) (*User, error) {
	query := "INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, name, email, password, role, created_at"
	var user User

	err := r.db.QueryRowContext(ctx, query, name, email, password, role).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &user, nil 
}
