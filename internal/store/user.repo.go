package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

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

type Students struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
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

func (r *userRepo) GetUser(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, name, email, password, created_at, role FROM users WHERE email = $1`

	var user User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}


func (r *userRepo) Me(ctx context.Context, id string) (*User, error){
	query := `SELECT id, name, email, role FROM users WHERE id = $1`

	var user User 
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		log.Printf("error on ME: %v", err)
		return nil, fmt.Errorf("faliled to get user me: %w", err)
	}

	return &user, nil 
}

func (r *userRepo) GetAllStudents(ctx context.Context) ([]*Students, error) {
	query := `
		SELECT id, name, email
		FROM users
		WHERE role = 'student'
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*Students

	for rows.Next() {
		var u Students
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
		)
		if err != nil {
			return nil, err
		}

		students = append(students, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}
