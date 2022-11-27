package app_calendar_sender

import (
	"context"
	"encoding/json"
	"fmt"

	amqp_pub "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/amqp/pub"
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
	logger           Logger
	pub              amqp_pub.AmqpPub
	uri              string
	consumer         string
	queue            string
	exchange         string
	exchangeType     string
	routingKey       string
	exchangeUser     string
	exchangeUserType string
}

func New(logger Logger, uri string, consumer string,
	queue string, exchange string, exchangeType string,
	routingKey string, exchangeUser string, exchangeUserType string) *AppCalendarSender {
	return &AppCalendarSender{
		logger:           logger,
		pub:              *amqp_pub.New(logger),
		uri:              uri,
		consumer:         consumer,
		queue:            queue,
		exchange:         exchange,
		exchangeType:     exchangeType,
		routingKey:       routingKey,
		exchangeUser:     exchangeUser,
		exchangeUserType: exchangeUserType,
	}
}

func (a *AppCalendarSender) Start(ctx context.Context) error {
	conn, err := amqp.Dial(a.uri)
	if err != nil {
		return err
	}

	c := amqp_sub.New(a.consumer, conn, a.logger)
	msgs, err := c.Consume(ctx, a.queue, a.exchange, a.exchangeType, a.routingKey)
	if err != nil {
		return err
	}

	err = a.pub.CreateExchange(a.uri, a.exchangeUser, a.exchangeUserType)
	if err != nil {
		return err
	}

	a.logger.Info("start consuming...")

	for m := range msgs {
		notifyEvent := app.NotifyEvent{}
		json.Unmarshal(m.Data, &notifyEvent)
		if err != nil {
			return err
		}
		a.logger.Info(fmt.Sprintf("receive new message:%+v\n", notifyEvent))

		if err := a.pub.Publish(a.uri, a.exchangeUser, a.exchangeUserType, "", string(m.Data), true); err != nil {
			return err
		}
		a.logger.Info("notified event for user queue is OK ( EventId: " + notifyEvent.ID + ")")
	}
	return nil
}
