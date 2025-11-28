package storage

import (
	"time"

	"github.com/HelplessPlacebo/backend/pkg/shared"
	"github.com/jmoiron/sqlx"
)

type RefreshTokenRepo struct {
	db *sqlx.DB
}

func NewRefreshRepo(db *sqlx.DB) *RefreshTokenRepo { return &RefreshTokenRepo{db: db} }

func (r *RefreshTokenRepo) SaveHash(hashedToken string, userID int, expiresAt time.Time) error {
	_, err := r.db.Exec(`INSERT INTO refresh_tokens (token_hash, user_id, expires_at)
		VALUES ($1,$2,$3)`, hashedToken, userID, expiresAt)
	if err != nil {
		return shared.Internal("failed to save refresh token", err)
	}
	return nil
}

func (r *RefreshTokenRepo) DeleteHashed(hashed string) error {

	_, err := r.db.Exec(`
		DELETE FROM refresh_tokens WHERE token_hash=$1
	`, hashed)

	if err != nil {
		return shared.Internal("failed to delete refresh token", err)
	}
	return nil
}

func (r *RefreshTokenRepo) Delete(raw string) error {
	hash := shared.HashSHA256(raw)

	_, err := r.db.Exec(`
		DELETE FROM refresh_tokens WHERE token_hash=$1
	`, hash)

	if err != nil {
		return shared.Internal("failed to delete refresh token", err)
	}
	return nil
}

func (r *RefreshTokenRepo) Find(raw string) (int, time.Time, *shared.AppError) {
	hash := shared.HashSHA256(raw)

	var userID int
	var exp time.Time

	row := r.db.QueryRowx(`
		SELECT user_id, expires_at 
		FROM refresh_tokens 
		WHERE token_hash=$1
	`, hash)

	if err := row.Scan(&userID, &exp); err != nil {
		return 0, time.Time{}, shared.NotFound("refresh token not found", err)
	}

	return userID, exp, nil
}

func (r *RefreshTokenRepo) FindByUserID(userID int) (string, time.Time, error) {
	var hash string
	var exp time.Time

	err := r.db.QueryRow(`
		SELECT token_hash, expires_at
		FROM refresh_tokens
		WHERE user_id=$1
		ORDER BY expires_at DESC
		LIMIT 1
	`, userID).Scan(&hash, &exp)

	return hash, exp, err
}
