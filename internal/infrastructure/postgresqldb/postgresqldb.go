package postgresqldb

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/Vasiliy82/otus-hla-homework/internal/config"
)

const (
	driverName               = "postgres"
	connectionStringTemplate = "postgres://%s:%s@%s:%s/%s"
)

func InitDB(ctx context.Context, cfg *config.DatabaseConfig) (*sql.DB, error) {

	connection := fmt.Sprintf(connectionStringTemplate, cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
	val := url.Values{}
	val.Add("sslmode", "disable")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
