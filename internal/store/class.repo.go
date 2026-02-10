package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Class struct {
	ID         string         `json:"id"`
	ClassName  string         `json:"class_name"`
	TeacherId  string         `json:"teacher_id"`
	StudentIds pq.StringArray `json:"student_ids"`
	CreatedAt  time.Time      `json:"created_at"`
}

type classRepo struct {
	db *sql.DB
}

func (r *classRepo) Create(ctx context.Context, className, teacherId string) (*Class, error) {
	query := `INSERT INTO classes (class_name, teacher_id) VALUES ($1, $2) RETURNING id, class_name, teacher_id, student_ids, created_at`

	var class Class

	err := r.db.QueryRowContext(ctx, query, className, teacherId).Scan(
		&class.ID,
		&class.ClassName,
		&class.TeacherId,
		&class.StudentIds,
		&class.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &class, nil

}
