package postgresqldb

import (
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	openConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_open_connections",
		Help: "Number of open database connections.",
	})
	inUseConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_in_use_connections",
		Help: "Number of in-use (active) database connections.",
	})
	idleConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_idle_connections",
		Help: "Number of idle database connections.",
	})
	waitCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_wait_count",
		Help: "Number of requests that had to wait for a connection.",
	})
	waitDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "db_wait_duration_seconds",
		Help: "Total time waiting for a connection in seconds.",
	})
)

func init() {
	prometheus.MustRegister(openConnections)
	prometheus.MustRegister(inUseConnections)
	prometheus.MustRegister(idleConnections)
	prometheus.MustRegister(waitCount)
	prometheus.MustRegister(waitDuration)
}

func recordDBStats(db *sql.DB) {
	stats := db.Stats()
	openConnections.Set(float64(stats.OpenConnections))
	inUseConnections.Set(float64(stats.InUse))
	idleConnections.Set(float64(stats.Idle))
	waitCount.Set(float64(stats.WaitCount))
	waitDuration.Set(stats.WaitDuration.Seconds())
}

func StartMonitoring(db *sql.DB, interval time.Duration) {
	go func() {
		for {
			recordDBStats(db)
			time.Sleep(interval)
		}
	}()
}
