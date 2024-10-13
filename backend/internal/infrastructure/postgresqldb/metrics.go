package postgresqldb

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Метрики с возможностью добавления лейблов (срезов)
var (
	openConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_open_connections",
		Help: "Number of open database connections.",
	}, []string{"role"})

	inUseConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_in_use_connections",
		Help: "Number of in-use (active) database connections.",
	}, []string{"role"})

	idleConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_idle_connections",
		Help: "Number of idle database connections.",
	}, []string{"role"})

	waitCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_wait_count",
		Help: "Number of requests that had to wait for a connection.",
	}, []string{"role"})

	waitDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_wait_duration_seconds",
		Help: "Total time waiting for a connection in seconds.",
	}, []string{"role"})
)

func init() {
	prometheus.MustRegister(openConnections)
	prometheus.MustRegister(inUseConnections)
	prometheus.MustRegister(idleConnections)
	prometheus.MustRegister(waitCount)
	prometheus.MustRegister(waitDuration)
}

// recordDBStats обновляет метрики для переданного DB и добавляет соответствующие лейблы (role: master/replica)
func recordDBStats(db *sql.DB, role string) {
	stats := db.Stats()
	openConnections.WithLabelValues(role).Set(float64(stats.OpenConnections))
	inUseConnections.WithLabelValues(role).Set(float64(stats.InUse))
	idleConnections.WithLabelValues(role).Set(float64(stats.Idle))
	waitCount.WithLabelValues(role).Set(float64(stats.WaitCount))
	waitDuration.WithLabelValues(role).Set(stats.WaitDuration.Seconds())
}

// StartMonitoring начинает мониторинг для master и реплик в DBCluster с заданным интервалом
func StartMonitoring(cluster *DBCluster, interval time.Duration) {
	go func() {
		for {
			// Обновляем метрики для master
			recordDBStats(cluster.masterDB, "master")

			// Обновляем метрики для каждой реплики
			for i, replica := range cluster.replicaDBs {
				role := fmt.Sprintf("replica_%d", i+1) // динамически создаем имя реплики
				recordDBStats(replica, role)
			}

			time.Sleep(interval)
		}
	}()
}
