package app_calendar_sender

import (
	"context"
	"encoding/json"
	"fmt"

	amqp_sub "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/amqp/sub"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/streadway/amqp"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type AppCalendarSender struct {
	logger   Logger
	uri      string
	consumer string
	queue    string
}

func New(logger Logger, uri string, consumer string, queue string) *AppCalendarSender {
	return &AppCalendarSender{logger: logger, uri: uri, consumer: consumer, queue: queue}
}

func (a *AppCalendarSender) Start(ctx context.Context) error {
	conn, err := amqp.Dial(a.uri)
	if err != nil {
		return err
	}

	c := amqp_sub.New(a.consumer, conn)
	msgs, err := c.Consume(ctx, a.queue)

	a.logger.Info("start consuming...")

	for m := range msgs {
		notifyEvent := app.NotifyEvent{}
		json.Unmarshal(m.Data, &notifyEvent)
		if err != nil {
			return err
		}
		a.logger.Info(fmt.Sprintf("receive new message:%+v\n", notifyEvent))
	}
	return nil
}
