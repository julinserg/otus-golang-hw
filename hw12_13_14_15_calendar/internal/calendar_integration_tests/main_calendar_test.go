//go:build integration

package calendar_integration_tests

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/integration_tests_utils"
)

const delay = 5 * time.Second

var dsn = "host=postgres port=5432 user=sergey password=sergey dbname=calendar sslmode=disable"

func TestMain(m *testing.M) {
	log.Printf("wait %s for service availability...", delay)
	time.Sleep(delay)
	integration_tests_utils.DropAndCreateSchema(dsn)
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
