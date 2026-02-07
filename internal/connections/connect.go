package connections

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	l "online-subscription-service/pkg/logger"
)

type ConnectionPool struct {
	l l.Logger
}

func NewConnectionPool(l l.Logger) *ConnectionPool {
	return &ConnectionPool{l: l}
}

func (c *ConnectionPool) InitPool(ctx context.Context) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(c.getConnectionString())
	if err != nil {
		c.l.Error("error to parse pool config", "error", err.Error())
	}
	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		c.l.Error("failed connect to pool", "error", err.Error())
	}

	for range 6 {
		err = pool.Ping(ctx)
		if err == nil {
			return pool
		}
		time.Sleep(time.Second)
	}

	if err != nil {
		c.l.Error("failed ping database:", "error", err.Error())
	}
	return pool
}

func (c *ConnectionPool) getConnectionString() string {
	conn := os.Getenv("CONNECTION_STRING")
	if conn == "" {
		c.l.Error("failed loading .env connection string")
	}
	return conn
}
