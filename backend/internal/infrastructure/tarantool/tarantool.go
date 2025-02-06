package tarantool

import (
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	tar "github.com/tarantool/go-tarantool"
)

func NewTarConn(cfg config.TarantoolConfig) (*tar.Connection, error) {

	connstr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := tar.Connect(connstr, tar.Opts{User: cfg.User, Pass: cfg.Password})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Tarantool: %w", err)
	}
	return conn, nil

}
