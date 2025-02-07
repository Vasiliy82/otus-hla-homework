package postgresqldb

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"sync"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"

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

func InitDBCluster(ctx context.Context, cfg *config.DatabaseConfig, appName string) (*DBCluster, error) {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	masterAppName := fmt.Sprintf("%s-master", appName)
	master, err := initDBInstance(ctx, cfg.Master, masterAppName)
	if err != nil {
		log.Debugw("Error in InitDBCluster: initDBInstance(master)", "err", err)
		return nil, fmt.Errorf("in InitDBCluster: initDBInstance(master) returned %w", err)
	}
	var replicas []*pgxpool.Pool
	for i, rcfg := range cfg.Replicas {
		slaveAppName := fmt.Sprintf("%s-slave-%d", appName, i)
		replica, err := initDBInstance(ctx, rcfg, slaveAppName)
		if err != nil {
			log.Warnw(fmt.Sprintf("in InitDBCluster: initDBInstance(replica[%d])", i), "err", err)
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

func initDBInstance(ctx context.Context, cfg *config.DBInstanceConfig, appName string) (*pgxpool.Pool, error) {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
	val := url.Values{}
	val.Add("sslmode", "disable")
	val.Add("application_name", appName)
	dsn := fmt.Sprintf("%s?%s", connStr, val.Encode())

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxConns)
	poolConfig.MinConns = int32(cfg.MinConns)
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	log.Infof("connection pool established: %s", connStr)

	return pool, nil
}
