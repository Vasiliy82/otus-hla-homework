package postgresqldb

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"sync"

	"github.com/Vasiliy82/otus-hla-homework/backend-messages/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend-messages/internal/observability/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	Read OperationType = iota
	ReadWrite
)

type OperationType int

type DBCluster struct {
	masterPool   *pgxpool.Pool
	replicaPools []*pgxpool.Pool
}

// NewDBCluster создает новый кластер базы данных
func NewDBCluster(masterPool *pgxpool.Pool, replicaPools []*pgxpool.Pool) *DBCluster {
	return &DBCluster{
		masterPool:   masterPool,
		replicaPools: replicaPools,
	}
}

func InitDBCluster(ctx context.Context, cfg *config.DatabaseConfig) (*DBCluster, error) {
	master, err := initDBInstance(ctx, cfg.Master)
	if err != nil {
		logger.Logger().Debugw("Error in InitDBCluster: initDBInstance(master)", "err", err)
		return nil, fmt.Errorf("in InitDBCluster: initDBInstance(master) returned %w", err)
	}
	var replicas []*pgxpool.Pool
	for i, rcfg := range cfg.Replicas {
		replica, err := initDBInstance(ctx, rcfg)
		if err != nil {
			logger.Logger().Warnw(fmt.Sprintf("in InitDBCluster: initDBInstance(replica[%d])", i), "err", err)
			continue
		}
		replicas = append(replicas, replica)
	}
	return NewDBCluster(master, replicas), nil
}

// GetDBPool возвращает подходящий пул базы данных в зависимости от типа операции
func (cluster *DBCluster) GetDBPool(opType OperationType) (*pgxpool.Pool, error) {
	switch opType {
	case ReadWrite:
		if cluster.masterPool == nil {
			return nil, errors.New("master DB is not configured")
		}
		return cluster.masterPool, nil
	case Read:
		if len(cluster.replicaPools) > 0 {
			index := rand.Intn(len(cluster.replicaPools)) // выбираем случайную реплику
			return cluster.replicaPools[index], nil
		}
		// если нет реплик, возвращаем master для операций чтения
		if cluster.masterPool == nil {
			return nil, errors.New("master DB is not configured")
		}
		return cluster.masterPool, nil
	default:
		return nil, errors.New("unknown operation type")
	}
}

func (c *DBCluster) Close() {
	var wg sync.WaitGroup

	// Закрываем реплики асинхронно
	for _, r := range c.replicaPools {
		wg.Add(1)
		go func(r *pgxpool.Pool) {
			defer wg.Done()
			r.Close()
		}(r)
	}

	// Закрываем master асинхронно
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.masterPool.Close()
	}()

	// Ожидаем завершения всех goroutine
	wg.Wait()
}

func initDBInstance(ctx context.Context, cfg *config.DBInstanceConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
	val := url.Values{}
	val.Add("sslmode", "disable")
	dsn := fmt.Sprintf("%s?%s", connStr, val.Encode())

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	logger.Logger().Infof("postgresqldb.initDBInstance: connection pool established: %s", connStr)

	return pool, nil
}
