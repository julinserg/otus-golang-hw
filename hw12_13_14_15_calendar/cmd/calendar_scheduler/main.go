package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app_calendar_scheduler"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/scheduler_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	err := config.Read(configFile)
	if err != nil {
		log.Fatalln("failed to read config: " + err.Error())
	}

	f, err := os.OpenFile("calendar_scheduler_logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalln("error opening file: " + err.Error())
	}
	defer f.Close()

	logg := logger.New(config.Logger.Level, f)

	sqlstor := sqlstorage.New()
	ctxDB, cancelDB := context.WithCancel(context.Background())
	defer cancelDB()
	if err := sqlstor.Connect(ctxDB, config.PSQL.DSN); err != nil {
		logg.Error("cannot connect to psql: " + err.Error())
		return
	}
	defer func() {
		if err := sqlstor.Close(ctxDB); err != nil {
			logg.Error("cannot close psql connection: " + err.Error())
		}
	}()

	calendarScheduler := app_calendar_scheduler.New(logg, sqlstor, config.AMQP.URI, config.AMQP.Exchange,
		config.AMQP.ExchangeType, config.AMQP.Key, config.Scheduler.TimeoutCheckNotify, config.Scheduler.TimeoutCheckRemove)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar_scheduler is running...")

	if err := calendarScheduler.Start(ctx); err != nil {
		logg.Error("calendar_scheduler failed: " + err.Error())
		cancel()
		return
	}

}
