package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/postgresqldb"
)

type blacklistRepository struct {
	ctx       context.Context
	dbCluster *postgresqldb.DBCluster
}

func NewBlacklistRepository(ctx context.Context, dbCluster *postgresqldb.DBCluster) *blacklistRepository {
	return &blacklistRepository{ctx: ctx, dbCluster: dbCluster}
}

func (r *blacklistRepository) AddToBlacklist(serial string, expireDate time.Time) error {
	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		return fmt.Errorf("blacklistRepository.AddToBlacklist: r.dbCluster.GetDB returned error %w", err)
	}

	if _, err := db.Exec(r.ctx, "INSERT INTO blacklisted (serial, expire_date) VALUES($1, $2)", serial, expireDate); err != nil {
		return fmt.Errorf("BlackListRepository.AddToBlacklist: r.db.Exec returned error: %w", err)
	}

	return nil

}

func (r *blacklistRepository) IsBlacklisted(serial string) (bool, error) {
	var result int

	db, err := r.dbCluster.GetDBPool(postgresqldb.Read)
	if err != nil {
		return false, fmt.Errorf("blacklistRepository.IsBlacklisted: r.dbCluster.GetDB returned error %w", err)
	}

	err = db.QueryRow(r.ctx, "SELECT 1 FROM blacklisted WHERE serial = $1", serial).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 0 rows is not error
			return false, nil
		}
		return false, fmt.Errorf("blacklistRepository.IsBlacklisted: r.db.QueryRow returned error %w", err)
	}
	return true, nil
}

func (r *blacklistRepository) NewSerial() (string, error) {

	var result string

	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		return "", fmt.Errorf("blacklistRepository.NewSerial: r.dbCluster.GetDB returned error %w", err)
	}

	err = db.QueryRow(r.ctx, "SELECT nextval('jwt_token')").Scan(&result)
	if err != nil {
		return "", fmt.Errorf("blacklistRepository.NewSerial: r.db.QueryRow returned error %w", err)
	}
	return result, nil
}
