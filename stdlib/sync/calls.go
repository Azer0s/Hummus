package main

import (
	"fmt"
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/lexer"
	"github.com/Azer0s/Hummus/parser"
	"sync"
	"time"
)

var CALL string = "--system-do-sync!"

var channelMap = make(map[int]chan interpreter.Node, 0)
var channelParent = make(map[int]int, 0)
var watcheeToWatcher = make(map[int]map[int]int, 0)
var watcherToWatchee = make(map[int]map[int]int, 0)
var currentPid = 0

var currentPidMu = &sync.RWMutex{}
var channelMapMu = &sync.RWMutex{}
var channelParentMu = &sync.RWMutex{}
var watcheeToWatcherMu = &sync.RWMutex{}
var watcherToWatcheeMu = &sync.RWMutex{}

// MAILBOX_BUFFER buffer for channel mailboxes
const MAILBOX_BUFFER = 1024

func createPidChannel(self int) (pid int) {
	channel := make(chan interpreter.Node, MAILBOX_BUFFER)

	currentPidMu.Lock()

	channelMapMu.RLock()
	for {
		if _, ok := channelMap[currentPid]; !ok {
			channelMapMu.RUnlock()
			break
		}

		currentPid++
	}

	pid = currentPid

	channelParentMu.Lock()
	channelParent[pid] = self
	channelParentMu.Unlock()

	channelMapMu.Lock()
	channelMap[pid] = channel
	channelMapMu.Unlock()

	currentPid++

	currentPidMu.Unlock()

	return
}

func createWatch(watcher, watchee int) {
	channelMapMu.RLock()
	if _, ok := channelMap[watchee]; !ok {
		panic(fmt.Sprintf("Process %d does not exist!", watchee))
	}
	channelMapMu.RUnlock()

	watcheeToWatcherMu.Lock()
	if _, ok := watcheeToWatcher[watchee]; !ok {
		watcheeToWatcher[watchee] = make(map[int]int, 0)
	}
	watcheeToWatcherMu.Unlock()

	watcherToWatcheeMu.Lock()
	if _, ok := watcherToWatchee[watcher]; !ok {
		watcherToWatchee[watcher] = make(map[int]int, 0)
	}
	watcherToWatcheeMu.Unlock()

	watcheeToWatcherMu.Lock()
	watcheeToWatcher[watchee][watcher] = 1
	watcheeToWatcherMu.Unlock()

	watcherToWatcheeMu.Lock()
	watcherToWatchee[watcher][watchee] = 1
	watcherToWatcheeMu.Unlock()

}

func doWatch(arg interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	if arg.NodeType != interpreter.NODETYPE_INT {
		panic(CALL + " :watch expects an int as first argument!")
	}

	createWatch((*variables)[interpreter.SELF].Value.(int), arg.Value.(int))
	return interpreter.Node{
		Value:    0,
		NodeType: 0,
	}
}

func doSend(pid, val interpreter.Node) interpreter.Node {
	if pid.NodeType != interpreter.NODETYPE_INT {
		panic(CALL + " :send expects an int as first argument!")
	}

	channel := channelMap[pid.Value.(int)]

	select {
	case channel <- val:
	default:
		//No data sent
	}

	return interpreter.Node{
		Value:    0,
		NodeType: 0,
	}
}

func doReceive(variables *map[string]interpreter.Node) interpreter.Node {
	self := (*variables)[interpreter.SELF].Value.(int)

	channelMapMu.RLock()
	channel := channelMap[self]
	channelMapMu.RUnlock()

	val := <-channel
	return val
}

