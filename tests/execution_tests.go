package tests

import (
	"github.com/Azer0s/Hummus/runner"
	"testing"
)

//TODO: Actually test something

// TestExecution test entrypoint
func TestExecution(t *testing.T) {
	runner.RunFile("examples/map_exists.hummus")
	runner.RunFile("examples/branching.hummus")
	runner.RunFile("examples/curry.hummus")
	runner.RunFile("examples/def_function.hummus")
	runner.RunFile("examples/def_sub.hummus")
	runner.RunFile("examples/fib.hummus")
	runner.RunFile("examples/identity.hummus")
	runner.RunFile("examples/import.hummus")
	runner.RunFile("examples/loop.hummus")
	runner.RunFile("examples/map_filter_reduce.hummus")
	runner.RunFile("examples/maps.hummus")
	runner.RunFile("examples/math.hummus")
	runner.RunFile("examples/option.hummus")
	runner.RunFile("examples/pipes.hummus")
	runner.RunFile("examples/structs.hummus")
	runner.RunFile("examples/variables.hummus")
	runner.RunFile("examples/watch.hummus")
	runner.RunFile("examples/watch2.hummus")
}
