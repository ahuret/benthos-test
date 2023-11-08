package output

import (
	"context"
	"fmt"

	"github.com/benthosdev/benthos/v4/public/service"
)

func testOutputConfig() *service.ConfigSpec {
	return service.NewConfigSpec()
}

func init() {
	err := service.RegisterOutput(
		"test_output", testOutputConfig(),
		func(conf *service.ParsedConfig, mgr *service.Resources) (service.Output, int, error) {
			w, err := newTestWriter(conf, mgr)
			if err != nil {
				return nil, 0, err
			}
			return w, 64, err
		},
	)
	if err != nil {
		panic(err)
	}
}

type testWriter struct {
	log *service.Logger
}

func newTestWriter(conf *service.ParsedConfig, mgr *service.Resources) (*testWriter, error) {
	n := testWriter{
		log: mgr.Logger(),
	}
	return &n, nil
}

func (n *testWriter) Connect(ctx context.Context) error {
	return nil
}

// Write attempts to write a message.
func (n *testWriter) Write(context context.Context, msg *service.Message) error {
	n.log.Info("test output write msg")
	b, err := msg.AsBytes()
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func (n *testWriter) Close(context.Context) (err error) {
	return
}
