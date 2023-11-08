package main

import (
	"context"

	_ "benthos-test/input"
	_ "benthos-test/output"
	_ "benthos-test/processor"

	_ "github.com/benthosdev/benthos/v4/public/components/all"
	"github.com/benthosdev/benthos/v4/public/service"
)

func main() {
	service.RunCLI(context.Background())
}
