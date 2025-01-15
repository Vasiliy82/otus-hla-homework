package postgresqldb

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

// Метрики для пула подключений
var (
	acquireCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pgxpool_acquire_count",
			Help: "Общее количество успешных получений подключений из пула",
		},
		[]string{"role"},
	)
	acquireDuration = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pgxpool_acquire_duration_seconds",
			Help: "Общее время, затраченное на получение подключений из пула в секундах",
		},
		[]string{"role"},
	)
	acquiredConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pgxpool_acquired_conns",
			Help: "Текущее количество занятых подключений в пуле",
		},
		[]string{"role"},
	)
	canceledAcquireCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pgxpool_canceled_acquire_count",
			Help: "Общее количество отмененных попыток получения подключения из пула",
		},
		[]string{"role"},
	)
	constructingConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pgxpool_constructing_conns",
			Help: "Текущее количество подключений, находящихся в процессе создания",
		},
		[]string{"role"},
	)
	emptyAcquireCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pgxpool_empty_acquire_count",
			Help: "Общее количество успешных получений подключений, ожидавших освобождения или создания ресурса из-за пустого пула",
		},
		[]string{"role"},
	)
	idleConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pgxpool_idle_conns",
			Help: "Текущее количество свободных подключений в пуле",
		},
		[]string{"role"},
	)
	maxConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pgxpool_max_conns",
			Help: "Максимальное количество подключений в пуле",
		},
		[]string{"role"},
	)
	totalConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pgxpool_total_conns",
			Help: "Общее количество подключений в пуле",
		},
		[]string{"role"},
	)
)

func init() {
	// Регистрация метрик в Prometheus
	prometheus.MustRegister(
		acquireCount,
		acquireDuration,
		acquiredConns,
		canceledAcquireCount,
		constructingConns,
		emptyAcquireCount,
		idleConns,
		maxConns,
		totalConns,
	)
}

// recordDBStats обновляет метрики для переданного пула подключений и роли
func recordDBStats(pool *pgxpool.Pool, role string) {
	stats := pool.Stat()

	acquireCount.WithLabelValues(role).Add(float64(stats.AcquireCount()))
	acquireDuration.WithLabelValues(role).Add(stats.AcquireDuration().Seconds())
	acquiredConns.WithLabelValues(role).Set(float64(stats.AcquiredConns()))
	canceledAcquireCount.WithLabelValues(role).Add(float64(stats.CanceledAcquireCount()))
	constructingConns.WithLabelValues(role).Set(float64(stats.ConstructingConns()))
	emptyAcquireCount.WithLabelValues(role).Add(float64(stats.EmptyAcquireCount()))
	idleConns.WithLabelValues(role).Set(float64(stats.IdleConns()))
	maxConns.WithLabelValues(role).Set(float64(stats.MaxConns()))
	totalConns.WithLabelValues(role).Set(float64(stats.TotalConns()))
}

// StartMonitoring запускает мониторинг метрик с указанным интервалом
func StartMonitoring(cluster *DBCluster, interval time.Duration) {
	go func() {
		for {
			recordDBStats(cluster.masterPool, "master")
			for i, replica := range cluster.replicaPools {
				role := fmt.Sprintf("replica_%d", i+1)
				recordDBStats(replica, role)
			}
			time.Sleep(interval)
		}
	}()
}
