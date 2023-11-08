package input

import (
	"context"
	"log"
	"net/http"

	"github.com/benthosdev/benthos/v4/public/service"
)

var testConfigSpec = service.NewConfigSpec().
	Summary("Creates an input that setup an HTTP server, waiting for remove_write requests and extract metrics from it.").
	Field(service.NewStringField("address").Default("http://127.0.0.1:4242")).
	Field(service.NewStringField("path").Default("/receive"))

func newTestInput(conf *service.ParsedConfig, mgr *service.Resources) (service.BatchInput, error) {
	address, err := conf.FieldString("address")
	if err != nil {
		return nil, err
	}
	path, err := conf.FieldString("path")
	if err != nil {
		return nil, err
	}
	return &testInput{
		logger:    mgr.Logger(),
		address:   address,
		path:      path,
		responses: make(chan response),
	}, nil
}

func init() {
	if err := service.RegisterBatchInput(
		"test_input", testConfigSpec,
		func(conf *service.ParsedConfig, mgr *service.Resources) (service.BatchInput, error) {
			return newTestInput(conf, mgr)
		}); err != nil {
		panic(err)
	}
}

type testInput struct {
	logger    *service.Logger
	address   string
	path      string
	responses chan response
}

type response struct {
	messages service.MessageBatch
	ack      chan error
}

func (p *testInput) Connect(ctx context.Context) error {
	go func() {
		http.HandleFunc(p.path, p.handleHTTPRequest())
		if err := http.ListenAndServe(p.address, nil); err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}

func (p *testInput) ReadBatch(ctx context.Context) (service.MessageBatch, service.AckFunc, error) {
	response := <-p.responses
	p.logger.Infof("test input read batch")
	return response.messages,
		func(ctx context.Context, err error) error {
			response.ack <- err
			return nil
		},
		nil
}

func (p *testInput) Close(ctx context.Context) error {
	p.logger.Info("Test Input: closing input")
	return nil
}

func (p *testInput) handleHTTPRequest() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var mm []*service.Message
		mm = append(mm, service.NewMessage([]byte("")))
		ack := make(chan error, 1)
		p.responses <- response{messages: service.MessageBatch(mm), ack: ack}
		err := <-ack
		close(ack)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}
