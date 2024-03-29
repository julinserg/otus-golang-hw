//go:build integration

package calendar_integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v16"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/streadway/amqp"
)

var amqpDSN string

func init() {
	amqpDSN = "amqp://guest:guest@rabbit:5672/"
}

const (
	queueName                 = "calendar-users-tests-queue"
	notificationsExchangeName = "calendar-users-exchange"
)

type Response struct {
	Data  []app.Event `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type notifyTest struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	messages      [][]byte
	messagesMutex *sync.RWMutex
	stopSignal    chan struct{}

	responseStatusCode int
	responseBody       []byte
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (test *notifyTest) startConsuming(ctx context.Context, _ *messages.Pickle) (context.Context, error) {
	test.messages = make([][]byte, 0)
	test.messagesMutex = new(sync.RWMutex)
	test.stopSignal = make(chan struct{})

	var err error

	test.conn, err = amqp.Dial(amqpDSN)
	panicOnErr(err)

	test.ch, err = test.conn.Channel()
	panicOnErr(err)

	// Consume
	_, err = test.ch.QueueDeclare(queueName, true, false, false, false, nil)
	panicOnErr(err)

	err = test.ch.QueueBind(queueName, "", notificationsExchangeName, false, nil)
	panicOnErr(err)

	events, err := test.ch.Consume(queueName, "", false, false, false, false, nil)
	panicOnErr(err)
	go func(stop <-chan struct{}) {
		for {
			select {
			case <-stop:
				return
			case event := <-events:
				test.messagesMutex.Lock()
				test.messages = append(test.messages, event.Body)
				test.messagesMutex.Unlock()
			}
		}
	}(test.stopSignal)
	return ctx, nil
}

func (test *notifyTest) stopConsuming(ctx context.Context, _ *messages.Pickle, _ error) (context.Context, error) {
	test.stopSignal <- struct{}{}

	panicOnErr(test.ch.Close())
	panicOnErr(test.conn.Close())
	test.messages = nil
	return ctx, nil
}

func (test *notifyTest) iSendRequestTo(httpMethod, addr string) (err error) {
	var r *http.Response

	switch httpMethod {
	case http.MethodGet:
		r, err = http.Get(addr)
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return
	}
	test.responseStatusCode = r.StatusCode
	test.responseBody, err = ioutil.ReadAll(r.Body)
	return
}

func (test *notifyTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *notifyTest) theResponseShouldMatchText(text string) error {
	if string(test.responseBody) != text {
		return fmt.Errorf("unexpected text: %s != %s", test.responseBody, text)
	}
	return nil
}

func (test *notifyTest) theResponseShouldMatchError(textError string) error {
	result := Response{}
	json.Unmarshal(test.responseBody, &result)
	if result.Error.Message != textError {
		return fmt.Errorf("unexpected error: %s != %s", result.Error.Message, textError)
	}
	return nil
}

func (test *notifyTest) theResponseShouldMatchJson(jsonText string) error {
	result := Response{}
	json.Unmarshal(test.responseBody, &result)

	idsReal := make(map[string]bool)
	for _, e := range result.Data {
		idsReal[e.ID] = true
	}

	resultTest := Response{}
	json.Unmarshal([]byte(jsonText), &resultTest)

	idsTest := make(map[string]bool)
	for _, e := range resultTest.Data {
		idsTest[e.ID] = true
	}

	if !reflect.DeepEqual(idsReal, idsTest) {
		return fmt.Errorf("unexpected error: %v != %v", idsReal, idsTest)
	}

	return nil
}

func (test *notifyTest) iSendRequestToWithData(httpMethod, addr, contentType string, data *messages.PickleDocString) (err error) {
	var r *http.Response

	switch httpMethod {
	case http.MethodPost:
		replacer := strings.NewReplacer("\n", "", "\t", "")
		cleanJson := replacer.Replace(data.Content)
		r, err = http.Post(addr, contentType, bytes.NewReader([]byte(cleanJson)))
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return
	}
	test.responseStatusCode = r.StatusCode
	test.responseBody, err = ioutil.ReadAll(r.Body)
	return
}

func (test *notifyTest) iReceiveEventWithJson(jsonText string) error {
	time.Sleep(6 * time.Second) // На всякий случай ждём обработки евента

	test.messagesMutex.RLock()
	defer test.messagesMutex.RUnlock()

	resultTest := app.NotifyEvent{}
	json.Unmarshal([]byte(jsonText), &resultTest)

	for _, msg := range test.messages {
		result := app.NotifyEvent{}
		json.Unmarshal(msg, &result)
		if resultTest.ID == result.ID && resultTest.UserID == result.UserID {
			return nil
		}
	}
	return fmt.Errorf("event with text '%s' was not found in %s", jsonText, test.messages)
}

func InitializeScenario(s *godog.ScenarioContext) {
	test := new(notifyTest)

	s.Before(test.startConsuming)

	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^The response should match text "([^"]*)"$`, test.theResponseShouldMatchText)

	s.Step(`^I send "([^"]*)" request to "([^"]*)" with "([^"]*)" data:$`, test.iSendRequestToWithData)
	s.Step(`^I receive event with json:$`, test.iReceiveEventWithJson)

	s.Step(`^The response should match error "([^"]*)"$`, test.theResponseShouldMatchError)

	s.Step(`^The response should match json:$`, test.theResponseShouldMatchJson)

	s.After(test.stopConsuming)

}
