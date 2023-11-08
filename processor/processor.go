package processor

import (
	"context"
	"strconv"

	"github.com/benthosdev/benthos/v4/public/service"
)

var testConfigSpec = service.NewConfigSpec().
	Summary("Creates a processor that transforms OSM way to Platform Track.")

func init() {
	constructor := func(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {
		return newTestProcessor(conf, mgr.Logger())
	}

	err := service.RegisterProcessor("processor_test", testConfigSpec, constructor)
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------

type testProcessor struct {
	log *service.Logger
}

func newTestProcessor(conf *service.ParsedConfig, log *service.Logger) (*testProcessor, error) {
	return &testProcessor{
		log: log,
	}, nil
}

func (wtt *testProcessor) Process(ctx context.Context, m *service.Message) (service.MessageBatch, error) {
	var msgs service.MessageBatch
	for i := 0; i < 3; i++ {
		msg := service.NewMessage([]byte(strconv.Itoa(i)))
		msgs = append(msgs, msg)
	}
	wtt.log.Info("custom processor")

	return msgs, nil
}

func (wtt *testProcessor) Close(ctx context.Context) error {
	return nil
}
