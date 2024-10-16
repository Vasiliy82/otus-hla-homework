package repository

import "github.com/prometheus/client_golang/prometheus"

var (
	// Счетчик для ДЗ № 3 (необходимо считать количество успешно записанных данных)
	recordsCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "db_records_count",
		Help: "Number of records created ",
	})
)

func init() {
	prometheus.MustRegister(recordsCount)
}

func incRecordsCount() {
	recordsCount.Add(1.0)
}
