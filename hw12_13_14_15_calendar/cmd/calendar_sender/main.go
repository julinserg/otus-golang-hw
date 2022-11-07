package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app_calendar_sender"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/sender_config.toml", "Path to configuration file")
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

	calendarSender := app_calendar_sender.New(logg, config.AMQP.URI,
		config.AMQP.Consumer, config.AMQP.Queue)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar_sender is running...")

	if err := calendarSender.Start(ctx); err != nil {
		logg.Error("failed to start calendar_sender: " + err.Error())
		cancel()
		return
	}
	logg.Info("calendar_sender is stopping...")
}
