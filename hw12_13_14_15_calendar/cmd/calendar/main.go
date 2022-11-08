package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app_calendar"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
}

/*
 goose -dir migrations postgres "user=sergey password=sergey dbname=calendar sslmode=disable" up
*/

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

	f, err := os.OpenFile("calendar_logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalln("error opening file: " + err.Error())
	}
	defer f.Close()

	logg := logger.New(config.Logger.Level, f)

	var storage app.Storage
	switch config.Storage.Type {
	case "inmemory":
		fmt.Println("use inmemory")
		storage = memorystorage.New()
	case "psql":
		fmt.Println("use psql")
		sqlstor := sqlstorage.New()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := sqlstor.Connect(ctx, config.PSQL.DSN); err != nil {
			logg.Error("cannot connect to psql: " + err.Error())
			return
		}
		defer func() {
			if err := sqlstor.Close(); err != nil {
				logg.Error("cannot close psql connection: " + err.Error())
			}
		}()
		storage = sqlstor
	}

	calendar := app_calendar.New(logg, storage)

	endpointHttp := net.JoinHostPort(config.HTTP.Host, config.HTTP.Port)
	serverHttp := internalhttp.NewServer(logg, calendar, endpointHttp)
	endpointGrpc := net.JoinHostPort(config.GRPC.Host, config.GRPC.Port)
	serverGrpc := internalgrpc.NewServer(logg, calendar, endpointGrpc)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverHttp.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
		if err := serverGrpc.Stop(ctx); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := serverHttp.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			return
		}
	}()
	go func() {
		defer wg.Done()
		if err := serverGrpc.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			return
		}
	}()
	wg.Wait()
}
