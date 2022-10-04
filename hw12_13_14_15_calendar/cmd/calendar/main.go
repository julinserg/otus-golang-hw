package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
}

//goose -dir migrations postgres "user=sergey password=sergey dbname=calendar sslmode=disable" up
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

	f, err := os.OpenFile("calendar_logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("error opening file: " + err.Error())
	}
	defer f.Close()

	logg := logger.New(config.Logger.Level, f)

	var storage app.Storage
	if config.Storage.IsInMemory {
		fmt.Println("use inmemory")
		storage = memorystorage.New()
	} else {
		fmt.Println("use psql")
		sqlstor := sqlstorage.New()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := sqlstor.Connect(ctx, config.PSQL.DSN); err != nil {
			logg.Error("cannot connect to psql: " + err.Error())
			os.Exit(1)
		}
		defer func() {
			if err := sqlstor.Close(ctx); err != nil {
				logg.Error("cannot close psql connection: " + err.Error())
			}
		}()
		storage = sqlstor
	}

	calendar := app.New(logg, storage)

	endpoint := net.JoinHostPort(config.Http.Host, config.Http.Port)
	server := internalhttp.NewServer(logg, calendar, endpoint)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