func doCleanup(p int, r interpreter.Node) {
	channelParentMu.Lock()
	delete(channelParent, p)
	channelParentMu.Unlock()

	channelMapMu.Lock()
	delete(channelMap, p)
	channelMapMu.Unlock()

	watcheeToWatcherMu.RLock()
	for watcher, doSend := range watcheeToWatcher[p] {
		if doSend != 1 {
			continue
		}

		channelMap[watcher] <- interpreter.Node{
			Value: interpreter.ListNode{Values: []interpreter.Node{
				{
					Value:    "dead",
					NodeType: interpreter.NODETYPE_ATOM,
				},
				{
					Value:    p,
					NodeType: interpreter.NODETYPE_INT,
				},
				r,
			}},
			NodeType: interpreter.NODETYPE_LIST,
		}
	}
	watcheeToWatcherMu.RUnlock()

	watcheeToWatcherMu.RLock()
	watchedByProcess := make([]int, 0)
	for k := range watcheeToWatcher[p] {
		watchedByProcess = append(watchedByProcess, k)
	}
	watcheeToWatcherMu.RUnlock()

	watcherToWatcheeMu.RLock()
	watchingProcess := make([]int, 0)
	for k := range watcherToWatchee[p] {
		watchingProcess = append(watchingProcess, k)
	}
	watcherToWatcheeMu.RUnlock()

	watcherToWatcheeMu.Lock()
	delete(watcherToWatchee, p)
	watcherToWatcheeMu.Unlock()

	watcheeToWatcherMu.Lock()
	delete(watcheeToWatcher, p)
	watcheeToWatcherMu.Unlock()

	watcherToWatcheeMu.Lock()
	for _, process := range watchedByProcess {
		delete(watcherToWatchee[process], p)
	}
	watcherToWatcheeMu.Unlock()

	watcheeToWatcherMu.Lock()
	for _, process := range watchingProcess {
		delete(watcheeToWatcher[process], p)
	}
	watcheeToWatcherMu.Unlock()
}

func doSpawn(arg interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	if arg.NodeType != interpreter.NODETYPE_FN {
		panic(CALL + " :spawn expects a function as first argument!")
	}

	ctx := make(map[string]interpreter.Node, 0)
	for k, v := range *variables {
		ctx[k] = v
	}

	//Do global mutex when inserting into chan map
	pid := createPidChannel((*variables)[interpreter.SELF].Value.(int))
	ctx[interpreter.SELF] = interpreter.Node{
		Value:    pid,
		NodeType: 0,
	}

	go func(p int) {
		defer func() {
			if r := recover(); r != nil {
				if val, ok := r.(string); ok {
					doCleanup(p, interpreter.Node{
						Value:    val,
						NodeType: interpreter.NODETYPE_STRING,
					})
				} else {
					doCleanup(p, r.(interpreter.Node))
				}
			}
		}()

		interpreter.DoVariableCall(parser.Node{
			Type:      0,
			Arguments: []parser.Node{},
			Token:     lexer.Token{},
		}, arg, &ctx)
		panic(interpreter.Node{
			Value:    0,
			NodeType: 0,
		})
	}(pid)

	return interpreter.Node{
		Value:    pid,
		NodeType: 0,
	}
}

func doSleep(duration, mode interpreter.Node) interpreter.Node {
	if duration.NodeType != interpreter.NODETYPE_INT {
		panic(CALL + " :sleep expects an int as first argument!")
	}

	if mode.NodeType != interpreter.NODETYPE_ATOM {
		panic(CALL + " :sleep expects an int as first argument!")
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
		panic(CALL + " :sleep only accepts :h, :min, :s or :ms as second argument!")
	}

	return interpreter.Node{
		Value:    0,
		NodeType: 0,
	}
}

func doUnwatch(watchee interpreter.Node, self int) interpreter.Node {
	if watchee.NodeType != interpreter.NODETYPE_INT {
		panic(CALL + " :unwatch expects an int as first argument!")
	}

	w := watchee.Value.(int)
	watcheeToWatcherMu.Lock()
	delete(watcheeToWatcher[w], self)
	watcheeToWatcherMu.Unlock()

	watcherToWatcheeMu.Lock()
	delete(watcherToWatchee[self], w)
	watcherToWatcheeMu.Unlock()

	return interpreter.Node{
		Value:    0,
		NodeType: 0,
	}
}

func Init(variables *map[string]interpreter.Node) {
	(*variables)[interpreter.SELF] = interpreter.Node{
		Value:    createPidChannel(0),
		NodeType: interpreter.NODETYPE_INT,
	}
}

func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "die":
		panic(args[1])
	case "watch":
		return doWatch(args[1], variables)
	case "unwatch":
		return doUnwatch(args[1], (*variables)[interpreter.SELF].Value.(int))
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
