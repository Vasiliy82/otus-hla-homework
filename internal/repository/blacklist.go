package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type blacklistRepository struct {
	db *sql.DB
}

func NewBlacklistRepository(db *sql.DB) *blacklistRepository {
	return &blacklistRepository{db: db}
}

func (r *blacklistRepository) AddToBlacklist(serial string, expireDate time.Time) error {

	if _, err := r.db.Exec("INSERT INTO blacklisted (serial, expire_date) VALUES($1, $2)", serial, expireDate); err != nil {
		return fmt.Errorf("BlackListRepository.AddToBlacklist: r.db.Exec returned error: %w", err)
	}

	return nil

}

func (r *blacklistRepository) IsBlacklisted(serial string) (bool, error) {
	var result int

	err := r.db.QueryRow("SELECT 1 FROM blacklisted WHERE serial = $1", serial).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 0 rows is not error
			return false, nil
		}
		return false, fmt.Errorf("blacklistRepository.IsBlacklisted: r.db.QueryRow returned error %w", err)
	}
	return true, nil
}

func (r *blacklistRepository) NewSerial() (int64, error) {

	var result int64

	err := r.db.QueryRow("SELECT nextval('jwt_token')").Scan(&result)
	if err != nil {
		return 0, fmt.Errorf("blacklistRepository.NewSerial: r.db.QueryRow returned error %w", err)
	}
	return result, nil
}
