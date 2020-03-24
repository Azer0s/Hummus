package interpreter

import (
	"fmt"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"sync"
	"time"
)

var channelMap = make(map[int]chan Node, 0)
var channelParent = make(map[int]int, 0)
var watcheeToWatcher = make(map[int]map[int]int, 0)
var watcherToWatchee = make(map[int]map[int]int, 0)
var currentPid = 0
var mu = &sync.Mutex{}

const MAILBOX_BUFFER = 1024

// CreatePidChannel creates a pid and a channel for a new process
func CreatePidChannel(self int) (pid int) {
	mu.Lock()

	channel := make(chan Node, MAILBOX_BUFFER)
	pid = currentPid

	channelParent[pid] = self

	channelMap[currentPid] = channel
	currentPid++
	mu.Unlock()

	return
}

func createWatch(watcher, watchee int) {
	if _, ok := channelMap[watchee]; !ok {
		panic(fmt.Sprintf("Process %d does not exist!", watchee))
	}

	if _, ok := watcheeToWatcher[watchee]; !ok {
		watcheeToWatcher[watchee] = make(map[int]int, 0)
	}

	if _, ok := watcherToWatchee[watcher]; !ok {
		watcherToWatchee[watcher] = make(map[int]int, 0)
	}

	watcheeToWatcher[watchee][watcher] = 1
	watcherToWatchee[watcher][watchee] = 1
}

func doWatch(arg Node, variables *map[string]Node) Node {
	if arg.NodeType != NODETYPE_INT {
		panic(SYSTEM_SYNC + " :watch expects an int as first argument!")
	}

	createWatch((*variables)[SELF].Value.(int), arg.Value.(int))
	return Node{
		Value:    0,
		NodeType: 0,
	}
}

func doSend(pid, val Node) Node {
	if pid.NodeType != NODETYPE_INT {
		panic(SYSTEM_SYNC + " :send expects an int as first argument!")
	}

	channel := channelMap[pid.Value.(int)]

	select {
	case channel <- val:
	default:
		//No data sent
	}

	return Node{
		Value:    0,
		NodeType: 0,
	}
}

func doReceive(variables *map[string]Node) Node {
	self := (*variables)[SELF].Value.(int)
	channel := channelMap[self]

	val := <-channel
	return val
}

func doCleanup(p int, r Node) {
	delete(channelParent, p)
	delete(channelMap, p)

	for watcher, doSend := range watcheeToWatcher[p] {
		if doSend != 1 {
			continue
		}

		channelMap[watcher] <- Node{
			Value: ListNode{Values: []Node{
				{
					Value:    "dead",
					NodeType: NODETYPE_ATOM,
				},
				{
					Value:    p,
					NodeType: NODETYPE_INT,
				},
				r,
			}},
			NodeType: NODETYPE_LIST,
		}
	}

	watchedByProcess := make([]int, 0)
	for k, _ := range watcheeToWatcher[p] {
		watchedByProcess = append(watchedByProcess, k)
	}

	watchingProcess := make([]int, 0)
	for k, _ := range watcherToWatchee[p] {
		watchingProcess = append(watchingProcess, k)
	}

	delete(watcherToWatchee, p)
	delete(watcheeToWatcher, p)

	for _, process := range watchedByProcess {
		delete(watcherToWatchee[process], p)
	}

	for _, process := range watchingProcess {
		delete(watcheeToWatcher[process], p)
	}
}

func doSpawn(arg Node, variables *map[string]Node) Node {
	if arg.NodeType != NODETYPE_FN {
		panic(SYSTEM_SYNC + " :spawn expects a function as first argument!")
	}

	ctx := make(map[string]Node, 0)
	for k, v := range *variables {
		ctx[k] = v
	}

	//Do global mutex when inserting into chan map
	pid := CreatePidChannel((*variables)[SELF].Value.(int))
	ctx[SELF] = Node{
		Value:    pid,
		NodeType: 0,
	}

	go func(p int) {
		defer func() {
			if r := recover(); r != nil {
				doCleanup(p, r.(Node))
			}
		}()

		doVariableCall(parser.Node{
			Type:      0,
			Arguments: []parser.Node{},
			Token:     lexer.Token{},
		}, arg, &ctx)
		panic(Node{
			Value:    0,
			NodeType: 0,
		})
	}(pid)

	return Node{
		Value:    pid,
		NodeType: 0,
	}
}

func doSleep(duration, mode Node) Node {
	if duration.NodeType != NODETYPE_INT {
		panic(SYSTEM_SYNC + " :sleep expects an int as first argument!")
	}

	if mode.NodeType != NODETYPE_ATOM {
		panic(SYSTEM_SYNC + " :sleep expects an int as first argument!")
	}

	d := time.Duration(int64(duration.Value.(int)))

	switch mode.Value.(string) {
	case "h":
		time.Sleep(d * time.Hour)
	case "min":
		time.Sleep(d * time.Minute)
	case "s":
		time.Sleep(d * time.Second)
	case "ms":
		time.Sleep(d * time.Millisecond)
	default:
		panic(SYSTEM_SYNC + " :sleep only accepts :h, :min, :s or :ms as second argument!")
	}

	return Node{
		Value:    0,
		NodeType: 0,
	}
}

func doUnwatch(watchee Node, self int) Node {
	if watchee.NodeType != NODETYPE_INT {
		panic(SYSTEM_SYNC + " :unwatch expects an int as first argument!")
	}

	w := watchee.Value.(int)
	delete(watcheeToWatcher[w], self)
	delete(watcherToWatchee[self], w)

	return Node{
		Value:    0,
		NodeType: 0,
	}
}

func doSystemCallSync(node parser.Node, variables *map[string]Node) Node {
	args := resolve(node.Arguments, variables, node.Token.Line)

	mode := args[0].Value.(string)

	switch mode {
	case "die":
		panic(args[1])
	case "watch":
		return doWatch(args[1], variables)
	case "unwatch":
		return doUnwatch(args[1], (*variables)[SELF].Value.(int))
	case "send":
		return doSend(args[1], args[2])
	case "receive":
		return doReceive(variables)
	case "spawn":
		return doSpawn(args[1], variables)
	case "sleep":
		return doSleep(args[1], args[2])
	default:
		panic("Unrecognized mode")
	}
}
