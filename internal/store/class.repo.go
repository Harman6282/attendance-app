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

type MyAttendance struct {
	ClassID   string `json:"class_id"`
	StudentID string `json:"student_id"`
	Status    string `json:"status"`
}

type AttendanceSession struct {
	ClassID    string            `json:"class_id"`
	StartedAt  time.Time         `json:"started_at"`
	Attendance map[string]string `json:"attendance"`
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

func (r *classRepo) AddStudent(ctx context.Context, studentId, classId string) (*Class, error) {

	query := `
		UPDATE classes
		SET student_ids = (
			SELECT ARRAY(
				SELECT DISTINCT unnest(
					COALESCE(student_ids, '{}'::uuid[]) || ARRAY[$1::uuid]
				)
			)
		)
		WHERE id = $2
		RETURNING id, class_name, teacher_id, student_ids, created_at
	`

	var class Class

	err := r.db.QueryRowContext(ctx, query, studentId, classId).Scan(
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

func (r *classRepo) Get(ctx context.Context, classId string) (*Class, error) {

	query := `SELECT id, class_name, teacher_id, student_ids, created_at from classes WHERE id = $1`

	var class Class

	err := r.db.QueryRowContext(ctx, query, classId).Scan(
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

func (r *classRepo) GetMyAttendance(ctx context.Context, classId, studentId string) (*MyAttendance, error) {
	query := `
		SELECT class_id, student_id, status
		FROM attendance
		WHERE class_id = $1 AND student_id = $2
		LIMIT 1
	`

	var attendance MyAttendance

	err := r.db.QueryRowContext(ctx, query, classId, studentId).Scan(
		&attendance.ClassID,
		&attendance.StudentID,
		&attendance.Status,
	)
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *classRepo) StartAttendance(ctx context.Context, classId string) (*AttendanceSession, error) {
	getStudentsQuery := `SELECT COALESCE(student_ids, '{}'::uuid[]) FROM classes WHERE id = $1`

	var studentIDs pq.StringArray
	err := r.db.QueryRowContext(ctx, getStudentsQuery, classId).Scan(&studentIDs)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM attendance WHERE class_id = $1`, classId); err != nil {
		return nil, err
	}

	if len(studentIDs) > 0 {
		insertAttendanceQuery := `
			INSERT INTO attendance (class_id, student_id, status)
			SELECT $1, unnest($2::uuid[]), 'absent'
		`
		if _, err := tx.ExecContext(ctx, insertAttendanceQuery, classId, pq.Array([]string(studentIDs))); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	session := &AttendanceSession{
		ClassID:    classId,
		StartedAt:  time.Now().UTC(),
		Attendance: make(map[string]string, len(studentIDs)),
	}
	for _, studentID := range studentIDs {
		session.Attendance[studentID] = "absent"
	}

	return session, nil
}
