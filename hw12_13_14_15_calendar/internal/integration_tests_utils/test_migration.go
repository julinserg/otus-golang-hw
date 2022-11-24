//go:build integration

package integration_tests_utils

import (
	"context"
	"log"

	// Register pgx driver for postgresql.
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/jmoiron/sqlx"
)

var schemaDropAndCreate = `
DROP table if exists events;
CREATE table events (
    id              text primary key,
    title           text not null,
    time_start      timestamp not null,
    time_stop       timestamp not null,
    description     text,
    user_id         text not null,    
    time_notify     bigint,
	is_notifyed     boolean,
	CONSTRAINT time_start_unique UNIQUE (time_start)
);`

func execSql(sql string, dsn string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbTestConnect, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	defer dbTestConnect.Close()

	err = dbTestConnect.PingContext(ctx)
	if err != nil {
		log.Fatal("cannot ping to db:", err)
	}

	dbTestConnect.MustExec(sql)
}

func DropAndCreateSchema(dsn string) {
	execSql(schemaDropAndCreate, dsn)
}
