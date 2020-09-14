package runner_test

import (
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/runner"
	"path"
	"runtime"
	"testing"
)

func TestRunFile(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	interpreter.BasePath = path.Join(path.Dir(filename), "../bin/")
	filename = path.Join(path.Dir(filename), "../examples")

	runner.RunFile(path.Join(filename, "map_exists.hummus"))
	runner.RunFile(path.Join(filename, "branching.hummus"))
	runner.RunFile(path.Join(filename, "curry.hummus"))
	runner.RunFile(path.Join(filename, "def_function.hummus"))
	runner.RunFile(path.Join(filename, "def_sub.hummus"))
	runner.RunFile(path.Join(filename, "/local_import/main.hummus"))
	runner.RunFile(path.Join(filename, "fib.hummus"))
	runner.RunFile(path.Join(filename, "free.hummus"))
	runner.RunFile(path.Join(filename, "identity.hummus"))
	runner.RunFile(path.Join(filename, "import.hummus"))
	runner.RunFile(path.Join(filename, "loop.hummus"))
	runner.RunFile(path.Join(filename, "logging.hummus"))
	runner.RunFile(path.Join(filename, "map_filter_reduce.hummus"))
	runner.RunFile(path.Join(filename, "maps.hummus"))
	runner.RunFile(path.Join(filename, "math.hummus"))
	runner.RunFile(path.Join(filename, "option.hummus"))
	runner.RunFile(path.Join(filename, "pipes.hummus"))
	runner.RunFile(path.Join(filename, "structs.hummus"))
	runner.RunFile(path.Join(filename, "variables.hummus"))
	runner.RunFile(path.Join(filename, "watch.hummus"))
	runner.RunFile(path.Join(filename, "watch2.hummus"))
}

func TestRunString(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	interpreter.BasePath = path.Join(path.Dir(filename), "../bin/")

	res := runner.RunString(`
(use :<base>)
(use :<str>)

(def fruits (list "Apple" "Mango" "Strawberry" "Orange"))

(map fruits (fn x
	(str/lower x)
))
`)

	if res.NodeType != interpreter.NODETYPE_LIST {
		t.FailNow()
	}

	resVals := res.Value.(interpreter.ListNode).Values

	fruits := []string{"apple", "mango", "strawberry", "orange"}

	for i := range resVals {
		if resVals[i].Value.(string) != fruits[i] {
			t.FailNow()
		}
	}
}
