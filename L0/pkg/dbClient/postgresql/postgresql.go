package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context) (pool *pgxpool.Pool, err error) {
	var (
		host     = viper.GetViper().GetString("postgreConfig.host")
		port     = viper.GetViper().GetString("postgreConfig.port")
		user     = viper.GetViper().GetString("postgreConfig.user")
		password = viper.GetViper().GetString("postgreConfig.password")
		dbname   = viper.GetViper().GetString("postgreConfig.dbname")
	)

	connstr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	ctx, cancel := context.WithTimeout(ctx, viper.GetViper().GetDuration("postgreConfig.connectionTimeout")*time.Second)
	defer cancel()

	pool, err = pgxpool.Connect(ctx, connstr)
	if err != nil {
		log.Fatal("Couldn't connect to postgres")
		return nil, err
	}
	return pool, nil
}
