package storage

import (
	"database/sql"
	"strings"
	"time"

	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type User struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Name         string    `db:"name"`
	CreatedAt    time.Time `db:"created_at"`
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(email, passHash, name string) error {
	_, err := r.db.Exec(`
		INSERT INTO users (email, password_hash, name)
		VALUES ($1, $2, $3)
	`, email, passHash, name)

	if err != nil {
		if pe, ok := err.(*pq.Error); ok {
			if pe.Code == "23505" { // unique_violation
				return shared.Conflict("email already exists", err)
			}
		}

		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return shared.Conflict("email already exists", err)
		}
		return shared.Internal("failed to create user", err)
	}

	return nil
}

func (r *UserRepo) GetByEmail(email string) (*User, error) {
	var u User
	err := r.db.Get(&u, `SELECT * FROM users WHERE email=$1`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.NotFound("user not found", err)
		}
		return nil, shared.Internal("db error", err)
	}
	return &u, nil
}
