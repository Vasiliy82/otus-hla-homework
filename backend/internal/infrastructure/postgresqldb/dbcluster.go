package postgresqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"sync"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"

	// _ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

const (
	driverName               = "postgres"
	connectionStringTemplate = "postgres://%s:%s@%s:%s/%s"
	passwordMasked           = "*hidden*"

	Read OperationType = iota
	ReadWrite
)

type OperationType int

type DBCluster struct {
	masterDB   *sql.DB
	replicaDBs []*sql.DB
}

// NewDBCluster создает новый кластер базы данных
func NewDBCluster(masterDB *sql.DB, replicaDBs []*sql.DB) *DBCluster {
	return &DBCluster{
		masterDB:   masterDB,
		replicaDBs: replicaDBs,
	}
}

func InitDBCluster(ctx context.Context, cfg *config.DatabaseConfig) (*DBCluster, error) {
	master, err := initDBInstance(cfg.Master)
	if err != nil {
		logger.Logger().Debugw("Error in InitDBCluster: InitDBInstance(master)", "err", err)
		return nil, fmt.Errorf("in InitDBCluster: InitDBInstance(master) returned %w", err)
	}
	var replicas []*sql.DB
	for i, rcfg := range cfg.Replicas {
		replica, err := initDBInstance(rcfg)
		if err != nil {
			logger.Logger().Warnw(fmt.Sprintf("in InitDBCluster: InitDBInstance(replica[%d])", i), "err", err)
			continue
		}
		replicas = append(replicas, replica)
	}
	return NewDBCluster(master, replicas), nil
}

// GetDB возвращает подходящий экземпляр базы данных в зависимости от типа операции
func (cluster *DBCluster) GetDB(opType OperationType) (*sql.DB, error) {
	switch opType {
	case ReadWrite:
		if cluster.masterDB == nil {
			return nil, errors.New("master DB is not configured")
		}
		return cluster.masterDB, nil
	case Read:
		if len(cluster.replicaDBs) > 0 {
			index := rand.Intn(len(cluster.replicaDBs)) // выбираем случайную реплику
			return cluster.replicaDBs[index], nil
		}
		// если нет реплик, возвращаем master для операций чтения
		if cluster.masterDB == nil {
			return nil, errors.New("master DB is not configured")
		}
		return cluster.masterDB, nil
	default:
		return nil, errors.New("unknown operation type")
	}
}

func (c *DBCluster) Close() error {
	var wg sync.WaitGroup
	errsChan := make(chan error, len(c.replicaDBs)+1)

	// Закрываем реплики асинхронно
	for _, r := range c.replicaDBs {
		wg.Add(1)
		go func(r *sql.DB) {
			defer wg.Done()
			if err := r.Close(); err != nil {
				errsChan <- err
			}
		}(r)
	}

	// Закрываем master асинхронно
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.masterDB.Close(); err != nil {
			errsChan <- err
		}
	}()

	// Ожидаем завершения всех goroutine
	wg.Wait()
	close(errsChan)

	// Собираем ошибки
	var errs []error
	for err := range errsChan {
		errs = append(errs, err)
	}

	// Возвращаем объединенные ошибки, если они есть
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func initDBInstance(cfg *config.DBInstanceConfig) (*sql.DB, error) {

	connStr := fmt.Sprintf(connectionStringTemplate, cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
	connStrMasked := fmt.Sprintf(connectionStringTemplate, cfg.User, passwordMasked, cfg.Host, cfg.Port, cfg.Name)
	val := url.Values{}
	val.Add("sslmode", "disable")
	dsn := fmt.Sprintf("%s?%s", connStr, val.Encode())
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to [%s]: %w", connStrMasked, err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database [%s]: %w", connStrMasked, err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.MaxConnIdleTime)
	db.SetConnMaxLifetime(cfg.MaxConnLifetime)

	logger.Logger().Infof("postgresqldb.initDBInstance: connection established: %s", connStrMasked)

	return db, nil
}
