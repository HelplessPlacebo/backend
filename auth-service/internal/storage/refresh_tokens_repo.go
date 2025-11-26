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

func (r *RefreshTokenRepo) Delete(token string) error {
	th := HashSHA256(token)
	_, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE token_hash=$1`, th)
	if err != nil {
		return shared.Internal("failed to delete refresh token", err)
	}
	return nil
}

func (r *RefreshTokenRepo) Find(token string) (int, time.Time, *shared.AppError) {
	th := HashSHA256(token)
	var userID int
	var expires time.Time

	row := r.db.QueryRowx(`SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash=$1`, th)

	if err := row.Scan(&userID, &expires); err != nil {
		return 0, time.Time{}, shared.NotFound("refresh token not found", err)
	}

	return userID, expires, nil
}
