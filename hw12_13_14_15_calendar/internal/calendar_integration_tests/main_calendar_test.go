//go:build integration

package calendar_integration_tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/jmoiron/sqlx"

	// Register pgx driver for postgresql.
	_ "github.com/jackc/pgx/v4/stdlib"
)

const delay = 5 * time.Second

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

var schemaDrop = `DROP table if exists events;`

var dsn = "host=postgres port=5432 user=sergey password=sergey dbname=calendar sslmode=disable"

func execSql(sql string) {
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

func TestMain(m *testing.M) {
	log.Printf("wait %s for service availability...", delay)
	time.Sleep(delay)
	execSql(schemaDropAndCreate)
	status := godog.TestSuite{
		Name:                "integration",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:    "pretty", // Замените на "pretty" для лучшего вывода
			Paths:     []string{"features"},
			Randomize: 0, // Последовательный порядок исполнения
		},
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	//execSql(schemaDrop)
	os.Exit(status)
}
