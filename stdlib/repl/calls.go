package main

import (
	"fmt"
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/project"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {}

// CALL string functions
var CALL string = "--system-do-repl!"

var docMap = make(map[string]string, 0)

// Init Hummus stdlib stub
func Init(variables *map[string]interpreter.Node) {
	docRe := regexp.MustCompile("(?m)(;; .*[\\n])*")
	functionRe := regexp.MustCompile(";; *([^ ]*)")
	labelRe := regexp.MustCompile(";; *[\\w /]+ *;;")

	log.Info("Parsing stdlib documentation...")

	err := filepath.Walk(path.Join(interpreter.BasePath, "stdlib"), func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".hummus") {
			buff, err := ioutil.ReadFile(path)

			if err != nil {
				return err
			}

			subMatches := docRe.FindAllString(string(buff), -1)

			filteredMatches := make([]string, 0, len(subMatches))
			for _, item := range subMatches {
				if item == "" || strings.Contains(item, "Copyright") || labelRe.MatchString(item) {
					continue
				}

				filteredMatches = append(filteredMatches, item)
			}

			if len(filteredMatches) > 1 {
				for _, doc := range filteredMatches {
					functionName := functionRe.FindStringSubmatch(strings.Split(doc, "\n")[0])[1]
					docMap[functionName] = doc
				}
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	log.Info("Parsed stdlib documentation successfully!")
}

// DoSystemCall Hummus stdlib stub
func DoSystemCall(args []interpreter.Node, variables *map[string]interpreter.Node) interpreter.Node {
	mode := args[0].Value.(string)

	switch mode {
	case "pull-lib":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :pull-lib")
		interpreter.EnsureType(&args, 2, interpreter.NODETYPE_STRING, CALL+" :pull-lib")

		p, err := os.Executable()

		if err != nil {
			panic(err)
		}

		project.PullPackage(path.Join(path.Dir(p), "lib"), args[1].Value.(string), args[2].Value.(string))

		return interpreter.Nothing

	case "help":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :help")

		functionName := args[1].Value.(string)
		doc, ok := docMap[functionName]

		if !ok {
			panic("Unrecognized function " + functionName + "!")
		}

		fmt.Print(doc)

		return interpreter.Nothing

	case "help-group":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :help-group")

		subName := args[1].Value.(string)
		docs := make([]string, 0)

		for name, doc := range docMap {
			if strings.Split(name, "/")[0] == subName {
				docs = append(docs, doc)
			}
		}

		fmt.Print(strings.Join(docs, "\n"))

		return interpreter.Nothing

	case "help-ungrouped":
		docs := make([]string, 0)

		for name, doc := range docMap {
			if !strings.Contains(name, "/") || name == "/" {
				docs = append(docs, doc)
			}
		}

		fmt.Print(strings.Join(docs, "\n"))

		return interpreter.Nothing

	case "search-fn":
		interpreter.EnsureType(&args, 1, interpreter.NODETYPE_STRING, CALL+" :search-fn")

		docs := make([]string, 0)
		searchString := args[1].Value.(string)

		for name, doc := range docMap {
			if strings.Contains(name, searchString) || strings.Contains(doc, searchString) {
				docs = append(docs, strings.ReplaceAll(doc, searchString, fmt.Sprintf("\033[31m%s\033[0m", searchString)))
			}
		}

		fmt.Print(strings.Join(docs, "\n"))

		return interpreter.Nothing

	default:
		panic("Unrecognized mode")
	}
}
