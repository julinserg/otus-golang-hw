package app_calendar_scheduler

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	amqp_pub "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/amqp/pub"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage interface {
	GetEventsForNotify(timeNow time.Time) ([]storage.Event, error)
	MarkEventIsNotifyed(id string) error
	RemoveOldYearEvent(timeLimit time.Time) (int64, error)
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type AppCalendarScheduler struct {
	logger                 Logger
	storage                Storage
	pub                    amqp_pub.AmqpPub
	uri                    string
	exchange               string
	exchangeType           string
	key                    string
	timeoutCheckNeedNotify int
	timeoutCheckNeedRemove int
}

func New(logger Logger, storage Storage,
	uri string, exchange string, exchangeType string,
	key string, timeoutCheckNeedNotify int, timeoutCheckNeedRemove int,
) *AppCalendarScheduler {
	return &AppCalendarScheduler{
		logger:                 logger,
		storage:                storage,
		pub:                    *amqp_pub.New(logger),
		uri:                    uri,
		exchange:               exchange,
		exchangeType:           exchangeType,
		key:                    key,
		timeoutCheckNeedNotify: timeoutCheckNeedNotify,
		timeoutCheckNeedRemove: timeoutCheckNeedRemove,
	}
}

func (a *AppCalendarScheduler) sendNotify() error {
	events, err := a.storage.GetEventsForNotify(time.Now().UTC())
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}
	for _, ev := range events {
		err := a.storage.MarkEventIsNotifyed(ev.ID)
		if err != nil {
			return err
		}
	}
	for _, ev := range events {
		nev := &app.NotifyEvent{
			ID:        ev.ID,
			Title:     ev.Title,
			TimeStart: ev.TimeStart,
			UserID:    ev.UserID,
		}
		data, err := json.Marshal(nev)
		if err != nil {
			return err
		}
		if err := a.pub.Publish(a.uri, a.exchange, a.exchangeType, a.key, string(data), true); err != nil {
			return err
		}
		a.logger.Info("published OK")
	}

	return nil
}

func (a *AppCalendarScheduler) removeOldEvents() error {
	rows, err := a.storage.RemoveOldYearEvent(time.Now().UTC())
	if err != nil {
		return err
	}
	if rows > 0 {
		a.logger.Info("removed " + strconv.FormatInt(rows, 10) + " old events")
	}
	return nil
}

func (a *AppCalendarScheduler) Start(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var mainErr error
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(a.timeoutCheckNeedNotify) * time.Second):
				err := a.sendNotify()
				if err != nil {
					mainErr = err
					return
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(a.timeoutCheckNeedRemove) * time.Second):
				err := a.removeOldEvents()
				if err != nil {
					mainErr = err
					return
				}
			}
		}
	}()
	wg.Wait()
	return mainErr
}
